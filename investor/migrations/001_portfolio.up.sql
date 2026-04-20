-- +migrate Up
-- Create moex_tickers reference table
CREATE TABLE IF NOT EXISTS moex_tickers (
    symbol VARCHAR(10) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    sector VARCHAR(100),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Create user portfolio table
CREATE TABLE IF NOT EXISTS portfolio (
    id SERIAL PRIMARY KEY,
    ticker_symbol VARCHAR(10) NOT NULL UNIQUE,
    added_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    FOREIGN KEY (ticker_symbol) REFERENCES moex_tickers(symbol)
);

-- Index for portfolio queries
CREATE INDEX IF NOT EXISTS idx_portfolio_ticker ON portfolio(ticker_symbol);

-- Seed data: Top MOEX tickers
INSERT INTO moex_tickers (symbol, name, sector) VALUES
    ('SBER', 'ПАО Сбербанк', 'Финансы'),
    ('GAZP', 'ПАО Газпром', 'Энергетика'),
    ('LKOH', 'ПАО ЛУКОЙЛ', 'Энергетика'),
    ('NVTK', 'ПАО НОВАТЭК', 'Энергетика'),
    ('TATN', 'ПАО Татнефть', 'Энергетика'),
    ('SNGSP', 'ПАО Сургутнефтегаз', 'Энергетика'),
    ('SNGS', 'ПАО Сургутнефтегаз (прив.)', 'Энергетика'),
    ('ROSN', 'ПАО Роснефть', 'Энергетика'),
    ('MGNT', 'ПАО Магнит', 'Ритейл'),
    ('MTSS', 'ПАО МТС', 'Телеком'),
    ('ALRS', 'ПАО Алроса', 'Материалы'),
    ('POLY', 'ПАО Полиметалл', 'Материалы'),
    ('NLMK', 'ПАО НЛМК', 'Материалы'),
    ('CHMF', 'ПАО Северсталь', 'Материалы'),
    ('GMKN', 'ПАО ГМК Норникель', 'Материалы'),
    ('PLZL', 'Полиметалл', 'Материалы'),
    ('MAGN', 'ПАО ММК', 'Материалы'),
    ('RUAL', 'РУСАЛ', 'Материалы'),
    ('YDEX', 'Яндекс', 'IT'),
    ('OZON', 'Ozon Holdings', 'Ритейл'),
    ('VK', 'VK', 'IT'),
    ('MAIL', 'Мэйл.Ру', 'IT'),
    ('AFKS', 'АФК Система', 'Финансы'),
    ('VTBR', 'Банк ВТБ', 'Финансы'),
    ('SBERP', 'Сбербанк (прив.)', 'Финансы'),
    ('TATNP', 'Татнефть (прив.)', 'Энергетика'),
    ('SNGS', 'Сургутнефтегаз (прив.)', 'Энергетика'),
    ('MTSS', 'МТС (прив.)', 'Телеком'),
    ('PHOR', 'ФосАгро', 'Материалы'),
    ('HYDR', 'РусГидро', 'Коммунальные услуги')
ON CONFLICT (symbol) DO NOTHING;