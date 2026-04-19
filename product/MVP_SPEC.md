# MVP Specification: Investor

## Scope Definition

### In Scope (MVP)
1. Получение котировок с MOEX ISS API
2. Graceful shutdown
3. Configuration via .env
4. Console output (fmt.Println)

### Out of Scope (Post-MVP)
- Telegram alerts
- Historical data storage
- Web dashboard
- Redis/Kafka integration

---

## User Stories

### Story 1: Basic Quote Ingestion
**As a** trader
**I want to** receive real-time quotes from MOEX
**So that** I can monitor stock prices without manual refresh

**Acceptance Criteria:**
- [x] App starts with `go run ./cmd/main.go`
- [x] Fetches SBER, GAZP, MOEX every 4 seconds
- [x] Outputs prices to console: `[APP] SBER: 324.50`
- [x] Handles network errors gracefully

**Priority:** P0
**SP:** 3

### Story 2: Configuration
**As a** trader
**I want to** configure symbols and poll interval
**So that** I can monitor different stocks

**Acceptance Criteria:**
- [ ] SYMBOLS env var (default: SBER,GAZP,MOEX)
- [ ] POLL_INTERVAL env var (default: 4s)
- [ ] .env file support via godotenv

**Priority:** P0
**SP:** 2

### Story 3: Graceful Shutdown
**As a** user
**I want to** stop app with Ctrl+C
**So that** no resources are leaked

**Acceptance Criteria:**
- [x] SIGINT/SIGTERM handled
- [x] PID file cleaned up
- [x] Channels closed

**Priority:** P0
**SP:** 1

---

## Technical Approach

### Architecture
```
┌─────────────┐     ┌─────────────┐     ┌─────────────┐
│   main.go   │────▶│     App     │────▶│   Ingestor  │
│  (entrypt)  │     │   (logic)   │     │   (MOEX)    │
└─────────────┘     └─────────────┘     └─────────────┘
                           │
                           ▼
                    ┌─────────────┐
                    │   stdout    │
                    │   (logs)    │
                    └─────────────┘
```

### Components

**ingestor/moex.go**
- `NewMOEXIngestor(symbols string)` — создаёт ingester
- `Start(ctx, interval, quotesCh)` — запускает fetch loop
- `fetchQuotes()` — делает HTTP GET к MOEX
- `parseQuotes()` — парсит JSON в Quote struct

**app/app.go**
- `App` struct — хранит config, ingestor, channel
- `Run(ctx)` — запускает ingestor и reader goroutines
- `Stop()` — graceful shutdown

**model/quote.go**
- `Quote{Symbol, Price, Time}`
- `ISSResponse` — структура для парсинга MOEX JSON

### Dependencies
- `github.com/joho/godotenv` — .env loading
- Standard library: net/http, encoding/json, time

---

## Non-Functional Requirements

### Performance
- Quote fetch: <2 seconds per request
- Memory: <50MB RSS
- CPU: <1% idle

### Reliability
- Retry on network error (1-2 attempts)
- Timeout: 5 seconds per request
- Graceful degradation if MOEX is down

### Security
- No hardcoded secrets
- Env vars for all configuration
- HTTPS only to MOEX API

---

## Open Issues

- [ ] Логирование — временно fmt.Println, нужно вернуть structured logging
- [ ] Telegram alerts — P1
- [ ] Unit tests — не написаны

---
**Version:** 0.1.0 | **Date:** 2026-04-19
