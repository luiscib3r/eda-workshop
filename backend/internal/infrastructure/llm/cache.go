package llm

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"time"

	"github.com/nats-io/nats.go/jetstream"
	"github.com/openai/openai-go/v3"
)

type CacheKey struct {
	Model    string
	Messages []openai.ChatCompletionMessageParamUnion
}

func GenerateCacheKey(params CacheKey) string {
	data, _ := json.Marshal(params)
	hash := sha256.Sum256(data)
	return fmt.Sprintf("llm.%x", hash[:16])
}

func NewLlmCache(
	js jetstream.JetStream,
) (jetstream.KeyValue, error) {
	ctx := context.Background()
	kv, err := js.CreateOrUpdateKeyValue(ctx, jetstream.KeyValueConfig{
		Bucket:  "LLM_CACHE",
		History: 1,
		TTL:     24 * 7 * time.Hour,
	})

	if err != nil {
		return nil, err
	}

	return kv, nil
}
