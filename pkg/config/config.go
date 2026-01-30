package config

import (
	"fmt"
	"os"
	"time"
)

type Config struct {
	AppEnv      string
	DB_DSN      string
	RedisAddr   string
	JWTSecret   string
	JWTExpiry   time.Duration
}

func Load() (*Config, error) {
	cfg := &Config{
		AppEnv:    getEnv("APP_ENV", "development"),
		DB_DSN:    os.Getenv("DB_DSN"),
		RedisAddr: os.Getenv("REDIS_ADDR"),
		JWTSecret: os.Getenv("JWT_SECRET"),
		JWTExpiry: 72 * time.Hour,
	}

	if cfg.DB_DSN == "" {
		return nil, fmt.Errorf("DB_DSN is required")
	}
	if cfg.RedisAddr == "" {
		return nil, fmt.Errorf("REDIS_ADDR is required")
	}
	if cfg.JWTSecret == "" {
		return nil, fmt.Errorf("JWT_SECRET is required")
	}

	return cfg, nil
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}