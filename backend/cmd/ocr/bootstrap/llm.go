package bootstrap

import (
	"backend/internal/infrastructure/llm"

	"go.uber.org/fx"
)

var LlmModule = fx.Module(
	"llm",
	fx.Provide(llm.NewClient),
)
