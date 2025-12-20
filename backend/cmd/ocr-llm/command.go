package ocrllm

import (
	"backend/cmd/ocr-llm/bootstrap"

	"github.com/spf13/cobra"
	"go.uber.org/fx"
)

var OcrLlmCmd = &cobra.Command{
	Use:   "ocr-llm",
	Short: "OCR LLM Service",
	Run:   run,
}

func run(cmd *cobra.Command, args []string) {
	app := fx.New(
		bootstrap.ConfigModule,
		bootstrap.OtelModule,
		bootstrap.NatsModule,
		bootstrap.StorageModule,
		bootstrap.PostgresModule,
		bootstrap.LlmModule,
		bootstrap.ServicesModule,
		bootstrap.ServerModule,
	)

	app.Run()
}
