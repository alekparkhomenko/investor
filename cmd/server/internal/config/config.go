package config

import (
	"os"
)

type Config struct {
	TelegramToken string
}

func Load() *Config {
	return &Config{
		TelegramToken: getEnv("TELEGRAM_TOKEN", ""),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
