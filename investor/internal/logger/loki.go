package logger

import (
	"context"
	"log"
	"sync"

	lokilogger "github.com/edaniel30/loki-logger-go"

	"github.com/alekparkhomenko/investor/investor/internal/config"
)

var (
	globalLogger *Logger
	initOnce    sync.Once
)

// Logger is a unified logger interface that supports both Loki and Zap backends.
type Logger struct {
	loki   *lokilogger.Logger
	mu     sync.Mutex
	closed bool
}

// Fields represents structured log fields.
type Fields map[string]interface{}

// New creates a new logger based on configuration.
func New(cfg *config.Config) (*Logger, error) {
	var err error
	initOnce.Do(func() {
		err = initLogger(cfg)
	})
	return globalLogger, err
}

func initLogger(cfg *config.Config) error {
	loggerCfg := cfg.Logger

	if loggerCfg.LokiEnabled() {
		l, err := createLokiLogger(loggerCfg)
		if err != nil {
			return err
		}
		globalLogger = &Logger{loki: l}
		return nil
	}

	// Fallback to zap via plantform logger
	// For now, we just use a no-op logger wrapper
	globalLogger = &Logger{}
	return nil
}

func createLokiLogger(loggerCfg config.LoggerSettings) (*lokilogger.Logger, error) {
	opts := []lokilogger.Option{
		lokilogger.WithAppName("investor"),
		lokilogger.WithAppVersion("1.0.0"),
		lokilogger.WithAppEnv(loggerCfg.LokiEnv()),
		lokilogger.WithLokiHost(loggerCfg.LokiHost()),
	}

	client, err := lokilogger.New(lokilogger.DefaultConfig(), opts...)
	if err != nil {
		return nil, err
	}

	return client, nil
}

// Info logs an info message.
func (l *Logger) Info(ctx context.Context, msg string, fields Fields) {
	if l.loki != nil {
		l.loki.Info(ctx, msg, fields)
	}
}

// Error logs an error message.
func (l *Logger) Error(ctx context.Context, msg string, fields Fields) {
	if l.loki != nil {
		l.loki.Error(ctx, msg, fields)
	}
}

// Warn logs a warning message.
func (l *Logger) Warn(ctx context.Context, msg string, fields Fields) {
	if l.loki != nil {
		l.loki.Warn(ctx, msg, fields)
	}
}

// Debug logs a debug message.
func (l *Logger) Debug(ctx context.Context, msg string, fields Fields) {
	if l.loki != nil {
		l.loki.Debug(ctx, msg, fields)
	}
}

// Fatal logs a fatal message and exits.
func (l *Logger) Fatal(ctx context.Context, msg string, fields Fields) {
	if l.loki != nil {
		l.loki.Fatal(ctx, msg, fields)
	}
	log.Fatal(msg)
}

// Close flushes and closes the logger.
func (l *Logger) Close() error {
	l.mu.Lock()
	defer l.mu.Unlock()

	if l.closed {
		return nil
	}

	if l.loki != nil {
		l.loki.Close()
	}

	l.closed = true
	return nil
}
