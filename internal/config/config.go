package config

import (
	"fmt"
	"time"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	ServerConfig        ServerConfig
	PostgresConfig      PostgresConfig
	RedisConfig         RedisConfig
	JwtConfig           JwtConfig
	LoggerConfig        LoggerConfig
	CurrencyRatesConfig CurrencyRatesConfig
}

type ServerConfig struct {
	HTTPHost string `envconfig:"HTTP_HOST" required:"true"`
	HTTPPort string `envconfig:"HTTP_PORT" required:"true"`
}

type LoggerConfig struct {
	LogLevel string `envconfig:"LOG_LEVEL" default:"debug"`
}

type PostgresConfig struct {
	ConnDb       string `envconfig:"CONN_DB" required:"true"`
	PostgresHost string `envconfig:"POSTGRES_HOST" required:"true"`
	PostgresPort string `envconfig:"POSTGRES_PORT" required:"true"`
	PostgresUser string `envconfig:"POSTGRES_USER" required:"true"`
	PostgresPass string `envconfig:"POSTGRES_PASS" required:"true"`
	PostgresDb   string `envconfig:"POSTGRES_DB" required:"true"`
}

type RedisConfig struct {
	RedisHost string `envconfig:"REDIS_HOST" required:"true"`
	RedisPort string `envconfig:"REDIS_PORT" required:"true"`
}

type JwtConfig struct {
	JwtSecret     string        `envconfig:"JWT_SECRET" default:"secret-key"`
	JwtAccessTTL  time.Duration `envconfig:"JWT_ACCESS_TTL" default:"15m"`
	JwtRefreshTTL time.Duration `envconfig:"JWT_REFRESH_TTL" default:"720h"`
}

type CurrencyRatesConfig struct {
	Host    string        `env:"CURRENCY_RATES_HOST" default:"https://api.frankfurter.dev/v2"`
	Timeout time.Duration `env:"	CURRENCY_RATES_CLIENT_TIMEOUT" default:"10s"`
}

func Load() (Config, error) {
	if err := godotenv.Load(); err != nil {
		return Config{}, fmt.Errorf("load env: %w", err)
	}

	var cfg Config
	if err := envconfig.Process("", &cfg); err != nil {
		return Config{}, fmt.Errorf("read Config: %w", err)
	}

	return cfg, nil
}
