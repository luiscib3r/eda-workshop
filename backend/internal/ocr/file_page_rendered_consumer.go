package ocr

import (
	"backend/gen/ocr"
	"backend/internal/infrastructure/nats"
	ocrdb "backend/internal/ocr/db"
	"backend/internal/ocr/events"
	"context"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/oklog/ulid/v2"
	"go.opentelemetry.io/otel"
	"google.golang.org/protobuf/encoding/protojson"
)

type FilePageRenderedConsumer struct {
	*nats.NatsConsumer[*events.FilePageRenderedEvent]
	db   *ocrdb.Queries
	pool *pgxpool.Pool
}

func NewFilePageRenderedConsumer(
	js jetstream.JetStream,
	db *ocrdb.Queries,
	pool *pgxpool.Pool,
) *FilePageRenderedConsumer {
	name := "ocr_file_page_rendered_consumer"
	numWorkers := 4
	workerBufferSize := 20

	consumer := &FilePageRenderedConsumer{
		db:   db,
		pool: pool,
	}

	consumer.NatsConsumer = nats.NewNatsConsumer(
		name,
		events.OCR_CHANNEL,
		events.FILE_PAGE_RENDERED_EVENT,
		numWorkers,
		workerBufferSize,
		events.NewFilePageRenderedEventFromMessage,
		consumer.handler,
		js,
		jetstream.ConsumerConfig{
			Name:          name,
			Durable:       name,
			Description:   "OCR File Page Rendered Event Consumer",
			FilterSubject: events.FILE_PAGE_RENDERED_EVENT,
			DeliverPolicy: jetstream.DeliverNewPolicy,
		},
	)

	return consumer
}

func (c *FilePageRenderedConsumer) handler(
	ctx context.Context,
	event *events.FilePageRenderedEvent,
) error {
	tracer := otel.Tracer("file_page_rendered_consumer")
	ctx, span := tracer.Start(
		ctx,
		"FilePageRenderedConsumer.handler",
	)
	defer span.End()

	tx, err := c.pool.Begin(ctx)
	if err != nil {
		span.RecordError(err)
		return err
	}
	defer tx.Rollback(ctx)

	qtx := c.db.WithTx(tx)

	id, err := ulid.Parse(event.Payload.PageKey)
	if err != nil {
		span.RecordError(err)
		return err
	}

	fileKey, err := ulid.Parse(event.Payload.FileKey)
	if err != nil {
		span.RecordError(err)
		return err
	}

	err = qtx.CreateFilePage(ctx, ocrdb.CreateFilePageParams{
		ID: pgtype.UUID{
			Bytes: id,
			Valid: true,
		},
		FileID: pgtype.UUID{
			Bytes: fileKey,
			Valid: true,
		},
		PageNumber:   event.Payload.PageNumber,
		PageImageKey: event.Payload.PageImageKey,
	})
	if err != nil {
		span.RecordError(err)
		return err
	}

	ev := events.NewFilePageRegisteredEvent(
		&ocr.FilePageRegisteredEventData{
			Id:           id.String(),
			FileId:       fileKey.String(),
			PageNumber:   event.Payload.PageNumber,
			PageImageKey: event.Payload.PageImageKey,
		},
	)

	eventId := event.Id
	eventType := ev.Type()
	payload, err := protojson.Marshal(ev.Payload)
	if err != nil {
		span.RecordError(err)
		return err
	}

	err = qtx.CreateOutboxEvent(ctx, ocrdb.CreateOutboxEventParams{
		EventID: pgtype.UUID{
			Bytes: eventId,
			Valid: true,
		},
		EventType: eventType,
		Payload:   payload,
	})

	if tx.Commit(ctx); err != nil {
		span.RecordError(err)
		return err
	}

	return nil
}
