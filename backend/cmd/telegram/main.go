package main

import (
	"backend/cmd/telegram/bootstrap"

	"go.uber.org/fx"
)

func main() {
	app := fx.New(
		fx.StopTimeout(0),
		bootstrap.ConfigModule,
		bootstrap.OtelModule,
		bootstrap.TelegramModule,
		bootstrap.NatsModule,
		bootstrap.ServicesModule,
		bootstrap.ServerModule,
	)

	app.Run()
}
