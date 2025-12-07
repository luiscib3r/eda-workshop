package otel

import (
	"context"

	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploggrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
)

func TraceExporter() (*otlptrace.Exporter, error) {
	exporter, err := otlptracegrpc.New(context.Background())
	if err != nil {
		return nil, err
	}

	return exporter, nil
}

func MetricExporter() (*otlpmetricgrpc.Exporter, error) {
	exporter, err := otlpmetricgrpc.New(context.Background())
	if err != nil {
		return nil, err
	}

	return exporter, nil
}

func LoggerExporter() (*otlploggrpc.Exporter, error) {
	exporter, err := otlploggrpc.New(context.Background())
	if err != nil {
		return nil, err
	}

	return exporter, nil
}
