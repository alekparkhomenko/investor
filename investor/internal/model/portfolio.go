// Package model defines shared data structures for the investor application.
package model

import "time"

// Portfolio represents a user's stock portfolio.
type Portfolio struct {
	ID        int       `json:"id"`
	UserID    string    `json:"user_id"`
	Name      string    `json:"name"`
	Tickers   []Ticker  `json:"tickers,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Ticker represents a stock ticker within a user's portfolio.
type Ticker struct {
	ID           int        `json:"id"`
	PortfolioID  int        `json:"portfolio_id"`
	Symbol       string     `json:"symbol"`
	AddedAt      time.Time  `json:"added_at"`
	CurrentPrice *float64   `json:"current_price,omitempty"`
	LastUpdate   *time.Time `json:"last_update,omitempty"`
}

// AvailableTicker represents a ticker available for trading on MOEX.
type AvailableTicker struct {
	Symbol string `json:"symbol"`
	Name   string `json:"name"`
	Market string `json:"market"`
	Board  string `json:"board"`
}

