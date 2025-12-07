package otel

import (
	"log/slog"
	"os"

	"go.opentelemetry.io/contrib/bridges/otelslog"
	"go.opentelemetry.io/otel/sdk/log"
)

func LoggerProvider() (*log.LoggerProvider, error) {
	exporter, err := LoggerExporter()
	if err != nil {
		return nil, err
	}

	res, err := Resource()
	if err != nil {
		return nil, err
	}

	provider := log.NewLoggerProvider(
		log.WithProcessor(log.NewBatchProcessor(exporter)),
		log.WithResource(res),
	)

	appName := os.Getenv("OTEL_SERVICE_NAME")
	logger := otelslog.NewLogger(
		appName,
		otelslog.WithLoggerProvider(provider),
	)
	slog.SetDefault(logger)

	return provider, nil
}
