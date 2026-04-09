package app

import (
	"context"
	"os"

	"github.com/alekparkhomenko/investor/investor/internal/config"
	"github.com/alekparkhomenko/investor/investor/internal/ingestor"
	"github.com/alekparkhomenko/investor/investor/internal/metrics"
	"github.com/alekparkhomenko/investor/investor/internal/model"
	"go.uber.org/zap"
)

type App struct {
	cfg      *config.Config
	log      *zap.Logger
	ing      ingestor.Ingestor
	quotesCh chan []model.Quote
	pidFile  string
}

func NewApp(cfg *config.Config, ing ingestor.Ingestor) *App {
	pidFile := "/tmp/investor.pid"
	if p := os.Getenv("PID_FILE"); p != "" {
		pidFile = p
	}

	logWithComponent := zap.L().With(zap.String("component", "app"))

	return &App{
		cfg:      cfg,
		log:      logWithComponent,
		ing:      ing,
		quotesCh: make(chan []model.Quote, 100),
		pidFile:  pidFile,
	}
}

func (a *App) Run(ctx context.Context) error {
	if err := metrics.WritePID(a.pidFile); err != nil {
		a.log.Warn("failed to write PID file", zap.Error(err))
	}

	a.log.Info("starting investor", zap.String("symbols", a.cfg.App.Symbols()))

	go func() {
		for {
			select {
			case quotes, ok := <-a.quotesCh:
				if !ok {
					return
				}
				for _, q := range quotes {
					a.log.Info("quote",
						zap.String("symbol", q.Symbol),
						zap.Float64("price", q.Price),
					)
				}
			case <-ctx.Done():
				return
			}
		}
	}()

	go a.ing.Start(ctx, a.cfg.App.PollInterval(), a.quotesCh)

	<-ctx.Done()

	return nil
}

func (a *App) Stop() error {
	a.log.Info("stopping investor")

	if a.ing != nil {
		a.ing.Stop()
	}

	if a.quotesCh != nil {
		close(a.quotesCh)
	}

	if a.pidFile != "" {
		os.Remove(a.pidFile)
	}

	return nil
}

func (a *App) Health() bool {
	return metrics.IsProcessRunning(a.pidFile)
}
