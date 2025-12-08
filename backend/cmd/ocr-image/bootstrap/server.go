package bootstrap

import (
	"backend/internal/infrastructure/server"
	"backend/internal/infrastructure/service"
	"context"
	"fmt"
	"net/http"

	"go.uber.org/fx"
)

var ServerModule = fx.Module(
	"server",
	fx.Provide(server.NewGatewayServeMux),
	fx.Provide(server.NewServeMux),
	fx.Provide(server.NewHttpServer),
	fx.Invoke(RunServer),
)

func RunServer(
	lc fx.Lifecycle,
	srv *http.Server,
	_ *service.ServicesDone,
) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			fmt.Println("Starting http server on " + srv.Addr)
			go func() {
				if err := srv.ListenAndServe(); err != nil {
					fmt.Println("Failed to start server: " + err.Error())
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			fmt.Println("Stopping http server")
			return srv.Shutdown(ctx)
		},
	})
}
