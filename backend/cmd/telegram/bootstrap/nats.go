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
	fx.Provide(telegram.NewFilesDeletedConsumer),
	fx.Invoke(SubcribeTelegramConsumers),
)

func SubcribeTelegramConsumers(
	lc fx.Lifecycle,
	fileUploadedConsumer *telegram.FileUploadedConsumer,
	filesDeletedConsumer *telegram.FilesDeletedConsumer,
) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			if err := fileUploadedConsumer.Subscribe(ctx); err != nil {
				return err
			}
			if err := filesDeletedConsumer.Subscribe(ctx); err != nil {
				return err
			}
			return nil
		},
		OnStop: func(ctx context.Context) error {
			fileUploadedConsumer.Stop()
			filesDeletedConsumer.Stop()
			return nil
		},
	})
}
