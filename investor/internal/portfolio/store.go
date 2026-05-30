package portfolio

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/alekparkhomenko/investor/investor/internal/model"
)

// Store defines the portfolio storage interface.
type Store interface {
	// GetPortfolio returns a user's portfolio with all associated tickers.
	// Returns ErrPortfolioNotFound if the portfolio does not exist.
	GetPortfolio(ctx context.Context, userID string) (*model.Portfolio, error)

	// CreatePortfolio creates a new portfolio from the given model.
	CreatePortfolio(ctx context.Context, p *model.Portfolio) (*model.Portfolio, error)

	// AddTickers adds the given ticker symbols to the portfolio.
	// Duplicate tickers are silently ignored.
	AddTickers(ctx context.Context, portfolioID int, symbols []string) ([]model.Ticker, error)

	// RemoveTicker removes a single ticker from a portfolio.
	// Returns ErrTickerNotFound if the ticker is not in the portfolio.
	//
	// NOTE: This differs from Architecture v2.0 §6.1 which specified
	// RemoveTickers (plural, accepting multiple symbols). Simplified to
	// single-ticker removal per REST API design for DELETE /portfolio/{ticker}.
	RemoveTicker(ctx context.Context, portfolioID int, symbol string) error

	// GetTickers returns all tickers for a given portfolio.
	GetTickers(ctx context.Context, portfolioID int) ([]model.Ticker, error)

	// ListAvailableTickers returns all tickers available on MOEX.
	ListAvailableTickers(ctx context.Context) ([]model.AvailableTicker, error)
}

// postgresStore implements Store using a PostgreSQL connection pool.
type postgresStore struct {
	db *pgxpool.Pool
}

// NewStore creates a new Store backed by PostgreSQL.
func NewStore(db *pgxpool.Pool) Store {
	return &postgresStore{db: db}
}

// GetPortfolio returns a user's portfolio with all associated tickers.
func (s *postgresStore) GetPortfolio(ctx context.Context, userID string) (*model.Portfolio, error) {
	query := `
		SELECT id, user_id, name, created_at, updated_at
		FROM portfolios
		WHERE user_id = $1`

	var p model.Portfolio
	err := s.db.QueryRow(ctx, query, userID).Scan(
		&p.ID,
		&p.UserID,
		&p.Name,
		&p.CreatedAt,
		&p.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("%w: %s", ErrPortfolioNotFound, userID)
		}
		return nil, fmt.Errorf("querying portfolio for user %s: %w", userID, err)
	}

	tickers, err := s.GetTickers(ctx, p.ID)
	if err != nil {
		return nil, fmt.Errorf("loading tickers for portfolio %d: %w", p.ID, err)
	}
	p.Tickers = tickers

	return &p, nil
}

// CreatePortfolio creates a new portfolio from the given model.
func (s *postgresStore) CreatePortfolio(ctx context.Context, p *model.Portfolio) (*model.Portfolio, error) {
	query := `
		INSERT INTO portfolios (user_id, name)
		VALUES ($1, $2)
		RETURNING id, user_id, name, created_at, updated_at`

	var result model.Portfolio
	err := s.db.QueryRow(ctx, query, p.UserID, p.Name).Scan(
		&result.ID,
		&result.UserID,
		&result.Name,
		&result.CreatedAt,
		&result.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("creating portfolio for user %s: %w", p.UserID, err)
	}

	return &result, nil
}

// AddTickers adds the given ticker symbols to the portfolio.
// Duplicate symbols are silently ignored due to the UNIQUE constraint.
func (s *postgresStore) AddTickers(ctx context.Context, portfolioID int, symbols []string) ([]model.Ticker, error) {
	query := `
		INSERT INTO portfolio_tickers (portfolio_id, symbol)
		VALUES ($1, $2)
		ON CONFLICT (portfolio_id, symbol) DO NOTHING
		RETURNING id, portfolio_id, symbol, added_at`

	var tickers []model.Ticker

	for _, symbol := range symbols {
		var t model.Ticker
		err := s.db.QueryRow(ctx, query, portfolioID, symbol).Scan(
			&t.ID,
			&t.PortfolioID,
			&t.Symbol,
			&t.AddedAt,
		)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				// Duplicate — skip, already exists.
				continue
			}
			return nil, fmt.Errorf("adding ticker %s to portfolio %d: %w", symbol, portfolioID, err)
		}
		tickers = append(tickers, t)
	}

	return tickers, nil
}

// RemoveTicker removes a ticker from the portfolio.
func (s *postgresStore) RemoveTicker(ctx context.Context, portfolioID int, symbol string) error {
	query := `
		DELETE FROM portfolio_tickers
		WHERE portfolio_id = $1 AND symbol = $2`

	result, err := s.db.Exec(ctx, query, portfolioID, symbol)
	if err != nil {
		return fmt.Errorf("removing ticker %s from portfolio %d: %w", symbol, portfolioID, err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("%w: symbol %s in portfolio %d", ErrTickerNotFound, symbol, portfolioID)
	}

	return nil
}

// GetTickers returns all tickers for a given portfolio.
func (s *postgresStore) GetTickers(ctx context.Context, portfolioID int) ([]model.Ticker, error) {
	query := `
		SELECT id, portfolio_id, symbol, added_at
		FROM portfolio_tickers
		WHERE portfolio_id = $1
		ORDER BY added_at`

	rows, err := s.db.Query(ctx, query, portfolioID)
	if err != nil {
		return nil, fmt.Errorf("querying tickers for portfolio %d: %w", portfolioID, err)
	}
	defer rows.Close()

	var tickers []model.Ticker
	for rows.Next() {
		var t model.Ticker
		if err := rows.Scan(&t.ID, &t.PortfolioID, &t.Symbol, &t.AddedAt); err != nil {
			return nil, fmt.Errorf("scanning ticker for portfolio %d: %w", portfolioID, err)
		}
		tickers = append(tickers, t)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterating tickers for portfolio %d: %w", portfolioID, err)
	}

	return tickers, nil
}

// ListAvailableTickers returns all tickers available on MOEX.
func (s *postgresStore) ListAvailableTickers(ctx context.Context) ([]model.AvailableTicker, error) {
	query := `
		SELECT symbol, name, market, board
		FROM moex_tickers
		ORDER BY name`

	rows, err := s.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("querying available tickers: %w", err)
	}
	defer rows.Close()

	var tickers []model.AvailableTicker
	for rows.Next() {
		var t model.AvailableTicker
		if err := rows.Scan(&t.Symbol, &t.Name, &t.Market, &t.Board); err != nil {
			return nil, fmt.Errorf("scanning available ticker: %w", err)
		}
		tickers = append(tickers, t)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterating available tickers: %w", err)
	}

	return tickers, nil
}
