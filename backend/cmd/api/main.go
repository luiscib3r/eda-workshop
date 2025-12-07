package main

import (
	"backend/cmd/api/bootstrap"

	"go.uber.org/fx"
)

func main() {
	app := fx.New(
		bootstrap.ConfigModule,
		bootstrap.OtelModule,
		bootstrap.PostgresModule,
		bootstrap.StorageModule,
		bootstrap.NatsModule,
		bootstrap.ServicesModule,
		bootstrap.ServerModule,
	)

	app.Run()
}
