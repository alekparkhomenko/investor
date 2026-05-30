# Execution Plan — Feature #16: User Portfolio Selection

**Version:** 1.1.0  
**Status:** READY FOR APPROVAL  
**Created:** 2026-05-30 (v1.1 — updated per ADR-009: CLI → HTTP API)  
**Feature:** GitHub Issue #16  
**Architecture:** ARCHITECTURE_PORTFOLIO.md v2.0  

---

## Overview

REST HTTP API with Swagger UI for managing personal stock ticker portfolio. The application remains a single daemon with an additional HTTP server.

## Implementation Tasks

### PORT-001: Database Migrations + Seed Data
**SP:** 2 | **Priority:** P0 | **Status:** TODO

**Deliverables:**
- Migration `000001_create_portfolios` (up/down):
  - Table: `portfolios` (id, user_id, name, created_at, updated_at)
  - Index: `idx_portfolios_user_id`
  - Trigger: `update_portfolios_updated_at`
- Migration `000002_create_portfolio_tickers` (up/down):
  - Table: `portfolio_tickers` (id, portfolio_id, symbol, added_at)
  - FK: `portfolio_id → portfolios(id) ON DELETE CASCADE`
  - Unique: `(portfolio_id, symbol)`
  - Indexes: `idx_portfolio_tickers_portfolio_id`, `idx_portfolio_tickers_symbol`
- Migration `000003_seed_moex_tickers` (up/down):
  - Table: `moex_tickers` (symbol, name, market, board)
  - Seed: 20+ popular MOEX tickers (SBER, GAZP, TATN, LKOH, YNDX, etc.)
- Files: `investor/migrations/000001_*.sql`, `000002_*.sql`, `000003_*.sql` (goose format)

**Verification:**
- Migration runs without errors via goose
- Tables created: `portfolios`, `portfolio_tickers`, `moex_tickers`
- Seed data populated (query returns 20+ rows)
- Rollback works correctly

**Skills:** `golang-database`

---

### PORT-002: Data Models + Storage (PostgreSQL)
**SP:** 3 | **Priority:** P0 | **Status:** TODO | **Depends:** PORT-001

**Deliverables:**
- Go models: `investor/internal/model/portfolio.go`
  - `Portfolio`, `Ticker`, `AvailableTicker` structs with JSON tags
- DB connection: `investor/internal/db/db.go`
  - `Initialize(ctx, dsn) (*pgxpool.Pool, error)` — connection pool
  - `RunMigrations(dsn) error` — goose runner
- Portfolio Store: `investor/internal/portfolio/store.go`
  - Interface `Store` with: `GetPortfolio`, `CreatePortfolio`, `AddTickers`, `RemoveTickers`, `GetTickers`, `ListAvailableTickers`
  - Implementation `postgresStore` using pgx
- Portfolio Service: `investor/internal/portfolio/service.go`
  - `Service` struct with business logic
  - `GetPortfolioWithPrices(ctx, userID)` — portfolio + current prices (stub)
  - `ValidateSymbols(ctx, symbols)` — check against moex_tickers table
- Portfolio Errors: `investor/internal/portfolio/errors.go`
  - Sentinel errors: `ErrPortfolioNotFound`, `ErrTickerNotFound`, `ErrInvalidSymbol`

**Verification:**
- Store interface compiles
- CRUD operations work (test with pgx mock or test container)
- Service unit tests pass (mock store)
- `ValidateSymbols` returns correct valid/invalid lists

**Skills:** `golang-database`, `golang-structs-interfaces`, `golang-error-handling`

---

### PORT-003: HTTP Server + Handlers + Swagger
**SP:** 3 | **Priority:** P0 | **Status:** TODO | **Depends:** PORT-002

**Deliverables:**
- HTTP Handler: `investor/internal/http/handler.go`
  - `NewHandler(svc, log) *Handler`
  - `RegisterRoutes(mux, swagger)` — register all routes
  - Handlers: `ListTickers`, `GetPortfolio`, `AddTickers`, `RemoveTicker`
