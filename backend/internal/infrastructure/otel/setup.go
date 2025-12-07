package otel

import (
	"context"
	"errors"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/log/global"
)

func Setup(
	ctx context.Context,
) (func(context.Context) error, error) {
	var shutdownFuncs []func(context.Context) error
	var err error

	shutdown := func(ctx context.Context) error {
		var err error
		for _, fn := range shutdownFuncs {
			err = errors.Join(err, fn(ctx))
		}
		shutdownFuncs = nil
		return err
	}

	// Error handler
	errorHandler := func(inErr error) {
		err = errors.Join(inErr, shutdown(ctx))
	}

	// Propagation
	prop := Propagator()
	otel.SetTextMapPropagator(prop)

	// Tracer provider
	tracerProvider, err := TracerProvider()
	if err != nil {
		errorHandler(err)
		return nil, err
	}
	shutdownFuncs = append(shutdownFuncs, tracerProvider.Shutdown)
	otel.SetTracerProvider(tracerProvider)

	// Meter provider
	meterProvider, err := MeterProvider()
	if err != nil {
		errorHandler(err)
		return nil, err
	}
	shutdownFuncs = append(shutdownFuncs, meterProvider.Shutdown)
	otel.SetMeterProvider(meterProvider)

	// Logger provider
	loggerProvider, err := LoggerProvider()
	if err != nil {
		errorHandler(err)
		return nil, err
	}
	shutdownFuncs = append(shutdownFuncs, loggerProvider.Shutdown)
	global.SetLoggerProvider(loggerProvider)

	SetupHttpClient()

	return shutdown, nil
}
