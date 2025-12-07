package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type AppConfig struct {
	Server   ServerConfig   `mapstructure:"server"`
	Storage  StorageConfig  `mapstructure:"storage"`
	Cors     CorsConfig     `mapstructure:"cors"`
	Nats     NatsConfig     `mapstructure:"nats"`
	Telegram TelegramConfig `mapstructure:"telegram"`
	Postgres PostgresConfig `mapstructure:"postgres"`
}

type ServerConfig struct {
	Port int `mapstructure:"port"`
}

type CorsConfig struct {
	AllowedOrigins []string `mapstructure:"allowed_origins"`
}

type StorageConfig struct {
	Endpoint       string `mapstructure:"endpoint"`
	PublicEndpoint string `mapstructure:"public_endpoint"`
	Region         string `mapstructure:"region"`
	AccessKey      string `mapstructure:"access_key"`
	SecretKey      string `mapstructure:"secret_key"`
	UsePathStyle   bool   `mapstructure:"use_path_style"`
}

type NatsConfig struct {
	Uri string `mapstructure:"uri"`
}

type TelegramConfig struct {
	BotToken string `mapstructure:"bot_token"`
	ChatID   int64  `mapstructure:"chat_id"`
}

type PostgresConfig struct {
	Uri string `mapstructure:"uri"`
	Dsn string `mapstructure:"dsn"`
}

func LoadAppConfig() (*AppConfig, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()
	viper.SetEnvPrefix("APP")

	viper.SetDefault("server.port", "8080")

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("error reading config: %w", err)
	}

	var config AppConfig
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("error unmarshaling config: %w", err)
	}

	return &config, nil
}
