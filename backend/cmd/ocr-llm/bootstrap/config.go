package bootstrap

import (
	"backend/internal/infrastructure/config"
	"backend/internal/infrastructure/llm"

	"go.uber.org/fx"
)

var ConfigModule = fx.Module(
	"config",
	fx.Provide(config.LoadAppConfig),
	fx.Provide(llm.LoadLlmConfig),
)
