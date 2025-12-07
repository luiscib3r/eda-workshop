package bootstrap

import (
	"backend/internal/infrastructure/nats"
	"backend/internal/storage"
	"context"

	"go.uber.org/fx"
)

var NatsModule = fx.Module(
	"nats",
	fx.Provide(nats.NewNatsClient),
	fx.Provide(nats.NewJetStreamClient),
	fx.Provide(storage.NewStorageProducer),
	fx.Provide(storage.NewFileUploadedEventConsumer),
	fx.Invoke(CreateStorageChannel),
	fx.Invoke(SubcribeStorageConsumers),
)

func CreateStorageChannel(
	lc fx.Lifecycle,
	storage *storage.StorageProducer,
) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			return storage.CreateChannel(ctx)
		},
	})
}

func SubcribeStorageConsumers(
	lc fx.Lifecycle,
	fileUploadedConsumer *storage.FileUploadedEventConsumer,
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
