package bootstrap

import (
	"backend/internal/infrastructure/llm"
	ocrllm "backend/internal/ocr-llm"

	"go.uber.org/fx"
)

var LlmModule = fx.Module(
	"llm",
	fx.Provide(llm.NewClient),
	fx.Provide(ocrllm.NewOcrAgent),
)
