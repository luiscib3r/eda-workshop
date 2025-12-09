package ocr

import (
	"backend/gen/ocr"
	"backend/internal/infrastructure/nats"
	ocrdb "backend/internal/ocr/db"
	ocrev "backend/internal/ocr/events"
	"backend/internal/storage/events"
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/oklog/ulid/v2"
	"github.com/samber/lo"
	"go.opentelemetry.io/otel"
	"google.golang.org/protobuf/encoding/protojson"
)

type FilesDeletedConsumer struct {
	*nats.NatsConsumer[*events.FilesDeletedEvent]
	db   *ocrdb.Queries
	pool *pgxpool.Pool
}

func NewFilesDeletedConsumer(
	js jetstream.JetStream,
	db *ocrdb.Queries,
	pool *pgxpool.Pool,
) *FilesDeletedConsumer {
	name := "ocr_files_deleted_consumer"
	numWorkers := 4
	workerBufferSize := 20

	consumer := &FilesDeletedConsumer{
		db:   db,
		pool: pool,
	}

	consumer.NatsConsumer = nats.NewNatsConsumer(
		name,
		events.STORAGE_CHANNEL,
		events.STORAGE_FILES_DELETED_EVENT,
		numWorkers,
		workerBufferSize,
		events.NewFilesDeletedEventFromMessage,
		consumer.handler,
		js,
		jetstream.ConsumerConfig{
			Name:          name,
			Durable:       name,
			Description:   "OCR Files Deleted Event Consumer",
			FilterSubject: events.STORAGE_FILES_DELETED_EVENT,
		},
	)

	return consumer
}

func (c *FilesDeletedConsumer) handler(
	ctx context.Context,
	event *events.FilesDeletedEvent,
) error {
	tracer := otel.Tracer("ocr.FilesDeletedConsumer")
	ctx, span := tracer.Start(ctx, "FilesDeletedConsumer.handler")
	defer span.End()

	tx, err := c.pool.Begin(ctx)
	if err != nil {
		span.RecordError(err)
		return err
	}
	defer tx.Rollback(ctx)

	qtx := c.db.WithTx(tx)

	ids := event.Payload.FileKeys
	for _, idStr := range ids {
		id, err := uuid.Parse(idStr)
		if err != nil {
			span.RecordError(err)
			return err
		}
		err = qtx.DeleteFilePagesByFileID(ctx, pgtype.UUID{
			Bytes: id,
			Valid: true,
		})
		if err != nil {
			span.RecordError(err)
			return err
		}
	}

	fileKeys := lo.Map(ids, func(key string, _ int) string {
		fileKey, err := uuid.Parse(key)
		if err != nil {
			return ""
		}
		return ulid.ULID(fileKey).String()
	})
	fileKeys = lo.Reject(fileKeys, func(id string, _ int) bool {
		return id == ""
	})
	ev := ocrev.NewFilePagesDeletedEvent(
		&ocr.FilePagesDeletedEventData{
			FileKeys: fileKeys,
		},
	)

	eventId := ev.Id
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
	if err != nil {
		span.RecordError(err)
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		span.RecordError(err)
		return err
	}

	return nil
}
