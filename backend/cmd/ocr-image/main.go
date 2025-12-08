package main

import (
	"backend/cmd/ocr-image/bootstrap"

	"go.uber.org/fx"
)

func main() {
	app := fx.New(
		bootstrap.ConfigModule,
		bootstrap.OtelModule,
		bootstrap.StorageModule,
		bootstrap.NatsModule,
		bootstrap.ServicesModule,
		bootstrap.ServerModule,
	)

	app.Run()
}
