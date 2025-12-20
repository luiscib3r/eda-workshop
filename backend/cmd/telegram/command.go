package telegram

import (
	"backend/cmd/telegram/bootstrap"

	"github.com/spf13/cobra"
	"go.uber.org/fx"
)

var TelegramCmd = &cobra.Command{
	Use:   "telegram",
	Short: "Telegram Bot Service",
	Run:   run,
}

func run(cmd *cobra.Command, args []string) {
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
