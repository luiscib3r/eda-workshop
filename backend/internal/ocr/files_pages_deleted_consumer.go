package ocr

import (
	"backend/internal/infrastructure/nats"
	"backend/internal/infrastructure/storage"
	"backend/internal/ocr/events"
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/nats-io/nats.go/jetstream"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"golang.org/x/sync/errgroup"
)

const (
	maxDeleteBatch    = 1000
	concurrentDeletes = 5
)

type FilePagesDeletedConsumer struct {
	*nats.NatsConsumer[*events.FilePagesDeletedEvent]
	s3 *s3.Client
}

func NewFilePagesDeletedConsumer(
	js jetstream.JetStream,
	s3 *s3.Client,
) *FilePagesDeletedConsumer {
	name := "ocr_file_pages_deleted_consumer"
	numWorkers := 5
	workerBufferSize := 10

	consumer := &FilePagesDeletedConsumer{
		s3: s3,
	}

	consumer.NatsConsumer = nats.NewNatsConsumer(
		name,
		events.OCR_CHANNEL,
		events.FILE_PAGES_DELETED_EVENT,
		numWorkers,
		workerBufferSize,
		events.NewFilePagesDeletedEventFromMessage,
		consumer.handler,
		js,
		jetstream.ConsumerConfig{
			Name:          name,
			Durable:       name,
			Description:   "OCR File Pages Deleted Event Consumer",
			FilterSubject: events.FILE_PAGES_DELETED_EVENT,
		},
	)

	return consumer
}

func (c *FilePagesDeletedConsumer) handler(
	ctx context.Context,
	event *events.FilePagesDeletedEvent,
) error {
	tracer := otel.Tracer("ocr.FilePagesDeletedConsumer")
	ctx, span := tracer.Start(ctx, "FilePagesDeletedConsumer.handler")
	defer span.End()

	g, ctx := errgroup.WithContext(ctx)
	g.SetLimit(concurrentDeletes)

	for _, fileKey := range event.Payload.FileKeys {
		fileKey := fileKey
		g.Go(func() error {
			return c.delete(ctx, fileKey)
		})
	}

	if err := g.Wait(); err != nil {
		span.RecordError(err)
		return err
	}

	return nil
}

func (c *FilePagesDeletedConsumer) delete(
	ctx context.Context,
	fileKey string,
) error {
	// get span from context
	tracer := otel.Tracer("ocr.FilePagesDeletedConsumer.delete")
	ctx, span := tracer.Start(ctx, "FilePagesDeletedConsumer.delete")
	defer span.End()

	span.SetAttributes(
		attribute.String("file_key", fileKey),
	)

	prefix := PageImagePrefix(fileKey)

	var objects []types.ObjectIdentifier
	paginator := s3.NewListObjectsV2Paginator(c.s3, &s3.ListObjectsV2Input{
		Bucket: aws.String(storage.BUCKET_NAME),
		Prefix: aws.String(prefix),
	})

	for paginator.HasMorePages() {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		page, err := paginator.NextPage(ctx)
		if err != nil {
			span.RecordError(err)
			return fmt.Errorf("error listando objetos para %s: %w", fileKey, err)
		}

		for _, obj := range page.Contents {
			objects = append(objects, types.ObjectIdentifier{
				Key: obj.Key,
			})
		}
	}

	objects = append(objects, types.ObjectIdentifier{
		Key: aws.String(prefix),
	})

	for i := 0; i < len(objects); i += maxDeleteBatch {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		end := min(i+maxDeleteBatch, len(objects))

		_, err := c.s3.DeleteObjects(ctx, &s3.DeleteObjectsInput{
			Bucket: aws.String(storage.BUCKET_NAME),
			Delete: &types.Delete{
				Objects: objects[i:end],
				Quiet:   aws.Bool(true),
			},
		})
		if err != nil {
			span.RecordError(err)
			return fmt.Errorf("error eliminando lote para %s: %w", fileKey, err)
		}
	}

	return nil
}
