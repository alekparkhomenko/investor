-- +goose Up
-- Create moex_tickers reference table
CREATE TABLE IF NOT EXISTS moex_tickers (
    symbol   VARCHAR(20) PRIMARY KEY,
    name     VARCHAR(255) NOT NULL,
    market   VARCHAR(20) NOT NULL DEFAULT 'TQBR',
    board    VARCHAR(20) NOT NULL DEFAULT 'EQNE'
);

-- Seed data: popular MOEX tickers
INSERT INTO moex_tickers (symbol, name, market, board) VALUES
    ('SBER',  'Сбер Банк',                    'TQBR', 'EQNE'),
    ('SBERP', 'Сбер Банк (прив.)',            'TQBR', 'EQNE'),
    ('GAZP',  'Газпром',                      'TQBR', 'EQNE'),
    ('TATN',  'Татнефть',                     'TQBR', 'EQNE'),
    ('LKOH',  'Лукойл',                        'TQBR', 'EQNE'),
    ('YNDX',  'Яндекс',                        'TQBR', 'EQNE'),
    ('ROSN',  'Роснефть',                      'TQBR', 'EQNE'),
    ('MOEX',  'Московская Биржа',              'TQBR', 'EQNE'),
    ('VTBR',  'ВТБ',                           'TQBR', 'EQNE'),
    ('NVTK',  'Новатэк',                       'TQBR', 'EQNE'),
    ('SNGS',  'Сургутнефтегаз',                'TQBR', 'EQNE'),
    ('SNGSP', 'Сургутнефтегаз (прив.)',        'TQBR', 'EQNE'),
    ('GMKN',  'Норникель',                     'TQBR', 'EQNE'),
    ('PLZL',  'Полюс',                         'TQBR', 'EQNE'),
    ('NLMK',  'НЛМК',                          'TQBR', 'EQNE'),
    ('MAGN',  'ММК',                           'TQBR', 'EQNE'),
    ('AFLT',  'Аэрофлот',                      'TQBR', 'EQNE'),
    ('RTKM',  'Ростелеком',                    'TQBR', 'EQNE'),
    ('PHOR',  'Фосагро',                       'TQBR', 'EQNE'),
    ('MGNT',  'Магнит',                        'TQBR', 'EQNE'),
    ('ALRS',  'АЛРОСА',                        'TQBR', 'EQNE'),
    ('IRAO',  'Интер РАО',                     'TQBR', 'EQNE'),
    ('RUAL',  'РУСАЛ',                         'TQBR', 'EQNE'),
    ('FEES',  'ФСК Россети',                   'TQBR', 'EQNE')
ON CONFLICT (symbol) DO NOTHING;
