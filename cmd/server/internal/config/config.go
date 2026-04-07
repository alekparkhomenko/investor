package config

import (
	"os"
)

type Config struct {
	FinnhubAPIKey string
	TelegramToken string
}

func Load() *Config {
	return &Config{
		FinnhubAPIKey: getEnv("FINNHUB_API_KEY", ""),
		TelegramToken: getEnv("TELEGRAM_TOKEN", ""),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
