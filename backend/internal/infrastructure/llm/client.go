package llm

import (
	"backend/internal/infrastructure/config"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
)

func NewClient(
	cfg *config.AppConfig,
) *openai.Client {
	client := openai.NewClient(
		option.WithBaseURL(cfg.LLM.BaseUrl),
		option.WithAPIKey(cfg.LLM.ApiKey),
	)

	return &client
}
