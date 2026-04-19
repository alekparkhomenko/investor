package config

import (
	"time"

	"github.com/alekparkhomenko/investor/platform/pkg/logger"
)

type AppSettings interface {
	Symbols() string
	PollInterval() time.Duration
}

type LoggerSettings interface {
	Level() string
	AsJson() bool
	LokiHost() string
	LokiEnv() string
	LokiEnabled() bool
	ToPlatformLoggerConfig() logger.Config
}
