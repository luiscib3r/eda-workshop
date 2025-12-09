package bootstrap

import (
	"backend/internal/infrastructure/nats"
	"backend/internal/telegram"
	"context"
	"fmt"

	"go.uber.org/fx"
)

var NatsModule = fx.Module(
	"nats",
	fx.Provide(nats.NewNatsClient),
	fx.Provide(nats.NewJetStreamClient),
	fx.Provide(telegram.NewFileUploadedConsumer),
	fx.Provide(telegram.NewFilesDeletedConsumer),
	fx.Provide(telegram.NewFilePageRenderedConsumer),
	fx.Provide(telegram.NewFilePageOcrGenerateConsumer),
	fx.Invoke(SubcribeTelegramConsumers),
)

func SubcribeTelegramConsumers(
	lc fx.Lifecycle,
	fileUploadedConsumer *telegram.FileUploadedConsumer,
	filesDeletedConsumer *telegram.FilesDeletedConsumer,
	filePageRenderedConsumer *telegram.FilePageRenderedConsumer,
	filePageOcrGenerateConsumer *telegram.FilePageOcrGenerateConsumer,
) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			if err := fileUploadedConsumer.Subscribe(ctx); err != nil {
				fmt.Println("Error subscribing to file uploaded consumer:", err)
			}
			if err := filesDeletedConsumer.Subscribe(ctx); err != nil {
				fmt.Println("Error subscribing to files deleted consumer:", err)
			}
			if err := filePageRenderedConsumer.Subscribe(ctx); err != nil {
				fmt.Println("Error subscribing to file page rendered consumer:", err)
			}
			if err := filePageOcrGenerateConsumer.Subscribe(ctx); err != nil {
				fmt.Println("Error subscribing to file page OCR generate consumer:", err)
			}
			return nil
		},
		OnStop: func(ctx context.Context) error {
			fileUploadedConsumer.Stop()
			filesDeletedConsumer.Stop()
			filePageRenderedConsumer.Stop()
			filePageOcrGenerateConsumer.Stop()
			return nil
		},
	})
}
