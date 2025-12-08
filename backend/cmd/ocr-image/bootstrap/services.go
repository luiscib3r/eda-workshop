package bootstrap

import (
	"backend/internal/health"
	"backend/internal/infrastructure/service"

	"go.uber.org/fx"
)

var ServicesModule = fx.Module(
	"services",
	// Provide services
	fx.Provide(service.AsService(health.NewHealthService)),
	// Register services
	fx.Provide(service.AsRegister(service.RegisterServices)),
)
