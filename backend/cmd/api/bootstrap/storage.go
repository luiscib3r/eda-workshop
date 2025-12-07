package bootstrap

import (
	"backend/internal/infrastructure/storage"

	"go.uber.org/fx"
)

var StorageModule = fx.Module(
	"storage",
	fx.Provide(storage.NewClient),
	fx.Provide(storage.NewPresignClient),
)
