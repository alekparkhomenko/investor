package portfolio

import (
	"context"
	"errors"
	"fmt"

	"github.com/alekparkhomenko/investor/investor/internal/model"
	"github.com/alekparkhomenko/investor/platform/pkg/logger"
)

// Service provides portfolio business logic.
type Service struct {
	store Store
	log   *logger.Logger
}

// NewService creates a new portfolio service.
func NewService(store Store, log *logger.Logger) *Service {
	return &Service{
		store: store,
		log:   log,
	}
}

// GetPortfolio returns a user's portfolio with all tickers and current prices.
func (s *Service) GetPortfolio(ctx context.Context, userID string) (*model.Portfolio, error) {
	p, err := s.store.GetPortfolio(ctx, userID)
	if err != nil {
		s.log.Warn(ctx, "failed to get portfolio", logger.Fields{
			"user_id": userID,
			"error":   err.Error(),
		})
		return nil, err
	}

	s.log.Info(ctx, "portfolio retrieved", logger.Fields{
		"user_id":      userID,
		"portfolio_id": p.ID,
		"ticker_count": len(p.Tickers),
	})

	return p, nil
}

// AddTickers adds symbols to a user's portfolio.
// If the portfolio does not exist, it is created automatically.
// Invalid symbols cause a ValidationError to be returned.
func (s *Service) AddTickers(ctx context.Context, userID string, symbols []string) (*model.Portfolio, error) {
	if len(symbols) == 0 {
		return nil, ErrEmptySymbols
	}

	// Validate all symbols before making any changes.
	valid, invalid, err := s.ValidateSymbols(ctx, symbols)
	if err != nil {
		return nil, fmt.Errorf("validating symbols: %w", err)
	}

	if len(invalid) > 0 {
		return nil, &ValidationError{InvalidSymbols: invalid}
	}

	// Get existing portfolio or create a new one.
	p, err := s.store.GetPortfolio(ctx, userID)
	if err != nil {
		if !errors.Is(err, ErrPortfolioNotFound) {
			return nil, fmt.Errorf("checking portfolio: %w", err)
		}

		p, err = s.store.CreatePortfolio(ctx, &model.Portfolio{UserID: userID, Name: "default"})
		if err != nil {
			return nil, fmt.Errorf("creating portfolio: %w", err)
		}

		s.log.Info(ctx, "portfolio created", logger.Fields{
			"user_id":      userID,
			"portfolio_id": p.ID,
		})
	}

	added, err := s.store.AddTickers(ctx, p.ID, valid)
	if err != nil {
		return nil, fmt.Errorf("adding tickers: %w", err)
	}

	s.log.Info(ctx, "tickers added to portfolio", logger.Fields{
		"user_id":      userID,
		"portfolio_id": p.ID,
		"added_count":  len(added),
		"symbols":      valid,
	})

	// Reload the full portfolio with tickers.
	return s.store.GetPortfolio(ctx, userID)
}

// RemoveTicker removes a ticker symbol from a user's portfolio.
func (s *Service) RemoveTicker(ctx context.Context, userID, symbol string) error {
	p, err := s.store.GetPortfolio(ctx, userID)
	if err != nil {
		return fmt.Errorf("getting portfolio: %w", err)
	}

	if err := s.store.RemoveTicker(ctx, p.ID, symbol); err != nil {
		return err
	}

	s.log.Info(ctx, "ticker removed from portfolio", logger.Fields{
		"user_id":      userID,
		"portfolio_id": p.ID,
		"symbol":       symbol,
	})

	return nil
}

// ListAvailableTickers returns all tickers available on MOEX.
func (s *Service) ListAvailableTickers(ctx context.Context) ([]model.AvailableTicker, error) {
	tickers, err := s.store.ListAvailableTickers(ctx)
	if err != nil {
		s.log.Error(ctx, "failed to list available tickers", logger.Fields{
			"error": err.Error(),
		})
		return nil, err
	}

	return tickers, nil
}

// ValidateSymbols checks which of the given symbols exist in moex_tickers.
// It returns the lists of valid and invalid symbols separately.
func (s *Service) ValidateSymbols(ctx context.Context, symbols []string) (valid []string, invalid []string, err error) {
	available, err := s.store.ListAvailableTickers(ctx)
	if err != nil {
		return nil, nil, fmt.Errorf("listing available tickers: %w", err)
	}

	availableSet := make(map[string]struct{}, len(available))
	for _, t := range available {
		availableSet[t.Symbol] = struct{}{}
	}

	for _, sym := range symbols {
		if _, ok := availableSet[sym]; ok {
			valid = append(valid, sym)
		} else {
			invalid = append(invalid, sym)
		}
	}

	return valid, invalid, nil
}
