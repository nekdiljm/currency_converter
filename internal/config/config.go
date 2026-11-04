package config

import (
	"fmt"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	HTTPHost     string `envconfig:"HTTP_HOST" required:"true"`
	HTTPPort     string `envconfig:"HTTP_PORT" required:"true"`
	ConnDb       string `envconfig:"CONN_DB" required:"true"`
	PostgresHost string `envconfig:"POSTGRES_HOST" required:"true"`
	PostgresPort string `envconfig:"POSTGRES_PORT" required:"true"`
	PostgresUser string `envconfig:"POSTGRES_USER" required:"true"`
	PostgresPass string `envconfig:"POSTGRES_PASS" required:"true"`
	PostgresDb   string `envconfig:"POSTGRES_DB" required:"true"`
}

func Load() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		return nil, fmt.Errorf("load env: %w", err)
	}

	var cfg Config
	if err := envconfig.Process("", &cfg); err != nil {
		return nil, fmt.Errorf("read Config: %w", err)
	}

	return &cfg, nil
}
