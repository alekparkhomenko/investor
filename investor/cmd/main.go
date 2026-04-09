package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/alekparkhomenko/investor/investor/internal/config"
	"github.com/alekparkhomenko/investor/investor/internal/ingestor"
	"github.com/alekparkhomenko/investor/investor/internal/model"
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

	appCtx, appCancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer appCancel()
	defer gracefulShutdown()

	closer.Configure(syscall.SIGINT, syscall.SIGTERM)
	closer.SetLogger(log)

	quotesChan := make(chan []model.Quote, 100)

	ing := ingestor.NewMOEXIngestor(cfg.App.Symbols())

	closer.AddNamed("moex-ingestor", func(ctx context.Context) error {
		ing.Stop()
		appCancel()
		return nil
	})

	go ing.Start(appCtx, cfg.App.PollInterval(), quotesChan)

	log.Info(appCtx, "started", zap.String("symbols", cfg.App.Symbols()))

	go func() {
		for quotes := range quotesChan {
			for _, q := range quotes {
				log.Info(appCtx, "quote",
					zap.String("symbol", q.Symbol),
					zap.Float64("price", q.Price),
				)
			}
		}
	}()

	<-appCtx.Done()
}

func gracefulShutdown() {
	log := logger.With(zap.String("component", "closer"))

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := closer.CloseAll(ctx); err != nil {
		log.Error(ctx, "❌ Ошибка при завершении работы", zap.Error(err))
	}
}
