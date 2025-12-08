package bootstrap

import (
	"backend/internal/infrastructure/config"

	"go.uber.org/fx"
)

var ConfigModule = fx.Module(
	"config",
	fx.Provide(config.LoadAppConfig),
)
