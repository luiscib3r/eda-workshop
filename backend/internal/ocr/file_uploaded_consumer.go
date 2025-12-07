package ocr

import (
	"backend/gen/ocr"
	"backend/internal/infrastructure/nats"
	"backend/internal/infrastructure/storage"
	ocrev "backend/internal/ocr/events"
	"backend/internal/storage/events"
	"bytes"
	"context"
	"image/png"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gen2brain/go-fitz"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/oklog/ulid/v2"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

type FileUploadedConsumer struct {
	*nats.NatsConsumer[*events.FileUploadedEvent]
	s3       *s3.Client
	producer *OcrProducer
}

func NewFileUploadedConsumer(
	js jetstream.JetStream,
	s3 *s3.Client,
	producer *OcrProducer,
) *FileUploadedConsumer {
	name := "ocr_file_uploaded_consumer"

	numWorkers := 10
	workerBufferSize := 5
	consumer := &FileUploadedConsumer{
		s3:       s3,
		producer: producer,
	}

	consumer.NatsConsumer = nats.NewNatsConsumer(
		name,
		events.STORAGE_CHANNEL,
		events.STORAGE_FILE_UPLOADED_EVENT,
		numWorkers,
		workerBufferSize,
		events.NewFileUploadedEventFromMessage,
		consumer.handler,
		js,
		jetstream.ConsumerConfig{
			Name:          name,
			Durable:       name,
			Description:   "OCR File Uploaded Event Consumer",
			FilterSubject: events.STORAGE_FILE_UPLOADED_EVENT,
			DeliverPolicy: jetstream.DeliverNewPolicy,
		},
	)

	return consumer
}

func (c *FileUploadedConsumer) handler(
	ctx context.Context,
	event *events.FileUploadedEvent,
) error {
	tracer := otel.Tracer("ocr_file_uploaded_consumer")

	ctx, span := tracer.Start(
		ctx,
		"FileUploadedConsumer.handler",
		trace.WithSpanKind(trace.SpanKindInternal),
	)
	defer span.End()

	// Download the file from S3 using the provided S3 client
	result, err := c.s3.GetObject(ctx, &s3.GetObjectInput{
		Key:    &event.Payload.FileKey,
		Bucket: aws.String(storage.BUCKET_NAME),
	})
	if err != nil {
		return err
	}
	defer result.Body.Close()

	// Transform PDF to image
	doc, err := fitz.NewFromReader(result.Body)
	if err != nil {
		return err
	}
	defer doc.Close()

	pageCount := doc.NumPage()

	for pageNum := range pageCount {
		// Render page to image
		img, err := doc.Image(pageNum)
		if err != nil {
			span.RecordError(err)
			continue
		}

		// Convert image to PNG bytes
		buf := new(bytes.Buffer)
		err = png.Encode(buf, img)
		if err != nil {
			span.RecordError(err)
			continue
		}

		// Upload the image to S3
		key := ulid.MustNew(
			ulid.Timestamp(time.Now()),
			ulid.DefaultEntropy(),
		).String()
		pageImageKey := PageImageKey(event.Payload.FileKey, key)

		if _, err = c.s3.PutObject(ctx, &s3.PutObjectInput{
			Bucket:      aws.String(storage.BUCKET_NAME),
			Key:         aws.String(pageImageKey),
			Body:        bytes.NewReader(buf.Bytes()),
			ContentType: aws.String("image/png"),
		}); err != nil {
			span.RecordError(err)
			continue
		}

		// Publish FilePageRenderedEvent
		event := ocrev.NewFilePageRenderedEvent(&ocr.FilePageRenderedEventData{
			FileKey:      event.Payload.FileKey,
			PageImageKey: pageImageKey,
			PageNumber:   int32(pageNum),
		})

		if err := c.producer.Publish(ctx, event); err != nil {
			span.RecordError(err)
			continue
		}
	}

	return nil
}
