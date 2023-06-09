package config

import (
	"fmt"

	"github.com/caarlos0/env/v8"
)

type Config struct {
	TelegramBotToken string `env:"TG_BOT_TOKEN"`
	BaseURL          string `env:"BASE_URL" envDefault:"http://server:10000"`
}

func NewConfig() (*Config, error) {
	var c Config
	if err := env.Parse(&c); err != nil {
		return nil, fmt.Errorf("parse config: %w", err)
	}

	return &c, nil
}
