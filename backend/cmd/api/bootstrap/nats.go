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
	fx.Provide(storage.NewOutboxProcessor),
	fx.Invoke(CreateStorageChannel),
	fx.Invoke(SubcribeStorageConsumers),
	fx.Invoke(RunOutboxProcessor),
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

func RunOutboxProcessor(
	lc fx.Lifecycle,
	processor *storage.OutboxProcessor,
) {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)

	lc.Append(fx.Hook{
		OnStart: func(_ context.Context) error {
			go processor.Start(ctx)
			return nil
		},
		OnStop: func(_ context.Context) error {
			cancel()
			return nil
		},
	})
}
