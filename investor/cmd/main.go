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
	"github.com/alekparkhomenko/investor/investor/internal/ingestor"
	"github.com/alekparkhomenko/investor/platform/pkg/closer"
	"github.com/alekparkhomenko/investor/platform/pkg/logger"
	"go.uber.org/zap"
)

func main() {
	err := config.Load()
	if err != nil {
		panic(fmt.Errorf("failed to load config: %w", err))
	}

	cfg := config.AppConfig()

	if err := logger.Init(cfg.Logger.Level(), cfg.Logger.AsJson()); err != nil {
		println("logger init error:", err.Error())
		os.Exit(1)
	}
	defer logger.Sync()

	log := logger.With(zap.String("component", "main"))

	ing := ingestor.NewMOEXIngestor(cfg.App.Symbols())
	a := app.NewApp(cfg, ing)

	closer.AddNamed("app", func(ctx context.Context) error {
		return a.Stop()
	})

	appCtx, appCancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer appCancel()
	defer gracefulShutdown()

	closer.Configure(syscall.SIGINT, syscall.SIGTERM)
	closer.SetLogger(log)

	a.Run(appCtx)
}

func gracefulShutdown() {
	log := logger.With(zap.String("component", "closer"))

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := closer.CloseAll(ctx); err != nil {
		log.Error(ctx, "error during shutdown", zap.Error(err))
	}
}
