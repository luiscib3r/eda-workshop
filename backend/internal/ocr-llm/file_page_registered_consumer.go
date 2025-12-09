package ocrllm

import (
	"backend/gen/ocr"
	"backend/internal/infrastructure/nats"
	"backend/internal/infrastructure/storage"
	ocrdb "backend/internal/ocr/db"
	"backend/internal/ocr/events"
	"context"
	"io"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/oklog/ulid/v2"
	"go.opentelemetry.io/otel"
	"google.golang.org/protobuf/encoding/protojson"
)

type FilePageRegisteredConsumer struct {
	*nats.NatsConsumer[*events.FilePageRegisteredEvent]
	db   *ocrdb.Queries
	pool *pgxpool.Pool
	ocr  *OcrAgent
	s3   *s3.Client
}

func NewFilePageRegisteredConsumer(
	js jetstream.JetStream,
	db *ocrdb.Queries,
	pool *pgxpool.Pool,
	ocr *OcrAgent,
	s3 *s3.Client,
) *FilePageRegisteredConsumer {
	name := "ocr_file_page_registered_consumer"
	numWorkers := 4
	workerBufferSize := 20

	consumer := &FilePageRegisteredConsumer{
		db:   db,
		pool: pool,
		ocr:  ocr,
		s3:   s3,
	}

	consumer.NatsConsumer = nats.NewNatsConsumer(
		name,
		events.OCR_CHANNEL,
		events.FILE_PAGE_REGISTERED_EVENT,
		numWorkers,
		workerBufferSize,
		events.NewFilePageRegisteredEventFromMessage,
		consumer.handler,
		js,
		jetstream.ConsumerConfig{
			Name:          name,
			Durable:       name,
			Description:   "OCR File Page Registered Event Consumer",
			FilterSubject: events.FILE_PAGE_REGISTERED_EVENT,
			DeliverPolicy: jetstream.DeliverNewPolicy,
		},
	)

	return consumer
}

func (c *FilePageRegisteredConsumer) handler(
	ctx context.Context,
	event *events.FilePageRegisteredEvent,
) error {
	// Start tracing span
	tracer := otel.Tracer("file_page_registered_consumer")
	ctx, span := tracer.Start(
		ctx,
		"FilePageRegisteredConsumer.handler",
	)
	defer span.End()

	// Validate UUID
	id, err := ulid.Parse(event.Payload.Id)
	if err != nil {
		span.RecordError(err)
		return err
	}

	// Get page image key
	pageKey := event.Payload.PageImageKey

	// Fetch image from S3
	result, err := c.s3.GetObject(ctx, &s3.GetObjectInput{
		Key:    &pageKey,
		Bucket: aws.String(storage.BUCKET_NAME),
	})
	if err != nil {
		span.RecordError(err)
		return err
	}
	defer result.Body.Close()

	// Get image data
	data, err := io.ReadAll(result.Body)
	if err != nil {
		span.RecordError(err)
		return err
	}

	// Generate OCR
	resp, err := c.ocr.Invoke(ctx, data)
	if err != nil {
		span.RecordError(err)
		return err
	}

	// Extract OCR text
	textContent := "Not recognized"
	if len(resp.Choices) > 0 {
		textContent = resp.Choices[0].Message.Content
	}

	// Begin db transaction
	tx, err := c.pool.Begin(ctx)
	if err != nil {
		span.RecordError(err)
		return err
	}
	defer tx.Rollback(ctx)
	qtx := c.db.WithTx(tx)

	// Store OCR result in DB
	if err := qtx.UpdateFilePageText(ctx, ocrdb.UpdateFilePageTextParams{
		ID: pgtype.UUID{
			Bytes: id,
			Valid: true,
		},
		TextContent: &textContent,
	}); err != nil {
		span.RecordError(err)
		return err
	}

	// Create FilePageOcrGeneratedEvent
	ev := events.NewFilePageOcrGeneratedEvent(
		&ocr.FilePageOcrGeneratedEventData{
			Id:           event.Payload.Id,
			FileId:       event.Payload.FileId,
			PageNumber:   event.Payload.PageNumber,
			PageImageKey: event.Payload.PageImageKey,
		},
	)

	eventId := ev.Id
	eventType := ev.Type()
	payload, err := protojson.Marshal(ev.Payload)
	if err != nil {
		span.RecordError(err)
		return err
	}

	// Save outbox event
	err = qtx.CreateOutboxEvent(ctx, ocrdb.CreateOutboxEventParams{
		EventID: pgtype.UUID{
			Bytes: eventId,
			Valid: true,
		},
		EventType: eventType,
		Payload:   payload,
	})

	// Commit transaction
	if tx.Commit(ctx); err != nil {
		span.RecordError(err)
		return err
	}

	return nil
}
