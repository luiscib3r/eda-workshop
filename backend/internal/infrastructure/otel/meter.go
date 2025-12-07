package otel

import (
	"go.opentelemetry.io/otel/sdk/metric"
)

func MeterProvider() (*metric.MeterProvider, error) {
	exporter, err := MetricExporter()
	if err != nil {
		return nil, err
	}

	res, err := Resource()
	if err != nil {
		return nil, err
	}

	provider := metric.NewMeterProvider(
		metric.WithReader(metric.NewPeriodicReader(exporter)),
		metric.WithResource(res),
	)

	return provider, nil
}
