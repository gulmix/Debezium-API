package config

import (
	"debezium_server/pkg/postgres"
	"fmt"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Environment string `env:"ENV" env-default:"development"`

	Port    int           `env:"PORT"         env-default:"8080"`
	Timeout time.Duration `env:"HTTP_TIMEOUT" env-default:"30s"`

	DebeziumBaseURL string `env:"DEBEZIUM_BASE_URL" env-default:"http://localhost:8080"`

	postgres.Config
}

func ParseConfigFromEnv() (*Config, error) {
	cfg := &Config{}

	if err := cleanenv.ReadEnv(cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config from env: %w", err)
	}

	return cfg, nil
}
