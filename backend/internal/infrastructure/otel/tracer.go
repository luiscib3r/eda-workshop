package otel

import (
	"go.opentelemetry.io/otel/sdk/trace"
)

func TracerProvider() (*trace.TracerProvider, error) {
	exporter, err := TraceExporter()
	if err != nil {
		return nil, err
	}

	res, err := Resource()
	if err != nil {
		return nil, err
	}

	provider := trace.NewTracerProvider(
		trace.WithBatcher(exporter),
		trace.WithResource(res),
	)

	return provider, nil
}
