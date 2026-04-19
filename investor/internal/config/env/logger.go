package env

import (
	"os"

	"github.com/caarlos0/env/v11"
)

type loggerEnvConfig struct {
	Level       string `env:"LOG_LEVEL"`
	AsJson      bool   `env:"LOG_AS_JSON"`
	LokiHost    string `env:"LOKI_HOST"`
	LokiEnv     string `env:"LOKI_ENV"`
	LokiEnabled bool   `env:"LOKI_ENABLED"`
}

type loggerConfig struct {
	raw loggerEnvConfig
}

func NewLoggerConfig() (*loggerConfig, error) {
	var raw loggerEnvConfig
	if err := env.Parse(&raw); err != nil {
		return nil, err
	}

	return &loggerConfig{raw: raw}, nil
}

func (c *loggerConfig) Level() string {
	return c.raw.Level
}

func (c *loggerConfig) AsJson() bool {
	return c.raw.AsJson
}

func (c *loggerConfig) LokiHost() string {
	if c.raw.LokiHost == "" {
		return "http://localhost:3100"
	}
	return c.raw.LokiHost
}

func (c *loggerConfig) LokiEnv() string {
	if c.raw.LokiEnv == "" {
		return "development"
	}
	return c.raw.LokiEnv
}

func (c *loggerConfig) LokiEnabled() bool {
	// Default to false if not set
	if _, ok := os.LookupEnv("LOKI_ENABLED"); !ok {
		return false
	}
	return c.raw.LokiEnabled
}
