package ingestor

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/alekparkhomenko/investor/investor/internal/model"
	"github.com/alekparkhomenko/investor/platform/pkg/logger"
	"go.uber.org/zap"
)

const (
	BaseURL = "https://iss.moex.com/iss/engines/stock/markets/shares"
)

type MOEXIngestor struct {
	client          *http.Client
	requiredSymbols map[string]bool
	done            chan struct{}
	mu              sync.Mutex
	stopped         bool
}

func NewMOEXIngestor(symbols string) *MOEXIngestor {
	symbolsMap := make(map[string]bool)
	for _, s := range strings.Split(symbols, ",") {
		s = strings.TrimSpace(s)
		if s != "" {
			symbolsMap[s] = true
		}
	}

	return &MOEXIngestor{
		client: &http.Client{
			Timeout: 5 * time.Second,
		},
		requiredSymbols: symbolsMap,
		done:            make(chan struct{}),
	}
}

func (m *MOEXIngestor) Start(ctx context.Context, interval time.Duration, out chan<- []model.Quote) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	log := logger.With(zap.String("component", "moex-ingestor"))

	for {
		select {
		case <-ctx.Done():
			log.Info(ctx, "stopped via context")
			return
		case <-m.done:
			log.Info(ctx, "stopped via Stop()")
			return
		case <-ticker.C:
		}

		quotes, err := m.fetchQuotes(ctx)
		if err != nil {
			log.Error(ctx, "fetch error", zap.Error(err))
			continue
		}
		if len(quotes) > 0 {
			select {
			case out <- quotes:
			case <-ctx.Done():
			case <-m.done:
				return
			}
		}
	}
}

func (m *MOEXIngestor) Stop() {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.stopped {
		return
	}
	m.stopped = true
	close(m.done)
}

func (m *MOEXIngestor) fetchQuotes(ctx context.Context) ([]model.Quote, error) {
	symbolsParam := strings.Join(func() []string {
		result := make([]string, 0, len(m.requiredSymbols))
		for s := range m.requiredSymbols {
			result = append(result, s)
		}
		return result
	}(), ",")

	url := fmt.Sprintf("%s/securities.json?secid=%s", BaseURL, symbolsParam)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, errors.Join(ErrMOEXUnavailable, err)
	}

	resp, err := m.client.Do(req)
	if err != nil {
		if ctx.Err() != nil {
			return nil, errors.Join(ErrTimeout, ctx.Err())
		}
		return nil, errors.Join(ErrMOEXUnavailable, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%w: status %d", ErrMOEXUnavailable, resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Join(ErrInvalidResponse, err)
	}

	var issResp model.ISSResponse
	if err := json.Unmarshal(body, &issResp); err != nil {
		return nil, errors.Join(ErrInvalidResponse, err)
	}

	return parseQuotes(issResp, m.requiredSymbols)
}

func parseQuotes(resp model.ISSResponse, requiredSymbols map[string]bool) ([]model.Quote, error) {
	columns := resp.MarketData.Columns
	data := resp.MarketData.Data

	secidIdx := -1
	boardIdx := -1
	lastIdx := -1
	for i, col := range columns {
		switch col {
		case "SECID":
			secidIdx = i
		case "BOARDID":
			boardIdx = i
		case "LAST":
			lastIdx = i
		}
	}

	if secidIdx == -1 || boardIdx == -1 || lastIdx == -1 {
		return nil, fmt.Errorf("%w: secid=%d, board=%d, last=%d", ErrInvalidResponse, secidIdx, boardIdx, lastIdx)
	}

	quotes := make([]model.Quote, 0, len(data))
	for _, row := range data {
		if len(row) <= secidIdx || len(row) <= boardIdx || len(row) <= lastIdx {
			continue
		}

		boardID, ok := row[boardIdx].(string)
		if !ok || boardID != "TQBR" {
			continue
		}

		symbol, ok := row[secidIdx].(string)
		if !ok {
			continue
		}

		if requiredSymbols != nil && !requiredSymbols[symbol] {
			continue
		}

		var price float64
		switch v := row[lastIdx].(type) {
		case float64:
			price = v
		case string:
			if v == "" {
				continue
			}
			p, err := strconv.ParseFloat(strings.ReplaceAll(v, ",", "."), 64)
			if err != nil {
				continue
			}
			price = p
		default:
			continue
		}

		if price <= 0 {
			continue
		}

		quotes = append(quotes, model.Quote{
			Symbol: symbol,
			Price:  price,
			Time:   time.Now(),
		})
	}

	return quotes, nil
}
