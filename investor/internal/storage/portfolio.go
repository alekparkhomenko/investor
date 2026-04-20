package storage

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/alekparkhomenko/investor/investor/internal/model"
	"github.com/jackc/pgx/v5/pgxpool"
)

// ErrTickerNotFound is returned when ticker doesn't exist in portfolio.
var ErrTickerNotFound = errors.New("ticker not found in portfolio")

// PortfolioStore handles portfolio CRUD operations.
type PortfolioStore struct {
	pool *pgxpool.Pool
}

// NewPortfolioStore creates a new PortfolioStore.
func NewPortfolioStore(pool *pgxpool.Pool) *PortfolioStore {
	return &PortfolioStore{
		pool: pool,
	}
}

// GetAllTickers returns all available MOEX tickers.
func (s *PortfolioStore) GetAllTickers(ctx context.Context) ([]model.Ticker, error) {
	query := `
		SELECT symbol, name, sector 
		FROM moex_tickers 
		ORDER BY name`

	rows, err := s.pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("querying tickers: %w", err)
	}
	defer rows.Close()

	var tickers []model.Ticker
	for rows.Next() {
		var t model.Ticker
		if err := rows.Scan(&t.Symbol, &t.Name, &t.Sector); err != nil {
			return nil, fmt.Errorf("scanning ticker: %w", err)
		}
		tickers = append(tickers, t)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterating tickers: %w", err)
	}

	return tickers, nil
}

// GetPortfolio returns user's portfolio tickers.
func (s *PortfolioStore) GetPortfolio(ctx context.Context) (*model.Portfolio, error) {
	query := `
		SELECT ticker_symbol 
		FROM portfolio 
		ORDER BY added_at`

	rows, err := s.pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("querying portfolio: %w", err)
	}
	defer rows.Close()

	var tickers []string
	for rows.Next() {
		var symbol string
		if err := rows.Scan(&symbol); err != nil {
			return nil, fmt.Errorf("scanning ticker: %w", err)
		}
		tickers = append(tickers, symbol)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterating portfolio: %w", err)
	}

	return &model.Portfolio{
		Tickers:    tickers,
		UpdatedAt: time.Now(),
	}, nil
}

// AddTickers adds tickers to user's portfolio.
func (s *PortfolioStore) AddTickers(ctx context.Context, tickers []string) ([]string, error) {
	var added []string

	for _, ticker := range tickers {
		// Validate ticker exists in moex_tickers
		var exists bool
		err := s.pool.QueryRow(ctx,
			"SELECT EXISTS(SELECT 1 FROM moex_tickers WHERE symbol = $1)",
			ticker,
		).Scan(&exists)
		if err != nil {
			return nil, fmt.Errorf("checking ticker %s: %w", ticker, err)
		}
		if !exists {
			continue // Skip invalid tickers
		}

		// Insert (ignore duplicates)
		_, err = s.pool.Exec(ctx,
			"INSERT INTO portfolio (ticker_symbol) VALUES ($1) ON CONFLICT (ticker_symbol) DO NOTHING",
			ticker,
		)
		if err != nil {
			return nil, fmt.Errorf("inserting ticker %s: %w", ticker, err)
		}

		added = append(added, ticker)
	}

	return added, nil
}

// RemoveTicker removes a ticker from user's portfolio.
func (s *PortfolioStore) RemoveTicker(ctx context.Context, ticker string) error {
	result, err := s.pool.Exec(ctx,
		"DELETE FROM portfolio WHERE ticker_symbol = $1",
		ticker,
	)
	if err != nil {
		return fmt.Errorf("removing ticker %s: %w", ticker, err)
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return ErrTickerNotFound
	}

	return nil
}