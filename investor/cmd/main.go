package main

import (
	"context"
	"fmt"
	"os/signal"
	"syscall"
	"time"

	"github.com/alekparkhomenko/investor/investor/internal/app"
	"github.com/alekparkhomenko/investor/investor/internal/config"
	"github.com/alekparkhomenko/investor/investor/internal/ingestor"
	"github.com/alekparkhomenko/investor/platform/pkg/closer"
)

func main() {
	err := config.Load()
	if err != nil {
		panic(fmt.Errorf("failed to load config: %w", err))
	}

	cfg := config.AppConfig()

	fmt.Println("Starting investor with symbols:", cfg.App.Symbols())

	ing := ingestor.NewMOEXIngestor(cfg.App.Symbols())
	a := app.NewApp(cfg, ing)

	closer.AddNamed("app", func(ctx context.Context) error {
		return a.Stop()
	})

	appCtx, appCancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer appCancel()
	defer gracefulShutdown()

	closer.Configure(syscall.SIGINT, syscall.SIGTERM)

	a.Run(appCtx)
}

func gracefulShutdown() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	fmt.Println("Graceful shutdown...")
	if err := closer.CloseAll(ctx); err != nil {
		fmt.Println("Error during shutdown:", err)
	}
	fmt.Println("Shutdown complete")
}
