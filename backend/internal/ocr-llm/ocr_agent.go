package ocrllm

import (
	"backend/internal/infrastructure/llm"
	"context"
	"encoding/base64"
	"fmt"

	"github.com/openai/openai-go/v3"
)

type OcrAgent struct {
	cfg *llm.AgentConfig
	api *openai.Client
}

func NewOcrAgent(
	cfg *llm.LlmConfig,
	api *openai.Client,
) *OcrAgent {
	return &OcrAgent{
		cfg: &cfg.Ocr,
		api: api,
	}
}

func (a *OcrAgent) Invoke(ctx context.Context, input []byte) (*openai.ChatCompletion, error) {
	image := base64.StdEncoding.EncodeToString(input)

	params := openai.ChatCompletionNewParams{
		Model: a.cfg.Model,
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.SystemMessage(a.cfg.System),
			openai.UserMessage([]openai.ChatCompletionContentPartUnionParam{
				{
					OfText: &openai.ChatCompletionContentPartTextParam{
						Text: a.cfg.User,
					},
				},
				{
					OfImageURL: &openai.ChatCompletionContentPartImageParam{
						ImageURL: openai.ChatCompletionContentPartImageImageURLParam{
							URL: fmt.Sprintf("data:image/png;base64,%s", image),
						},
					},
				},
			}),
		},
	}

	params.SetExtraFields(map[string]any{
		"provider": map[string]any{
			"order":           a.cfg.Providers,
			"allow_fallbacks": false,
		},
	})

	response, err := a.api.Chat.Completions.New(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("error creating chat completion: %w", err)
	}

	return response, nil
}
