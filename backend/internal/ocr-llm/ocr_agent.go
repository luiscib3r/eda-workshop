package ocrllm

import (
	"backend/internal/infrastructure/llm"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"

	"github.com/nats-io/nats.go/jetstream"
	"github.com/openai/openai-go/v3"
)

type OcrAgent struct {
	cfg *llm.AgentConfig
	api *openai.Client
	kv  jetstream.KeyValue
}

func NewOcrAgent(
	cfg *llm.LlmConfig,
	api *openai.Client,
	kv jetstream.KeyValue,
) *OcrAgent {
	return &OcrAgent{
		cfg: &cfg.Ocr,
		api: api,
		kv:  kv,
	}
}

func (a *OcrAgent) Invoke(ctx context.Context, input []byte) (*openai.ChatCompletion, error) {
	image := base64.StdEncoding.EncodeToString(input)

	messages := []openai.ChatCompletionMessageParamUnion{
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
	}

	cacheKey := llm.GenerateCacheKey(llm.CacheKey{
		Model:    a.cfg.Model,
		Messages: messages,
	})

	entry, err := a.kv.Get(ctx, cacheKey)
	if err == nil {
		var cachedResponse openai.ChatCompletion
		if err := json.Unmarshal(entry.Value(), &cachedResponse); err == nil {
			return &cachedResponse, nil
		}
	}

	params := openai.ChatCompletionNewParams{
		Model:    a.cfg.Model,
		Messages: messages,
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

	data, _ := json.Marshal(response)
	a.kv.Put(ctx, cacheKey, data)

	return response, nil
}
