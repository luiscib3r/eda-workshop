package service

import (
	"context"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"go.uber.org/fx"
)

type Service interface {
	Register(ctx context.Context, mux *runtime.ServeMux)
}

type ServicesDone struct{}

func AsService(f any) any {
	return fx.Annotate(
		f,
		fx.As(new(Service)),
		fx.ResultTags(`group:"services"`),
	)
}

func AsRegister(
	f any,
) any {
	return fx.Annotate(
		f,
		fx.ParamTags(`group:"services"`),
	)
}

func RegisterServices(
	services []Service,
	mux *runtime.ServeMux,
) *ServicesDone {
	ctx := context.Background()

	// Register all services
	for _, service := range services {
		service.Register(ctx, mux)
	}

	// Return done signal
	return &ServicesDone{}
}
