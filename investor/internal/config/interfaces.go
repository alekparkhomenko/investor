package config

import "time"

type AppSettings interface {
	Symbols() string
	PollInterval() time.Duration
}

type LoggerSettings interface {
	Level() string
	AsJson() bool
}
