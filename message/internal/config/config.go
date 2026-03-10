package config

import (
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type Config struct {
	Host             string `env:"HOST"`
	Port             string `env:"PORT"`
	PostgresHost     string `env:"POSTGRES_HOST"`
	PostgresPort     string `env:"POSTGRES_PORT"`
	PostgresUser     string `env:"POSTGRES_USER"`
	PostgresPassword string `env:"POSTGRES_PASSWORD"`
	PostgresDatabase string `env:"POSTGRES_DATABASE"`
	KafkaHost        string `env:"KAFKA_HOST"`
	KafkaPort        string `env:"KAFKA_PORT"`
}

func NewConfig() (*Config, error) {
	var cfg Config
	godotenv.Load(".env")
	if err := cleanenv.ReadEnv(&cfg); err != nil {
		return nil, fmt.Errorf("failed to read config: %w", err)
	}
	return &cfg, nil
}
