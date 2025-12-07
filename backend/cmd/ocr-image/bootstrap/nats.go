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
	fx.Provide(ocr.NewFileUploadedConsumer),
	fx.Invoke(CreateOcrChannel),
	fx.Invoke(SubscribeOcrConsumers),
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

func SubscribeOcrConsumers(
	lc fx.Lifecycle,
	fileUploadedConsumer *ocr.FileUploadedConsumer,
) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			if err := fileUploadedConsumer.Subscribe(ctx); err != nil {
				return err
			}
			return nil
		},
		OnStop: func(ctx context.Context) error {
			fileUploadedConsumer.Stop()
			return nil
		},
	})
}