- HTTP Server: `investor/internal/http/server.go`
  - `NewServer(cfg, svc, log) *Server`
  - `Start(ctx)` / `Stop(ctx)` — lifecycle
- HTTP Middleware: `investor/internal/http/middleware.go`
  - Logging middleware (request duration, status, path)
  - CORS middleware (allow all origins for dev)
  - Panic recovery middleware
- HTTP Response: `investor/internal/http/response.go`
  - `jsonResponse(w, status, data)` helper
  - `errorResponse(w, status, code, message)` helper
- Swagger: swaggo annotations + `/swagger/index.html` endpoint

**Endpoints:**
| Method | Path | Description |
|--------|------|-------------|
| GET | `/api/v1/tickers` | List available MOEX tickers |
| GET | `/api/v1/portfolio` | Get user's portfolio |
| POST | `/api/v1/portfolio` | Add tickers to portfolio |
| DELETE | `/api/v1/portfolio/{ticker}` | Remove ticker from portfolio |

**Verification:**
- All 4 endpoints respond correctly
- JSON request/response format matches spec
- Error codes returned correctly (400, 404, 422, 500)
- Swagger UI accessible at `/swagger/index.html`
- Middleware logs requests
- Unit tests with httptest pass

**Skills:** `golang-api`, `golang-error-handling`

---

### PORT-004: Integration in main.go + Config
**SP:** 2 | **Priority:** P0 | **Status:** TODO | **Depends:** PORT-003

**Deliverables:**
- Config: `investor/internal/config/env/database.go`
  - `DSN`, `MaxOpenConns`, `MaxIdleConns`
- Config: `investor/internal/config/env/http.go`
  - `HTTP_HOST` (default: `0.0.0.0`), `HTTP_PORT` (default: `8080`)
- `investor/cmd/main.go`:
  - Initialize config → logger → DB pool → migrations
  - Wire: store → service → HTTP handler → HTTP server
  - Start HTTP server (goroutine)
  - Start MOEX poller (existing)
  - Graceful shutdown (signal.NotifyContext + closer)
- `.env.example` update with new variables

**Verification:**
- Application builds without errors
- `go vet ./...` passes
- Application starts, HTTP endpoints respond
- Graceful shutdown works (SIGINT/SIGTERM)

**Skills:** `golang-database`, `golang-design-patterns`

---

## Dependencies

```
PORT-001 (migrations) → PORT-002 (models + store) → PORT-003 (HTTP) → PORT-004 (integration)
```

| Task | Dependency | Description |
|------|-----------|-------------|
| PORT-001 | None | Create tables and seed data |
| PORT-002 | PORT-001 | Models, store, service, DB connection |
| PORT-003 | PORT-002 | HTTP handlers, server, middleware, swagger |
| PORT-004 | PORT-003 | Wire everything in main.go |

---

## Skills Required

| Task | Skills |
|------|--------|
| PORT-001 | `golang-database` |
| PORT-002 | `golang-database`, `golang-structs-interfaces`, `golang-error-handling` |
| PORT-003 | `golang-api`, `golang-error-handling` |
| PORT-004 | `golang-database`, `golang-design-patterns` |

---

## Verification Checklist

- [ ] PORT-001: Migrations execute, tables created, seed data populated
- [ ] PORT-002: Store CRUD works, service business logic works
- [ ] PORT-003: All 4 HTTP endpoints respond correctly, Swagger UI accessible
- [ ] PORT-004: Application builds, starts, and handles graceful shutdown
- [ ] Full flow: start app → `POST /api/v1/portfolio` → `GET /api/v1/portfolio` → `DELETE /api/v1/portfolio/{ticker}`

---

## Timeline Estimate

| Task | Estimate |
|------|----------|
| PORT-001 | 2 hours |
| PORT-002 | 3 hours |
| PORT-003 | 3 hours |
| PORT-004 | 2 hours |
| **Total** | **~10 hours** |

---

**Created:** Planning Agent  
**Ready for:** Human Approval