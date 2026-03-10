package config

import (
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type Config struct {
	ServicePort string `env:"SERVICE_PORT" required:"true"`
	ServiceHost string `env:"SERVICE_HOST" required:"true"`

	PostgresHost string `env:"POSTGRES_HOST" required:"true"`
	PostgresPort string `env:"POSTGRES_PORT" required:"true"`
	PostgresUser string `env:"POSTGRES_USER" required:"true"`
	PostgresPass string `env:"POSTGRES_PASSWORD" required:"true"`
	PostgresDB   string `env:"POSTGRES_DB_NAME" required:"true"`

	AppEnv string `env:"APP_ENV" required:"true"`
}

func NewConfig() (*Config, error) {
	var cfg Config
	godotenv.Load(".env")

	if err := cleanenv.ReadEnv(&cfg); err != nil {
		return nil, fmt.Errorf("unable to read config: %w", err)
	}
	return &cfg, nil
}
