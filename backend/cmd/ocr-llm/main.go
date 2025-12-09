package main

import (
	"backend/cmd/ocr-llm/bootstrap"

	"go.uber.org/fx"
)

func main() {
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
