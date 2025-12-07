package nats

import (
	"backend/internal/core"
	"context"
	"fmt"

	"github.com/nats-io/nats.go/jetstream"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
	"go.opentelemetry.io/otel/trace"
)

type NatsConsumer[T core.EventSpec] struct {
	name       string
	channel    string
	event      string
	numWorkers int
	builder    core.EventBuilder[T]
	handler    core.EventHandler[T]
	js         jetstream.JetStream
	cfg        jetstream.ConsumerConfig
	consumer   jetstream.ConsumeContext
	msgCh      chan jetstream.Msg
	stopCh     chan struct{}
	workerSem  chan struct{}
}

func NewNatsConsumer[T core.EventSpec](
	name string,
	channel string,
	event string,
	numWorkers int,
	workerBufferSize int,
	builder core.EventBuilder[T],
	handler core.EventHandler[T],
	js jetstream.JetStream,
	cfg jetstream.ConsumerConfig,
) *NatsConsumer[T] {
	if numWorkers <= 0 {
		numWorkers = 1
	}
	if workerBufferSize <= 0 {
		workerBufferSize = 10
	}

	return &NatsConsumer[T]{
		name:       name,
		channel:    channel,
		event:      event,
		js:         js,
		cfg:        cfg,
		numWorkers: numWorkers,
		builder:    builder,
		handler:    handler,
		msgCh:      make(chan jetstream.Msg, numWorkers*workerBufferSize),
		stopCh:     make(chan struct{}),
		workerSem:  make(chan struct{}, numWorkers),
	}
}

func (c *NatsConsumer[T]) Subscribe(
	ctx context.Context,
) error {
	tracer := otel.Tracer(c.name)

	ctx, span := tracer.Start(
		ctx, fmt.Sprintf("%s.subcribe", c.channel),
		trace.WithSpanKind(trace.SpanKindConsumer),
		trace.WithAttributes(
			semconv.MessagingSystemKey.String("nats"),
			semconv.MessagingDestinationNameKey.String(c.event),
			attribute.String("consumer", c.name),
			attribute.String("channel", c.channel),
		),
	)
	defer span.End()

	stream, err := c.js.Stream(ctx, StreamName(c.channel))
	if err != nil {
		span.SetStatus(codes.Error, "failed to get stream")
		span.RecordError(err)
		return fmt.Errorf("failed to get stream: %w", err)
	}

	consumer, err := stream.CreateOrUpdateConsumer(ctx, c.cfg)
	if err != nil {
		span.SetStatus(codes.Error, "failed to create or update consumer")
		span.RecordError(err)
		return fmt.Errorf("failed to create or update consumer: %w", err)
	}

	c.consumer, err = consumer.Consume(func(msg jetstream.Msg) {
		select {
		case c.msgCh <- msg:
			c.startWorker()
		case <-c.stopCh:
			return
		}
	})

	return nil
}

func (c *NatsConsumer[T]) startWorker() {
	select {
	case c.workerSem <- struct{}{}:
		go c.worker()
	default:
		// There are already max workers running
	}
}

func (c *NatsConsumer[T]) worker() {
	defer func() { <-c.workerSem }()

	for {
		select {
		case msg := <-c.msgCh:
			c.process(msg)
		default:
			// There are no more messages to process
			return
		}
	}
}

func (c *NatsConsumer[T]) process(msg jetstream.Msg) {
	tracer := otel.Tracer(c.name)
	ctx, span := tracer.Start(
		otel.GetTextMapPropagator().Extract(
			context.Background(), &natsCarrier{msg.Headers()},
		),
		c.event,
		trace.WithSpanKind(trace.SpanKindConsumer),
		trace.WithAttributes(
			semconv.MessagingSystemKey.String("nats"),
			semconv.MessagingDestinationNameKey.String(c.event),
			attribute.String("consumer", c.name),
			attribute.String("channel", c.channel),
		),
	)
	defer span.End()

	event, err := c.builder(msg)
	if err != nil {
		span.SetStatus(codes.Error, "failed to build event from message")
		span.RecordError(err)
		if err := msg.Nak(); err != nil {
			span.RecordError(err)
		}
		return
	}

	span.SetAttributes(
		attribute.String("event.id", event.ID()),
		attribute.String("event.type", event.Type()),
	)

	if err := c.handler(ctx, event); err != nil {
		span.SetStatus(codes.Error, "failed to handle event")
		span.RecordError(err)
		if err := msg.Nak(); err != nil {
			span.RecordError(err)
		}
		return
	}

	if err := msg.Ack(); err != nil {
		span.SetStatus(codes.Error, "failed to acknowledge message")
		span.RecordError(err)
		return
	}

	span.SetStatus(codes.Ok, "event processed successfully")
}

func (c *NatsConsumer[T]) Stop() {
	if c.consumer != nil {
		c.consumer.Drain()
		c.consumer.Stop()
	}
	close(c.stopCh)
}
