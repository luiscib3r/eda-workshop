package nats

import (
	"backend/internal/core"
	"context"
	"fmt"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/protobuf/proto"
)

type NatsProducer struct {
	name    string
	channel string
	js      jetstream.JetStream
	cfg     jetstream.StreamConfig
}

func NewNatsProducer(
	name string,
	channel string,
	js jetstream.JetStream,
	cfg jetstream.StreamConfig,
) *NatsProducer {
	return &NatsProducer{
		name:    name,
		channel: channel,
		js:      js,
		cfg:     cfg,
	}
}

func (p *NatsProducer) CreateChannel(ctx context.Context) error {
	tracer := otel.Tracer(p.name)

	ctx, span := tracer.Start(
		ctx, fmt.Sprintf("%s.init", p.channel),
		trace.WithSpanKind(trace.SpanKindProducer),
		trace.WithAttributes(
			semconv.MessagingSystemKey.String("nats"),
			semconv.MessagingDestinationNameKey.String(p.channel),
			attribute.String("nats.retention_period", p.cfg.MaxAge.String()),
		),
	)
	defer span.End()

	if _, err := p.js.CreateOrUpdateStream(ctx, p.cfg); err != nil {
		span.SetStatus(codes.Error, "failed to create or update channel")
		span.RecordError(err)
		return fmt.Errorf("failed to create or update channel: %w", err)
	}

	span.SetStatus(codes.Ok, fmt.Sprintf("channel %s created or updated successfully", StreamName(p.channel)))

	return nil
}

func (p *NatsProducer) Publish(
	ctx context.Context,
	event core.EventSpec,
) error {
	tracer := otel.Tracer(p.name)

	ctx, span := tracer.Start(
		ctx, event.Type(),
		trace.WithSpanKind(trace.SpanKindProducer),
		trace.WithAttributes(
			semconv.MessagingSystemKey.String("nats"),
			semconv.MessagingDestinationNameKey.String(event.Type()),
			attribute.String("producer", p.name),
			attribute.String("channel", p.channel),
			attribute.String("event.id", event.ID()),
			attribute.String("event.type", event.Type()),
		),
	)
	defer span.End()
	data, err := proto.Marshal(event.Data())
	if err != nil {
		span.SetStatus(codes.Error, "failed to marshal event data")
		span.RecordError(err)
		return fmt.Errorf("failed to marshal event data: %w", err)
	}

	headers := nats.Header{}
	headers.Set(core.EVENT_ID_HEADER, event.ID())
	headers.Set("Event-Type", event.Type())
	headers.Set("Content-Type", "application/protobuf")

	otel.GetTextMapPropagator().Inject(ctx, &natsCarrier{headers})

	msg := &nats.Msg{
		Subject: event.Type(),
		Header:  headers,
		Data:    data,
	}

	if _, err := p.js.PublishMsg(ctx, msg); err != nil {
		span.SetStatus(codes.Error, "failed to publish event")
		span.RecordError(err)
		return fmt.Errorf("failed to publish event: %w", err)
	}
	span.SetStatus(codes.Ok, "event published successfully")

	return nil
}
