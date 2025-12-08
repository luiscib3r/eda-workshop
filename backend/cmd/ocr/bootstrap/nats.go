package bootstrap

import (
	"backend/internal/infrastructure/nats"
	"backend/internal/ocr"
	"context"

	"go.uber.org/fx"
)

var NatsModule = fx.Module(
	"nats",
	fx.Provide(nats.NewNatsClient),
	fx.Provide(nats.NewJetStreamClient),
	fx.Provide(ocr.NewOcrProducer),
	fx.Provide(ocr.NewFilePageRenderedConsumer),
	fx.Provide(ocr.NewFilesDeletedConsumer),
	fx.Provide(ocr.NewOutboxProcessor),
	fx.Invoke(CreateOcrChannel),
	fx.Invoke(SubcribeOcrConsumers),
	fx.Invoke(RunOutboxProcessor),
)

func CreateOcrChannel(
	lc fx.Lifecycle,
	ocrProducer *ocr.OcrProducer,
) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			return ocrProducer.CreateChannel(ctx)
		},
	})
}

func SubcribeOcrConsumers(
	lc fx.Lifecycle,
	filePageRenderedConsumer *ocr.FilePageRenderedConsumer,
	filesDeletedConsumer *ocr.FilesDeletedConsumer,
) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			if err := filePageRenderedConsumer.Subscribe(ctx); err != nil {
				return err
			}
			if err := filesDeletedConsumer.Subscribe(ctx); err != nil {
				return err
			}
			return nil
		},
		OnStop: func(ctx context.Context) error {
			filePageRenderedConsumer.Stop()
			filesDeletedConsumer.Stop()
			return nil
		},
	})
}

func RunOutboxProcessor(
	lc fx.Lifecycle,
	processor *ocr.OutboxProcessor,
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
