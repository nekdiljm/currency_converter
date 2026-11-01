package config

import (
	"fmt"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	HTTPHost     string `envconfig:"HTTP_HOST" default:"localhost"`
	HTTPPort     string `envconfig:"HTTP_PORT" default:"8080"`
	ConnDb       string `envconfig:"CONN_DB" default:"postgres://currency_converter:123456@localhost:5432/currency_converter?sslmode=disable" `
	PostgresHost string `envconfig:"POSTGRES_HOST" default:"localhost"`
	PostgresPort string `envconfig:"POSTGRES_PORT" default:"5432"`
	PostgresUser string `envconfig:"POSTGRES_USER" default:"currency_converter"`
	PostgresPass string `envconfig:"POSTGRES_PASS" default:"123456"`
	PostgresDb   string `envconfig:"POSTGRES_DB" default:"currency_converter"`
}

func Load() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		return nil, fmt.Errorf("error load env")
	}

	var cfg Config
	if err := envconfig.Process("", &cfg); err != nil {
		return nil, fmt.Errorf("error to read Config")
	}

	return &cfg, nil
}
