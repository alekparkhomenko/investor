-- +goose Down
DROP INDEX IF EXISTS idx_portfolio_tickers_symbol;
DROP INDEX IF EXISTS idx_portfolio_tickers_portfolio_id;
DROP TABLE IF EXISTS portfolio_tickers;
