package bootstrap

import (
	"backend/internal/health"
	"backend/internal/infrastructure/service"
	"backend/internal/storage"
	"context"

	"go.uber.org/fx"
)

var ServicesModule = fx.Module(
	"services",
	// Provide services
	fx.Provide(service.AsService(health.NewHealthService)),
	fx.Provide(storage.NewStorageService),
	fx.Provide(service.AsService(func(storage *storage.StorageService) *storage.StorageService {
		return storage
	})),
	fx.Provide(service.AsService(storage.NewFilesService)),
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
