package ocr

import (
	"backend/gen/ocr"
	"backend/internal/core"
	ocrdb "backend/internal/ocr/db"
	"backend/internal/ocr/events"
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
	db       *ocrdb.Queries
	producer *OcrProducer
}

func NewOutboxProcessor(
	pool *pgxpool.Pool,
	db *ocrdb.Queries,
	producer *OcrProducer,
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
	_, err = conn.Exec(ctx, "LISTEN ocr_outbox_channel")
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
	tracer := otel.Tracer("ocr_outbox_processor")
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
	event ocrdb.GetOutboxUnpublishedEventsRow,
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
	event ocrdb.GetOutboxUnpublishedEventsRow,
) (core.EventSpec, error) {
	switch event.EventType {
	case events.FILE_PAGES_DELETED_EVENT:
		data := &ocr.FilePagesDeletedEventData{}
		if err := protojson.Unmarshal(event.Payload, data); err != nil {
			return nil, err
		}
		return &events.FilePagesDeletedEvent{
			Id:      event.EventID.Bytes,
			Payload: data,
		}, nil
	case events.FILE_PAGE_REGISTERED_EVENT:
		data := &ocr.FilePageRegisteredEventData{}
		if err := protojson.Unmarshal(event.Payload, data); err != nil {
			return nil, err
		}
		return &events.FilePageRegisteredEvent{
			Id:      event.EventID.Bytes,
			Payload: data,
		}, nil
	case events.FILE_PAGE_OCR_GENERATED_EVENT:
		data := &ocr.FilePageOcrGeneratedEventData{}
		if err := protojson.Unmarshal(event.Payload, data); err != nil {
			return nil, err
		}
		return &events.FilePageOcrGeneratedEvent{
			Id:      event.EventID.Bytes,
			Payload: data,
		}, nil
	}

	return nil, fmt.Errorf("unknown outbox event type: %s", event.EventType)
}
