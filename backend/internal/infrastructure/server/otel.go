package server

import (
	"context"
	"net/http"

	"go.opentelemetry.io/otel/trace"
	"google.golang.org/protobuf/proto"
)

func OtelTraceIDHeader(
	ctx context.Context,
	w http.ResponseWriter,
	_ proto.Message,
) error {
	span := trace.SpanFromContext(ctx)
	if sc := span.SpanContext(); sc.IsValid() {
		w.Header().Set("X-Trace-Id", sc.TraceID().String())
	}

	return nil
}
