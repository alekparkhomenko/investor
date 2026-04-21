package model

import "time"

// Ticker represents an available MOEX stock ticker.
type Ticker struct {
	Symbol string `json:"symbol"`
	Name   string `json:"name"`
	Sector string `json:"sector,omitempty"`
}

// TickersResponse is the response for GET /api/v1/tickers.
type TickersResponse struct {
	Tickers []Ticker `json:"tickers"`
}

// Portfolio represents user's ticker portfolio.
type Portfolio struct {
	Tickers    []string   `json:"tickers"`
	UpdatedAt time.Time `json:"updated_at"`
}

// AddTickersRequest is the request for POST /api/v1/portfolio.
type AddTickersRequest struct {
	Tickers []string `json:"tickers" validate:"required,min=1"`
}

// AddTickersResponse is the response for POST /api/v1/portfolio.
type AddTickersResponse struct {
	Added   int      `json:"added"`
	Tickers []string `json:"tickers"`
}

// RemoveTickerResponse is the response for DELETE /api/v1/portfolio/{ticker}.
type RemoveTickerResponse struct {
	Removed bool   `json:"removed"`
	Ticker  string `json:"ticker"`
}

// ErrorResponse is the standard error response.
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
	Code    int    `json:"code"`
}