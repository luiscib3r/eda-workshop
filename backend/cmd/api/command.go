package api

import (
	"backend/cmd/api/bootstrap"

	"github.com/spf13/cobra"
	"go.uber.org/fx"
)

var ApiCmd = &cobra.Command{
	Use:   "api",
	Short: "API Service",
	Run:   run,
}

func run(cmd *cobra.Command, args []string) {
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
