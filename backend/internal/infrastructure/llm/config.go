package llm

import (
	"fmt"

	"github.com/spf13/viper"
)

type LlmConfig struct {
	Ocr AgentConfig `json:"ocr"`
}

type AgentConfig struct {
	Model     string   `json:"model"`
	Providers []string `json:"providers"`
	System    string   `json:"system"`
	User      string   `json:"user"`
}

func LoadLlmConfig() (*LlmConfig, error) {
	viper.SetConfigName("prompts")
	viper.SetConfigType("yaml")

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("error reading prompts config: %w", err)
	}

	var config LlmConfig
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("error unmarshaling prompts config: %w", err)
	}

	return &config, nil
}
