// Package logger provides a unified logging interface with Loki backend support.
//
// Usage example:
//
//	cfg := &logger.Config{
//		LokiEnabled: true,
//		LokiHost:    "http://localhost:3100",
//		LokiEnv:     "development",
//		AppName:     "my-app",
//		AppVersion:  "1.0.0",
//	}
//
//	log, err := logger.New(cfg)
//	if err != nil {
//		log.Fatal(err)
//	}
//	defer log.Close()
//
//	log.Info(ctx, "application started", logger.Fields{
//		"version": "1.0.0",
//		"env":     "development",
//	})
package logger

import (
	"context"
	"log"
	"sync"

	lokilogger "github.com/edaniel30/loki-logger-go"
)

// Fields represents structured log fields for contextual logging.
type Fields map[string]interface{}

// Logger provides a unified logging interface with Loki backend support.
// When Loki is disabled, it acts as a no-op logger.
type Logger struct {
	loki   *lokilogger.Logger
	mu     sync.Mutex
	closed bool
}

// Config holds logger configuration settings.
type Config struct {
	LokiEnabled bool
	LokiHost    string
	LokiEnv     string
	AppName     string
	AppVersion  string
}

// New creates a new logger instance based on the provided configuration.
// Returns a no-op logger if Loki is not enabled.
func New(cfg *Config) (*Logger, error) {
	if cfg.LokiEnabled {
		l, err := createLokiLogger(cfg)
		if err != nil {
			return nil, err
		}
		return &Logger{loki: l}, nil
	}

	// Return no-op logger if Loki is not enabled
	return &Logger{}, nil
}

func createLokiLogger(cfg *Config) (*lokilogger.Logger, error) {
	opts := []lokilogger.Option{
		lokilogger.WithAppName(cfg.AppName),
		lokilogger.WithAppVersion(cfg.AppVersion),
		lokilogger.WithAppEnv(cfg.LokiEnv),
		lokilogger.WithLokiHost(cfg.LokiHost),
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
