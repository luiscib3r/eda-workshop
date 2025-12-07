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
	fx.Provide(storage.NewFileUploadedConsumer),
	fx.Provide(storage.NewFilesDeletedConsumer),
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
	fileUploadedConsumer *storage.FileUploadedConsumer,
	filesDeletedConsumer *storage.FilesDeletedConsumer,
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
