package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/alekparkhomenko/investor/investor/internal/app"
	"github.com/alekparkhomenko/investor/investor/internal/config"
	"github.com/alekparkhomenko/investor/investor/internal/http"
	"github.com/alekparkhomenko/investor/investor/internal/ingestor"
	"github.com/alekparkhomenko/investor/investor/internal/storage"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/alekparkhomenko/investor/platform/pkg/closer"
	"github.com/alekparkhomenko/investor/platform/pkg/logger"
)

func main() {
	// Load configuration
	if err := config.Load(); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "failed to load config: %v\n", err)
		os.Exit(1)
	}

	cfg := config.AppConfig()

	// Initialize logger from platform
	loggerCfg := cfg.Logger.ToPlatformLoggerConfig()
	appLogger, err := logger.New(&loggerCfg)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "failed to initialize logger: %v\n", err)
		os.Exit(1)
	}

	// Register logger in closer for graceful shutdown
	closer.AddNamed("logger", func(ctx context.Context) error {
		return appLogger.Close()
	})

	appCtx, appCancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer appCancel()

	// Log startup
	appLogger.Info(appCtx, "starting investor", logger.Fields{
		"component": "main",
		"symbols":   cfg.App.Symbols(),
	})

	// Initialize PostgreSQL connection pool
	pool, err := pgxpool.New(appCtx, cfg.DB.URL)
	if err != nil {
		appLogger.Error(appCtx, "failed to connect to database", logger.Fields{
			"component": "main",
			"error":     err.Error(),
		})
		os.Exit(1)
	}
	defer pool.Close()

	// Configure connection pool
	pool.Config().MaxConns = int32(cfg.DB.MaxOpenConns)
	pool.Config().MinConns = int32(cfg.DB.MaxIdleConns)

	// Initialize portfolio storage
	store := storage.NewPortfolioStore(pool)

	// Initialize HTTP handler
	handler := http.NewHandler(store)

	// Initialize HTTP server
	httpServer := http.NewHTTPServer(http.Config{
		Address: cfg.App.HTTPAddress(),
		Handler: handler,
		Logger:  appLogger,
	})

	// Register HTTP server in closer
	closer.AddNamed("http", httpServer.Stop)
	closer.AddNamed("db", func(ctx context.Context) error {
		pool.Close()
		return nil
	})

	// Start HTTP server in goroutine
	go func() {
		if err := httpServer.Start(); err != nil {
			appLogger.Error(appCtx, "HTTP server error", logger.Fields{
				"component": "main",
				"error":     err.Error(),
			})
		}
	}()

	// Initialize components with logger injection
	ing := ingestor.NewMOEXIngestor(cfg.App.Symbols(), appLogger)
	a := app.NewApp(cfg, ing, appLogger)

	// Register app in closer
	closer.AddNamed("app", func(ctx context.Context) error {
		return a.Stop()
	})

	closer.Configure(syscall.SIGINT, syscall.SIGTERM)

	// Run application
	if err := a.Run(appCtx); err != nil {
		appLogger.Error(appCtx, "application error", logger.Fields{
			"component": "main",
			"error":     err.Error(),
		})
		os.Exit(1)
	}

	// Graceful shutdown
	gracefulShutdown(appCtx, appLogger)
}

func gracefulShutdown(ctx context.Context, log *logger.Logger) {
	shutdownCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	log.Info(shutdownCtx, "graceful shutdown started", logger.Fields{
		"component": "main",
	})

	if err := closer.CloseAll(shutdownCtx); err != nil {
		log.Error(shutdownCtx, "error during shutdown", logger.Fields{
			"component": "main",
			"error":     err.Error(),
		})
	}

	log.Info(shutdownCtx, "shutdown complete", logger.Fields{
		"component": "main",
	})
}