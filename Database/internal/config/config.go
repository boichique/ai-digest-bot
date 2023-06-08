package config

import (
	"fmt"

	"github.com/caarlos0/env/v8"
)

type Config struct {
	DBUrl           string `env:"DB_URL"`
	Port            int    `env:"PORT"`
	Local           bool   `env:"LOCAL" envDefault:"false"`
	LogLevel        string `env:"LOG_LEVEL" envDefault:"info"`
	YoutubeApiToken string `env:"YOUTUBE_API_TOKEN"`
	ChatGPTApiToken string `env:"CHAT_GPT_API_TOKEN"`
}

func NewConfig() (*Config, error) {
	var c Config
	if err := env.Parse(&c); err != nil {
		return nil, fmt.Errorf("parse config: %w", err)
	}

	return &c, nil
}
