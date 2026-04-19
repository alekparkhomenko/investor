package app

import (
	"context"
	"os"

	"github.com/alekparkhomenko/investor/investor/internal/config"
	"github.com/alekparkhomenko/investor/investor/internal/ingestor"
	"github.com/alekparkhomenko/investor/investor/internal/metrics"
	"github.com/alekparkhomenko/investor/investor/internal/model"
	"github.com/alekparkhomenko/investor/platform/pkg/logger"
)

type App struct {
	cfg      *config.Config
	ing      ingestor.Ingestor
	quotesCh chan []model.Quote
	pidFile  string
	log      *logger.Logger
}

func NewApp(cfg *config.Config, ing ingestor.Ingestor, log *logger.Logger) *App {
	pidFile := "/tmp/investor.pid"
	if p := os.Getenv("PID_FILE"); p != "" {
		pidFile = p
	}

	return &App{
		cfg:      cfg,
		ing:      ing,
		quotesCh: make(chan []model.Quote, 100),
		pidFile:  pidFile,
		log:      log,
	}
}

func (a *App) Run(ctx context.Context) error {
	// Write PID file with optional logging
	if err := metrics.WritePID(a.pidFile, a.log); err != nil {
		a.log.Warn(ctx, "failed to write PID file", logger.Fields{
			"component": "app",
			"error":     err.Error(),
			"pid_file":  a.pidFile,
		})
	}

	a.log.Info(ctx, "starting investor", logger.Fields{
		"component": "app",
		"symbols":   a.cfg.App.Symbols(),
	})

	go func() {
		for {
			select {
			case quotes, ok := <-a.quotesCh:
				if !ok {
					a.log.Info(ctx, "quotes channel closed", logger.Fields{
						"component": "app",
					})
					return
				}
				a.log.Info(ctx, "received quotes", logger.Fields{
					"component": "app",
					"count":     len(quotes),
				})
				for _, q := range quotes {
					a.log.Debug(ctx, "quote received", logger.Fields{
						"component": "app",
						"symbol":    q.Symbol,
						"price":     q.Price,
					})
				}
			case <-ctx.Done():
				a.log.Info(ctx, "context cancelled, stopping", logger.Fields{
					"component": "app",
				})
				return
			}
		}
	}()

	go a.ing.Start(ctx, a.cfg.App.PollInterval(), a.quotesCh)

	<-ctx.Done()

	return nil
}

func (a *App) Stop() error {
	a.log.Info(context.Background(), "stopping investor", logger.Fields{
		"component": "app",
	})

	if a.ing != nil {
		a.ing.Stop()
	}

	if a.quotesCh != nil {
		close(a.quotesCh)
	}

	if a.pidFile != "" {
		if err := os.Remove(a.pidFile); err != nil {
			a.log.Warn(context.Background(), "failed to remove PID file", logger.Fields{
				"component": "app",
				"error":     err.Error(),
				"pid_file":  a.pidFile,
			})
		}
	}

	return nil
}

func (a *App) Health() bool {
	return metrics.IsProcessRunning(a.pidFile)
}
