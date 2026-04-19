package app

import (
	"context"
	"fmt"
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
	if err := metrics.WritePID(a.pidFile); err != nil {
		fmt.Println("WARNING: failed to write PID file:", err)
	}

	fmt.Println("[APP] Starting investor with symbols:", a.cfg.App.Symbols())

	go func() {
		for {
			select {
			case quotes, ok := <-a.quotesCh:
				if !ok {
					fmt.Println("[APP] quotes channel closed")
					return
				}
				fmt.Printf("[APP] received %d quotes:\n", len(quotes))
				for _, q := range quotes {
					fmt.Printf("[APP]   %s: %.2f\n", q.Symbol, q.Price)
				}
			case <-ctx.Done():
				fmt.Println("[APP] context cancelled, stopping")
				return
			}
		}
	}()

	go a.ing.Start(ctx, a.cfg.App.PollInterval(), a.quotesCh)

	<-ctx.Done()

	return nil
}

func (a *App) Stop() error {
	fmt.Println("[APP] Stopping investor")

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
