-- +goose Up
-- Create portfolio_tickers table
CREATE TABLE IF NOT EXISTS portfolio_tickers (
    id              SERIAL PRIMARY KEY,
    portfolio_id    INTEGER NOT NULL REFERENCES portfolios(id) ON DELETE CASCADE,
    symbol          VARCHAR(20) NOT NULL,
    added_at        TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    UNIQUE(portfolio_id, symbol)
);

CREATE INDEX IF NOT EXISTS idx_portfolio_tickers_portfolio_id ON portfolio_tickers(portfolio_id);
CREATE INDEX IF NOT EXISTS idx_portfolio_tickers_symbol ON portfolio_tickers(symbol);
