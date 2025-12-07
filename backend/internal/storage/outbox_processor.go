package storage

import (
	"backend/gen/storage"
	"backend/internal/core"
	storagedb "backend/internal/storage/db"
	"backend/internal/storage/events"
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/protobuf/encoding/protojson"
)

type OutboxProcessor struct {
	pool     *pgxpool.Pool
	db       *storagedb.Queries
	producer *StorageProducer
}

func NewOutboxProcessor(
	pool *pgxpool.Pool,
	db *storagedb.Queries,
	producer *StorageProducer,
) *OutboxProcessor {
	return &OutboxProcessor{
		pool:     pool,
		db:       db,
		producer: producer,
	}
}

func (p *OutboxProcessor) Start(
	ctx context.Context,
) error {
	conn, err := p.pool.Acquire(ctx)
	if err != nil {
		return err
	}
	defer conn.Release()

	// Listen to outbox notifications
	_, err = conn.Exec(ctx, "LISTEN storage_outbox_channel")
	if err != nil {
		return err
	}

	// Outbox notifications channel
	notifyChan := make(chan struct{})
	go func() {
		for {
			_, err := conn.Conn().WaitForNotification(ctx)
			if err != nil {
				if ctx.Err() != nil {
					return // Context canceled, exit
				}
				continue // Ignore errors and continue listening
			}
			notifyChan <- struct{}{}
		}
	}()

	// Initial backlog
	p.process(ctx)

	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-notifyChan:
			p.process(ctx)
		case <-ticker.C:
			p.process(ctx)
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

func (p *OutboxProcessor) process(ctx context.Context) error {
	tracer := otel.Tracer("storage_outbox_processor")
	ctx, span := tracer.Start(
		ctx,
		"OutboxProcessor.process",
		trace.WithSpanKind(trace.SpanKindInternal),
	)
	defer span.End()

	tx, err := p.pool.Begin(ctx)
	if err != nil {
		span.RecordError(err)
		return err
	}
	defer tx.Rollback(ctx)

	qtx := p.db.WithTx(tx)

	events, err := qtx.GetOutboxUnpublishedEvents(ctx, 100)
	if err != nil {
		span.RecordError(err)
		return err
	}

	span.SetAttributes(attribute.Int("outbox.events_count", len(events)))

	if len(events) == 0 {
		return nil
	}

	successCount := 0
	failureCount := 0
	publishedButNotMarked := 0

	for _, event := range events {
		if err := p.publish(ctx, event); err != nil {
			span.RecordError(err)
			failureCount++
			continue
		}
		if err := qtx.MarkEventAsPublished(ctx, event.EventID); err != nil {
			span.RecordError(err)
			publishedButNotMarked++
			continue
		}
		successCount++
	}

	span.SetAttributes(
		attribute.Int("outbox.published_count", successCount),
		attribute.Int("outbox.failed_count", failureCount),
		attribute.Int("outbox.published_but_not_marked", publishedButNotMarked),
	)

	if err := tx.Commit(ctx); err != nil {
		span.RecordError(err)
		return err
	}

	return nil
}

func (p *OutboxProcessor) publish(
	ctx context.Context,
	event storagedb.GetOutboxUnpublishedEventsRow,
) error {
	ev, err := p.event(event)
	if err != nil {
		return err
	}

	if err := p.producer.Publish(ctx, ev); err != nil {
		return err
	}

	return nil
}

func (p *OutboxProcessor) event(
	event storagedb.GetOutboxUnpublishedEventsRow,
) (core.EventSpec, error) {
	switch event.EventType {
	case events.STORAGE_FILE_UPLOADED_EVENT:
		data := &storage.FileUploadedEventData{}
		if err := protojson.Unmarshal(event.Payload, data); err != nil {
			return nil, err
		}
		return &events.FileUploadedEvent{
			Id:      event.EventID,
			Payload: data,
		}, nil
	case events.STORAGE_FILES_DELETED_EVENT:
		data := &storage.FilesDeletedEventData{}
		if err := protojson.Unmarshal(event.Payload, data); err != nil {
			return nil, err
		}
		return &events.FilesDeletedEvent{
			Id:      event.EventID,
			Payload: data,
		}, nil
	}

	return nil, fmt.Errorf("unknown outbox event type: %s", event.EventType)
}
