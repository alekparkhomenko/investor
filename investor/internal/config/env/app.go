package env

import (
	"time"

	"github.com/caarlos0/env/v11"
)

type appEnvConfig struct {
	Symbols      string `env:"SYMBOLS,required"`
	PollInterval string `env:"POLL_INTERVAL,required"`
}

type appConfig struct {
	raw appEnvConfig
}

func NewAppConfig() (*appConfig, error) {
	var raw appEnvConfig
	if err := env.Parse(&raw); err != nil {
		return nil, err
	}

	return &appConfig{raw: raw}, nil
}

func (c *appConfig) Symbols() string {
	return c.raw.Symbols
}

func (c *appConfig) PollInterval() time.Duration {
	d, err := time.ParseDuration(c.raw.PollInterval)
	if err != nil {
		return 2 * time.Second
	}
	return d
}
