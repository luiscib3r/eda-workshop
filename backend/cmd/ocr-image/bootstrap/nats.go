package bootstrap

import (
	"backend/internal/infrastructure/nats"
	"backend/internal/ocr"
	ocrimage "backend/internal/ocr-image"
	"context"

	"go.uber.org/fx"
)

var NatsModule = fx.Module(
	"nats",
	fx.Provide(nats.NewNatsClient),
	fx.Provide(nats.NewJetStreamClient),
	fx.Provide(ocr.NewOcrProducer),
	fx.Provide(ocrimage.NewFileUploadedConsumer),
	fx.Invoke(SubscribeOcrImageConsumers),
)

func SubscribeOcrImageConsumers(
	lc fx.Lifecycle,
	fileUploadedConsumer *ocrimage.FileUploadedConsumer,
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
