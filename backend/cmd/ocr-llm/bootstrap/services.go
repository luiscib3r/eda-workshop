package bootstrap

import (
	"backend/internal/health"
	"backend/internal/infrastructure/service"
	ocrllm "backend/internal/ocr-llm"

	"go.uber.org/fx"
)

var ServicesModule = fx.Module(
	"services",
	// Provide services
	fx.Provide(service.AsService(health.NewHealthService)),
	fx.Provide(service.AsService(ocrllm.NewLlmDebugService)),
	// Register services
	fx.Provide(service.AsRegister(service.RegisterServices)),
)
