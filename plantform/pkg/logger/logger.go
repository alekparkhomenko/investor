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
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	lokilogger "github.com/edaniel30/loki-logger-go"
)

// Fields represents structured log fields for contextual logging.
type Fields map[string]interface{}

// Logger provides a unified logging interface with Loki and console output.
// Always logs to console (stdout/stderr), optionally also to Loki if enabled.
// When Loki is disabled, the logger still works and outputs to console only.
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
// Always logs to console, optionally also to Loki if enabled.
func New(cfg *Config) (*Logger, error) {
	logger := &Logger{}

	// Optionally create Loki logger
	if cfg.LokiEnabled {
		l, err := createLokiLogger(cfg)
		if err != nil {
			// Log error to console but continue with console-only logger
			consoleLog("ERROR", "failed to initialize Loki logger, using console-only", err)
			return logger, nil
		}
		logger.loki = l
	}

	return logger, nil
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

// consoleLog writes a log message to stdout/stderr with timestamp and level.
func consoleLog(level, msg string, err error) {
	timestamp := time.Now().Format(time.RFC3339)
	var sb strings.Builder

	sb.WriteString(timestamp)
	sb.WriteString(" ")
	sb.WriteString(level)
	sb.WriteString(" ")
	sb.WriteString(msg)

	if err != nil {
		sb.WriteString(" error=")
		sb.WriteString(err.Error())
	}

	sb.WriteString("\n")

	if level == "ERROR" || level == "FATAL" {
		fmt.Fprint(log.Writer(), sb.String())
	} else {
		fmt.Print(sb.String())
	}
}

// formatFields converts Fields to a readable string for console output.
func formatFields(fields Fields) string {
	if len(fields) == 0 {
		return ""
	}

	var sb strings.Builder
	sb.WriteString(" ")

	first := true
	for k, v := range fields {
		if !first {
			sb.WriteString(" ")
		}
		fmt.Fprintf(&sb, "%s=%v", k, v)
		first = false
	}

	return sb.String()
}

// Info logs an info message to console and Loki (if enabled).
func (l *Logger) Info(ctx context.Context, msg string, fields Fields) {
	// Always log to console
	consoleLog("INFO", msg+formatFields(fields), nil)

	// Also log to Loki if enabled
	if l.loki != nil {
		l.loki.Info(ctx, msg, fields)
	}
}

// Error logs an error message to console and Loki (if enabled).
func (l *Logger) Error(ctx context.Context, msg string, fields Fields) {
	// Always log to console
	consoleLog("ERROR", msg+formatFields(fields), nil)

	// Also log to Loki if enabled
	if l.loki != nil {
		l.loki.Error(ctx, msg, fields)
	}
}

// Warn logs a warning message to console and Loki (if enabled).
func (l *Logger) Warn(ctx context.Context, msg string, fields Fields) {
	// Always log to console
	consoleLog("WARN", msg+formatFields(fields), nil)

	// Also log to Loki if enabled
	if l.loki != nil {
		l.loki.Warn(ctx, msg, fields)
	}
}

// Debug logs a debug message to console and Loki (if enabled).
func (l *Logger) Debug(ctx context.Context, msg string, fields Fields) {
	// Always log to console
	consoleLog("DEBUG", msg+formatFields(fields), nil)

	// Also log to Loki if enabled
	if l.loki != nil {
		l.loki.Debug(ctx, msg, fields)
	}
}

// Fatal logs a fatal message to console and Loki (if enabled), then exits.
func (l *Logger) Fatal(ctx context.Context, msg string, fields Fields) {
	// Always log to console
	consoleLog("FATAL", msg+formatFields(fields), nil)

	// Also log to Loki if enabled
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
