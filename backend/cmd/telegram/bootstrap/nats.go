package bootstrap

import (
	"backend/internal/infrastructure/nats"
	"backend/internal/telegram"
	"context"

	"go.uber.org/fx"
)

var NatsModule = fx.Module(
	"nats",
	fx.Provide(nats.NewNatsClient),
	fx.Provide(nats.NewJetStreamClient),
	fx.Provide(telegram.NewFileUploadedConsumer),
	fx.Invoke(SubcribeTelegramConsumers),
)

func SubcribeTelegramConsumers(
	lc fx.Lifecycle,
	fileUploadedConsumer *telegram.FileUploadedConsumer,
) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			return fileUploadedConsumer.Subscribe(ctx)
		},
		OnStop: func(ctx context.Context) error {
			fileUploadedConsumer.Stop()
			return nil
		},
	})
}
