// Package portfolio provides the portfolio domain: storage, business logic,
// and validation for managing user stock portfolios and MOEX tickers.
package portfolio

import "errors"

var (
	// ErrPortfolioNotFound is returned when a portfolio does not exist for a user.
	ErrPortfolioNotFound = errors.New("portfolio not found")

	// ErrTickerNotFound is returned when a ticker is not found in the portfolio.
	ErrTickerNotFound = errors.New("ticker not found")

	// ErrInvalidSymbol is returned when a ticker symbol is not found in moex_tickers.
	ErrInvalidSymbol = errors.New("invalid symbol")

	// ErrEmptySymbols is returned when an empty symbol list is provided.
	ErrEmptySymbols = errors.New("symbols list is empty")
)

// ValidationError is returned when one or more ticker symbols fail validation.
type ValidationError struct {
	InvalidSymbols []string
}

// Error implements the error interface.
func (e *ValidationError) Error() string {
	return "invalid symbols: " + func() string {
		if len(e.InvalidSymbols) == 0 {
			return ""
		}
		result := e.InvalidSymbols[0]
		for _, s := range e.InvalidSymbols[1:] {
			result += ", " + s
		}
		return result
	}()
}
