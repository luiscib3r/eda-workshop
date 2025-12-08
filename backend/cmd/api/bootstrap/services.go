package bootstrap

import (
	"backend/internal/health"
	"backend/internal/infrastructure/service"
	"backend/internal/ocr"
	"backend/internal/storage"
	"context"

	"go.uber.org/fx"
)

var ServicesModule = fx.Module(
	"services",
	// Provide services
	// Health
	fx.Provide(service.AsService(health.NewHealthService)),
	// Storage
	fx.Provide(storage.NewStorageService),
	fx.Provide(service.AsService(func(storage *storage.StorageService) *storage.StorageService {
		return storage
	})),
	fx.Provide(service.AsService(storage.NewFilesService)),
	// Ocr
	fx.Provide(service.AsService(ocr.NewFilesService)),
	// Register services
	fx.Provide(service.AsRegister(service.RegisterServices)),
	// Create buckets on startup
	fx.Invoke(CreateBuckets),
)

func CreateBuckets(
	lc fx.Lifecycle,
	storage *storage.StorageService,
) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			return storage.CreateBuckets(ctx)
		},
	})
}
