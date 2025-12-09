package bootstrap

import (
	"backend/internal/infrastructure/otel"
	"context"

	"go.uber.org/fx"
)

var OtelModule = fx.Module(
	"otel",
	fx.Invoke(RunOtel),
)

func RunOtel(
	lc fx.Lifecycle,
) error {
	shutdown, err := otel.Setup(context.Background())
	if err != nil {
		return err
	}

	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			err := shutdown(ctx)
			return err
		},
	})

	return nil
}
