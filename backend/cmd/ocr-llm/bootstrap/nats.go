package bootstrap

import (
	"backend/internal/infrastructure/nats"
	ocrllm "backend/internal/ocr-llm"
	"context"

	"go.uber.org/fx"
)

var NatsModule = fx.Module(
	"nats",
	fx.Provide(nats.NewNatsClient),
	fx.Provide(nats.NewJetStreamClient),
	fx.Provide(ocrllm.NewFilePageRegisteredConsumer),
	fx.Invoke(SubcribeOcrLlmConsumers),
)

func SubcribeOcrLlmConsumers(
	lc fx.Lifecycle,
	filePageRegisteredConsumer *ocrllm.FilePageRegisteredConsumer,
) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			if err := filePageRegisteredConsumer.Subscribe(ctx); err != nil {
				return err
			}

			return nil
		},
		OnStop: func(ctx context.Context) error {
			filePageRegisteredConsumer.Stop()
			return nil
		},
	})
}
