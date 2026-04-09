package ingestor

import (
	"context"
	"time"

	"github.com/alekparkhomenko/investor/investor/internal/model"
)

type Ingestor interface {
	Start(ctx context.Context, interval time.Duration, out chan<- []model.Quote)
	Stop()
}
