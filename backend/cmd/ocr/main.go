package main

import (
	"backend/cmd/ocr/bootstrap"

	"go.uber.org/fx"
)

func main() {
	app := fx.New(
		bootstrap.ConfigModule,
		bootstrap.OtelModule,
		bootstrap.StorageModule,
		bootstrap.PostgresModule,
		bootstrap.NatsModule,
		bootstrap.LlmModule,
	)

	app.Run()
}
