package config

import (
	"fmt"

	"github.com/caarlos0/env/v11"
)

type Config struct {
	ServerPort      string `env:"SERVER_PORT" envDefault:"8080"`
	DatabaseURL     string `env:"DATABASE_URL" envDefault:"postgres://user:pass@localhost:5432/recdb?sslmode=disable"`
	RedisURL        string `env:"REDIS_URL" envDefault:"redis://localhost:6379/0"`
	RabbitMQURL     string `env:"RABBITMQ_URL" envDefault:"amqp://guest:guest@localhost:5672/"`
	JWTSecret       string `env:"JWT_SECRET" envDefault:"super-secret-key"`
	AccessTokenTTL  int    `env:"ACCESS_TOKEN_TTL" envDefault:"15"`     // 15 минут
	RefreshTokenTTL int    `env:"REFRESH_TOKEN_TTL" envDefault:"43200"` // 30 дней
}

func Load() (*Config, error) {
	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		return nil, fmt.Errorf("config load: %w", err)
	}
	return cfg, nil
}
