package main

import (
	"backend/cmd/api/bootstrap"

	"go.uber.org/fx"
)

func main() {
	app := fx.New(
		fx.StopTimeout(0),
		fx.NopLogger,
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
