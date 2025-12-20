package ocr

import (
	"backend/cmd/ocr/bootstrap"

	"github.com/spf13/cobra"
	"go.uber.org/fx"
)

var OcrCmd = &cobra.Command{
	Use:   "ocr",
	Short: "OCR Service",
	Run:   run,
}

func run(cmd *cobra.Command, args []string) {
	app := fx.New(
		bootstrap.ConfigModule,
		bootstrap.OtelModule,
		bootstrap.StorageModule,
		bootstrap.PostgresModule,
		bootstrap.NatsModule,
		bootstrap.ServicesModule,
		bootstrap.ServerModule,
	)

	app.Run()
}
