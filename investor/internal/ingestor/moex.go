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
	log             *logger.Logger
}

func NewMOEXIngestor(symbols string, log *logger.Logger) *MOEXIngestor {
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
		log:             log,
	}
}

func (m *MOEXIngestor) Start(ctx context.Context, interval time.Duration, out chan<- []model.Quote) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	m.log.Info(ctx, "starting ingestor", logger.Fields{
		"component": "moex-ingestor",
	})

	for {
		select {
		case <-ctx.Done():
			m.log.Info(ctx, "stopped via context", logger.Fields{
				"component": "moex-ingestor",
			})
			return
		case <-m.done:
			m.log.Info(ctx, "stopped via Stop()", logger.Fields{
				"component": "moex-ingestor",
			})
			return
		case <-ticker.C:
			m.log.Debug(ctx, "fetching quotes", logger.Fields{
				"component": "moex-ingestor",
			})
		}

		quotes, err := m.fetchQuotes(ctx)
		if err != nil {
			m.log.Error(ctx, "fetch error", logger.Fields{
				"component": "moex-ingestor",
				"error":     err.Error(),
			})
			continue
		}
		if len(quotes) > 0 {
			m.log.Info(ctx, "quotes fetched", logger.Fields{
				"component": "moex-ingestor",
				"count":     len(quotes),
			})
			select {
			case out <- quotes:
			case <-ctx.Done():
			case <-m.done:
				return
			}
		} else {
			m.log.Warn(ctx, "no quotes fetched", logger.Fields{
				"component": "moex-ingestor",
			})
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
	start := time.Now()

	symbolsParam := strings.Join(func() []string {
		result := make([]string, 0, len(m.requiredSymbols))
		for s := range m.requiredSymbols {
			result = append(result, s)
		}
		return result
	}(), ",")

	url := fmt.Sprintf("%s/securities.json?secid=%s", BaseURL, symbolsParam)
	m.log.Debug(ctx, "fetching quotes from MOEX", logger.Fields{
		"component": "moex-ingestor",
		"url":       url,
		"symbols":   len(m.requiredSymbols),
	})

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, errors.Join(ErrMOEXUnavailable, err)
	}

	resp, err := m.client.Do(req)
	if err != nil {
		if ctx.Err() != nil {
			m.log.Error(ctx, "request timeout", logger.Fields{
				"component": "moex-ingestor",
				"error":     err.Error(),
				"url":       url,
			})
			return nil, errors.Join(ErrTimeout, ctx.Err())
		}
		m.log.Error(ctx, "request failed", logger.Fields{
			"component": "moex-ingestor",
			"error":     err.Error(),
			"url":       url,
		})
		return nil, errors.Join(ErrMOEXUnavailable, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		m.log.Error(ctx, "unexpected status code", logger.Fields{
			"component":   "moex-ingestor",
			"url":         url,
			"status_code": resp.StatusCode,
		})
		return nil, fmt.Errorf("%w: status %d", ErrMOEXUnavailable, resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		m.log.Error(ctx, "failed to read response body", logger.Fields{
			"component": "moex-ingestor",
			"error":     err.Error(),
			"url":       url,
		})
		return nil, errors.Join(ErrInvalidResponse, err)
	}

	var issResp model.ISSResponse
	if err := json.Unmarshal(body, &issResp); err != nil {
		m.log.Error(ctx, "failed to parse response JSON", logger.Fields{
			"component": "moex-ingestor",
			"error":     err.Error(),
			"url":       url,
		})
		return nil, errors.Join(ErrInvalidResponse, err)
	}

	duration := time.Since(start).Milliseconds()
	m.log.Debug(ctx, "quotes parsed successfully", logger.Fields{
		"component":   "moex-ingestor",
		"duration_ms": duration,
		"url":         url,
	})

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
