package storage

import (
	"backend/gen/storage"
	"backend/internal/core"
	storagedb "backend/internal/storage/db"
	"backend/internal/storage/events"
	"context"
	"fmt"
	"log/slog"
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
	interval time.Duration
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
		// This could be made configurable later
		interval: 5 * time.Second,
	}
}

func (p *OutboxProcessor) Start(
	ctx context.Context,
) error {
	ticker := time.NewTicker(p.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if err := p.process(ctx); err != nil {
				slog.ErrorContext(ctx, "failed to process outbox", "error", err)
			}
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

	events, err := p.db.GetOutboxUnpublishedEvents(ctx, 100)
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

	for _, event := range events {
		if err := p.publish(ctx, event); err != nil {
			slog.ErrorContext(ctx, "failed to publish outbox event",
				"event_id", event.EventID, "error", err)
			failureCount++
			continue
		}
		successCount++
	}

	span.SetAttributes(
		attribute.Int("outbox.published_count", successCount),
		attribute.Int("outbox.failed_count", failureCount),
	)

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

	if err := p.db.MarkEventAsPublished(ctx, event.EventID); err != nil {
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
