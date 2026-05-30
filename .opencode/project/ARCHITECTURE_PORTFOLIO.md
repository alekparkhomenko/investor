# Feature #11: User Portfolio Selection — Architecture Design

**Version:** 2.0.0  
**Status:** ✅ ACCEPTED  
**Created:** 2026-05-30 (v2.0 — replaced CLI with HTTP API per ADR-009)  
**Supersedes:** v1.0.0 (2026-04-20)  
**Feature:** BACKLOG.md #11  
**Skills Required:** `golang-database`, `golang-api`

---

## 1. Overview

This architecture enables users to manage their personal stock portfolio through a REST HTTP API. The application remains a daemon (existing behavior) with an additional HTTP server exposing portfolio management endpoints.

**Key Architectural Changes:**
- HTTP REST API for portfolio management (replaces CLI)
- PostgreSQL storage for user portfolios
- Swagger UI for API documentation and testing
- Integration with existing MOEX ingestor
- New HTTP handlers package

---

## 2. Component Diagram

```
┌──────────────────────────────────────────────────────────────────┐
│                      investor/cmd/main.go                       │
│  ┌────────────────────────────────────────────────────────────┐ │
│  │                    Daemon Mode (single)                    │ │
│  │  - Initialize config, logger, DB pool                     │ │
│  │  - Start MOEX poller (existing)                            │ │
│  │  - Start HTTP server (new)                                 │ │
│  │  - Graceful shutdown (both)                                │ │
│  └────────────────────────┬───────────────────────────────────┘ │
└───────────────────────────┼─────────────────────────────────────┘
                            │
                            ▼
          ┌──────────────────────────────────────────┐
          │         HTTP Server (net/http or chi)    │
          │  ┌────────────────────────────────────┐  │
          │  │        /api/v1/ Router            │  │
          │  │  ┌─────────────────────────────┐  │  │
          │  │  │  GET    /tickers            │  │  │
          │  │  │  GET    /portfolio          │  │  │
          │  │  │  POST   /portfolio          │  │  │
          │  │  │  DELETE /portfolio/{ticker}  │  │  │
          │  │  └─────────────────────────────┘  │  │
          │  └────────────────────────────────────┘  │
          │                                          │
          │  ┌────────────────────────────────────┐  │
          │  │  Swagger UI (/swagger/index.html)  │  │
          │  └────────────────────────────────────┘  │
          └──────────────────────────────────────────┘
                            │
                            ▼
          ┌──────────────────────────────────────────┐
          │        investor/internal/http/           │
          │  ┌────────────────────────────────────┐  │
          │  │  handler.go     — HTTP handlers    │  │
          │  │  server.go      — Server lifecycle │  │
          │  │  middleware.go  — Logging, CORS    │  │
          │  │  response.go   — Response helpers  │  │
          │  └────────────────────────────────────┘  │
          └──────────────────────────────────────────┘
                            │
                            ▼
          ┌──────────────────────────────────────────┐
          │        investor/internal/portfolio/      │
          │  ┌────────────────────────────────────┐  │
          │  │  store.go    — DB operations       │  │
          │  │  service.go  — Business logic      │  │
          │  │  errors.go   — Portfolio errors    │  │
          │  └────────────────────────────────────┘  │
          └──────────────────────────────────────────┘
                            │
                            ▼
          ┌──────────────────────────────────────────┐
          │         investor/internal/db/            │
          │  ┌────────────────────────────────────┐  │
          │  │  db.go  — Connection pool (pgx)    │  │
          │  └────────────────────────────────────┘  │
          └──────────────────────────────────────────┘
                            │
                            ▼
          ┌──────────────────────────────────────────┐
          │              PostgreSQL                  │
          │  (Loki Stack docker-compose)             │
          │  - portfolios                            │
          │  - portfolio_tickers                     │
          │  - moex_tickers (seed data)              │
          └──────────────────────────────────────────┘
```

---

## 3. Database Schema

### 3.1 Tables

```sql
-- portfolios: User portfolio metadata
CREATE TABLE portfolios (
    id              SERIAL PRIMARY KEY,
    user_id         VARCHAR(255) NOT NULL UNIQUE,  -- User identifier (telegram chat_id or similar)
    name            VARCHAR(255) NOT NULL DEFAULT 'default',
    created_at      TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at     TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_portfolios_user_id ON portfolios(user_id);

-- portfolio_tickers: Tickers in each portfolio
CREATE TABLE portfolio_tickers (
    id              SERIAL PRIMARY KEY,
    portfolio_id    INTEGER NOT NULL REFERENCES portfolios(id) ON DELETE CASCADE,
    symbol          VARCHAR(20) NOT NULL,
    added_at        TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    
    UNIQUE(portfolio_id, symbol)
);

CREATE INDEX idx_portfolio_tickers_portfolio_id ON portfolio_tickers(portfolio_id);
CREATE INDEX idx_portfolio_tickers_symbol ON portfolio_tickers(symbol);

-- Function to update updated_at automatically
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_portfolios_updated_at
    BEFORE UPDATE ON portfolios
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();
```

### 3.2 Migration Strategy

**Tool:** goose (github.com/pressly/goose)

**Migration Files Structure:**
```
investor/
├── migrations/
│   ├── 000001_create_portfolios.up.sql
│   ├── 000001_create_portfolios.down.sql
│   ├── 000002_create_portfolio_tickers.up.sql
│   ├── 000002_create_portfolio_tickers.down.sql
│   ├── 000003_seed_moex_tickers.up.sql
│   └── 000003_seed_moex_tickers.down.sql
```

**Note:** Schema design follows best practices from `golang-database` skill. Annotations use `-- +goose Up` / `-- +goose Down` format.

---

## 4. REST API Specification

### 4.1 Endpoints

| Method | Path | Description | Status |
|--------|------|-------------|--------|
| `GET` | `/api/v1/tickers` | List available MOEX tickers | ✅ P0 |
| `GET` | `/api/v1/portfolio` | Get user's portfolio with prices | ✅ P0 |
| `POST` | `/api/v1/portfolio` | Add tickers to portfolio | ✅ P0 |
| `DELETE` | `/api/v1/portfolio/{ticker}` | Remove ticker from portfolio | ✅ P0 |

### 4.2 Request/Response Schemas

#### `GET /api/v1/tickers`

List all available MOEX tickers from seed data.

**Response `200 OK`:**
```json
{
  "tickers": [
    {"symbol": "SBER", "name": "Сбер Банк", "market": "TQBR", "board": "EQNE"},
    {"symbol": "GAZP", "name": "Газпром", "market": "TQBR", "board": "EQNE"}
  ]
}
```

---

#### `GET /api/v1/portfolio?user_id=default`

Get user's portfolio with current market prices.

**Query parameters:**
- `user_id` (string, optional, default: `"default"`) — user identifier

**Response `200 OK`:**
```json
{
  "id": 1,
  "user_id": "default",
  "name": "default",
  "tickers": [
    {
      "symbol": "SBER",
      "current_price": 285.50,
      "last_update": "2026-05-30T10:00:00Z"
    }
  ],
  "created_at": "2026-05-30T00:00:00Z",
  "updated_at": "2026-05-30T10:00:00Z"
}
```

**Response `404 Not Found`:** Portfolio not yet created (first `POST` will create it).

---

#### `POST /api/v1/portfolio`

Add tickers to portfolio. Creates portfolio if it doesn't exist.

**Request body:**
```json
{
  "user_id": "default",
  "symbols": ["SBER", "GAZP", "TATN"]
}
```

**Validation rules:**
- `user_id` — required, non-empty string
- `symbols` — required, at least 1 symbol, max 50
- Each symbol — uppercase, 1-20 characters, letters only
- Duplicates within request — silently ignored
- Already owned symbols — silently ignored
- Invalid MOEX symbols — return 422 with list of invalid symbols

**Response `200 OK`:**
```json
{
  "portfolio": {
    "id": 1,
    "user_id": "default",
    "name": "default",
    "tickers": [
      {"symbol": "SBER", "current_price": null, "last_update": null},
      {"symbol": "GAZP", "current_price": null, "last_update": null},
      {"symbol": "TATN", "current_price": null, "last_update": null}
    ],
    "created_at": "2026-05-30T00:00:00Z",
    "updated_at": "2026-05-30T10:05:00Z"
  }
}
```

**Response `422 Unprocessable Entity`:**
```json
{
  "error": "invalid symbols",
  "invalid_symbols": ["INVALID1", "FAKE2"]
}
```

---

#### `DELETE /api/v1/portfolio/{ticker}?user_id=default`

Remove a ticker from portfolio.

**Path parameters:**
- `ticker` (string, required) — symbol to remove

**Query parameters:**
- `user_id` (string, optional, default: `"default"`)

**Response `200 OK`:**
```json
{
  "message": "ticker SBER removed from portfolio"
}
```

**Response `404 Not Found`:** Ticker not found in portfolio.

---

### 4.3 HTTP Error Response Format

All error responses follow a consistent format:

```json
{
  "error": "human-readable error message",
  "code": "ERROR_CODE"
}
```

**Standard error codes:**
| HTTP Status | Code | Description |
|-------------|------|-------------|
| 400 | `INVALID_REQUEST` | Malformed request body |
| 404 | `NOT_FOUND` | Resource not found |
| 422 | `VALIDATION_ERROR` | Invalid input data |
| 500 | `INTERNAL_ERROR` | Internal server error |

---

## 5. Data Models

### 5.1 Domain Models

```go
// investor/internal/model/portfolio.go

// Portfolio represents a user's stock portfolio
type Portfolio struct {
    ID        int       `json:"id"`
    UserID    string    `json:"user_id"`
    Name      string    `json:"name"`
    Tickers   []Ticker  `json:"tickers"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}

// Ticker represents a stock ticker
type Ticker struct {
    ID           int       `json:"id"`
    PortfolioID  int       `json:"portfolio_id"`
    Symbol      string    `json:"symbol"`
    AddedAt     time.Time `json:"added_at"`
    // Current price (populated when fetching quotes)
    CurrentPrice *float64 `json:"current_price,omitempty"`
    LastUpdate  *time.Time `json:"last_update,omitempty"`
}

// AvailableTicker represents a ticker available on MOEX
type AvailableTicker struct {
    Symbol    string `json:"symbol"`
    Name      string `json:"name"`
    Market   string `json:"market"`
    Board    string `json:"board"`
}
```

### 5.2 Database Models

```go
// investor/internal/db/model.go (for sqlx scanning)

type PortfolioRow struct {
    ID        int       `db:"id"`
    UserID    string   `db:"user_id"`
    Name     string   `db:"name"`
    CreatedAt time.Time `db:"created_at"`
    UpdatedAt time.Time `db:"updated_at"`
}

type PortfolioTickerRow struct {
    ID           int       `db:"id"`
    PortfolioID int       `db:"portfolio_id"`
    Symbol     string    `db:"symbol"`
    AddedAt    time.Time `db:"added_at"`
}
```

---

## 6. API/Interface Contracts

### 6.1 PortfolioStore Interface

```go
// investor/internal/portfolio/store.go

// Store defines the portfolio storage interface
type Store interface {
    // GetPortfolio retrieves a portfolio by user ID
    GetPortfolio(ctx context.Context, userID string) (*Portfolio, error)
    
    // CreatePortfolio creates a new portfolio
    CreatePortfolio(ctx context.Context, p *Portfolio) (*Portfolio, error)
    
    // AddTickers adds tickers to a portfolio
    AddTickers(ctx context.Context, portfolioID int, symbols []string) ([]Ticker, error)
    
    // RemoveTickers removes tickers from a portfolio
    RemoveTickers(ctx context.Context, portfolioID int, symbols []string) error
    
    // GetTickers returns all tickers in a portfolio
    GetTickers(ctx context.Context, portfolioID int) ([]Ticker, error)
}

// NewStore creates a new portfolio store
func NewStore(db *pgxpool.Pool) Store {
    return &postgresStore{db: db}
}
```

### 6.2 TickerService Interface

```go
// investor/internal/portfolio/service.go

// Service provides portfolio business logic
type Service struct {
    store Store
    log   *logger.Logger
}

func NewService(store Store, log *logger.Logger) *Service {
    return &Service{
        store: store,
        log:   log,
    }
}

// GetPortfolioWithPrices returns portfolio with current market prices
func (s *Service) GetPortfolioWithPrices(ctx context.Context, userID string) (*Portfolio, error)

// ValidateSymbols checks if symbols are valid MOEX tickers
func (s *Service) ValidateSymbols(ctx context.Context, symbols []string) (valid []string, invalid []string, err error)
```

### 6.3 HTTP Handlers Interface

```go
// investor/internal/http/handler.go

// Handler defines the HTTP handlers for portfolio API
type Handler struct {
    svc *portfolio.Service
    log *logger.Logger
}

func NewHandler(svc *portfolio.Service, log *logger.Logger) *Handler

// RegisterRoutes registers all API routes on the given mux
func (h *Handler) RegisterRoutes(mux *http.ServeMux, swagger http.Handler)
```

**Route Registration Pattern:**
```go
func (h *Handler) RegisterRoutes(mux *http.ServeMux, swagger http.Handler) {
    mux.HandleFunc("GET /api/v1/tickers", h.ListTickers)
    mux.HandleFunc("GET /api/v1/portfolio", h.GetPortfolio)
    mux.HandleFunc("POST /api/v1/portfolio", h.AddTickers)
    mux.HandleFunc("DELETE /api/v1/portfolio/{ticker}", h.RemoveTicker)
    mux.Handle("GET /swagger/", swagger)
}
```

---

## 7. Skill Mapping

| Component | Skill | Usage |
|-----------|-------|-------|
| HTTP handlers | `golang-api` | REST handlers, request parsing, response encoding |
| DB operations | `golang-database` | pgx connection, parameterized queries, error handling |
| Logging | `golang-observability` | Structured logging (existing) |
| Error handling | `golang-error-handling` | Error wrapping, sentinel errors |
| Testing | `golang-testing` | Unit tests, integration tests, HTTP tests |

### Skill-Specific Patterns

**golang-api patterns to follow:**
- Use standard `net/http` or `chi` for routing
- Structured request/response with JSON
- Consistent error response format (see 4.3)
- Middleware for logging, CORS, recovery
- Swagger/OpenAPI docs with swaggo
- HTTP handler tests with httptest

**golang-database patterns to follow:**
- Use `pgx` for PostgreSQL
- Parameterized queries ($1, $2, ...)
- Context propagation to all DB operations
- Handle sql.ErrNoRows explicitly
- Connection pool configuration
- Use migrations, not schema generation

---

## 8. Dependencies on Existing Code

### 8.1 Reused Components

| Component | Location | Reuse |
|-----------|----------|-------|
| Logger | `plantform/pkg/logger` | Full reuse |
| Config | `investor/internal/config` | Extend with DB config |
| Closer | `plantform/pkg/closer` | Full reuse |
| Model | `investor/internal/model` | Add portfolio types |
| Ingestor | `investor/internal/ingestor` | Integration for price lookups |

### 8.2 New Components

| Component | Path | Purpose |
|-----------|------|---------|
| DB Connection | `investor/internal/db/` | PostgreSQL connection pool |
| Portfolio Store | `investor/internal/portfolio/store.go` | DB operations |
| Portfolio Service | `investor/internal/portfolio/service.go` | Business logic |
| HTTP Handlers | `investor/internal/http/handler.go` | REST API handlers |
| HTTP Server | `investor/internal/http/server.go` | Server lifecycle |
| HTTP Middleware | `investor/internal/http/middleware.go` | Logging, CORS, recovery |
| HTTP Response | `investor/internal/http/response.go` | JSON response helpers |
| Migrations | `investor/migrations/` | SQL migrations |

### 8.3 Configuration Extensions

```go
// investor/internal/config/env/database.go

type databaseConfig struct {
    DSN          string `env:"DATABASE_URL,required"`
    MaxOpenConns int    `env:"DB_MAX_OPEN_CONNS" env-default:"25"`
    MaxIdleConns int    `env:"DB_MAX_IDLE_CONNS" env-default:"10"`
}

func (c *databaseConfig) ToPGXPoolConfig() pgxpool.Config {
    return pgxpool.Config{
        ConnString: c.DSN,
        MaxConns:  int32(c.MaxOpenConns),
        MinConns:  int32(c.MaxIdleConns),
    }
}
```

---

## 9. Architecture Decision Records (ADRs)

### ADR-006: Dual-Mode Application

**Status:** 🔴 SUPERSEDED (by ADR-009)  
**Date:** 2026-04-20 → 2026-05-30

**Context:**
The investor application currently runs only as a long-running daemon. Feature #11 requires portfolio management capabilities while maintaining daemon capability.

**Decision (original):**
Implement dual-mode architecture with CLI commands.

**Revised Decision:**
Single daemon mode with HTTP API (see ADR-009 in DECISIONS.md).

**Consequences:**
- ✅ SUPERSEDED by ADR-009
- ✅ HTTP API provides richer interface than CLI
- ✅ No mode complexity

---

### ADR-007: PostgreSQL with pgx

**Status:** 📋 PROPOSED  
**Date:** 2026-04-20

**Context:**
Feature requires persistent storage for user portfolios. PostgreSQL is available in Loki Stack.

**Decision:**
Use `pgx` (github.com/jackc/pgx) as the database driver.

**Consequences:**
- ✅ Native PostgreSQL support
- ✅ Connection pooling included
- ✅ Better performance than sqlx
- ⚠️ PostgreSQL-specific (not portable)

**Skills Applied:** `golang-database`

---

### ADR-008: CLI with Cobra + Viper

**Status:** 🔴 SUPERSEDED (by ADR-009)  
**Date:** 2026-04-20 → 2026-05-30

**Context:**
Need portfolio management interface.

**Decision (original):**
Use Cobra for CLI command structure.

**Revised Decision:**
Use REST HTTP API (net/http or chi router) instead of CLI.

**Consequences:**
- ✅ SUPERSEDED by ADR-009
- ✅ HTTP API is more extensible
- ✅ Swagger provides self-documenting interface

**Skills Applied:** `golang-api`

---

### ADR-010: Manual Dependency Injection

**Status:** ✅ ACCEPTED  
**Date:** 2026-05-30  
**Supersedes:** ADR-009 (v1.0 numbering)

**Context:**
Need to inject dependencies (DB pool, store, service) into HTTP handlers.

**Decision:**
Continue using manual constructor injection (existing pattern).

**Consequences:**
- ✅ Consistent with existing code
- ✅ No DI library required
- ✅ Easy to test with httptest
- ⚠️ Manual wiring in main.go

**Skills Applied:** `golang-dependency-injection`

---

## 10. Implementation Notes

### 10.1 Application Initialization Flow

```go
// investor/cmd/main.go

func main() {
    cfg, err := config.Load()
    if err != nil {
        fmt.Fprintf(os.Stderr, "failed to load config: %v\n", err)
        os.Exit(1)
    }
    
    log, err := logger.New(&cfg.Logger)
    if err != nil {
        fmt.Fprintf(os.Stderr, "failed to initialize logger: %v\n", err)
        os.Exit(1)
    }
    
    ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
    defer cancel()
    
    // Initialize DB pool
    pool, err := db.Initialize(ctx, cfg.Database.DSN)
    if err != nil {
        log.Fatal(ctx, "failed to connect to database", zap.Error(err))
    }
    defer pool.Close()
    
    // Run migrations
    if err := db.RunMigrations(cfg.Database.DSN); err != nil {
        log.Fatal(ctx, "failed to run migrations", zap.Error(err))
    }
    
    // Initialize services
    store := portfolio.NewStore(pool)
    svc := portfolio.NewService(store, log)
    
    // Start HTTP server (non-blocking)
    httpServer := http.NewServer(cfg.HTTP, svc, log)
    closer.AddNamed("http", httpServer.Stop)
    go httpServer.Start(ctx)
    
    // Start MOEX poller (existing behavior)
    app := app.New(cfg, log, pool)
    closer.AddNamed("app", app.Stop)
    app.Run(ctx)
    
    <-ctx.Done()
    log.Info(ctx, "shutting down...")
    closer.Close(ctx)
}
```

### 10.2 Database Initialization

```go
// investor/internal/db/db.go

func Initialize(ctx context.Context, dsn string) (*pgxpool.Pool, error) {
    pgxCfg, err := pgxpool.ParseConfig(dsn)
    if err != nil {
        return nil, fmt.Errorf("parsing database config: %w", err)
    }
    
    pool, err := pgxpool.NewWithConfig(ctx, pgxCfg)
    if err != nil {
        return nil, fmt.Errorf("creating connection pool: %w", err)
    }
    
    // Verify connection
    if err := pool.Ping(ctx); err != nil {
        return nil, fmt.Errorf("pinging database: %w", err)
    }
    
    return pool, nil
}
```

### 10.3 HTTP Server Configuration

```go
// investor/internal/config/env/http.go

type httpConfig struct {
    Host string `env:"HTTP_HOST" env-default:"0.0.0.0"`
    Port int    `env:"HTTP_PORT" env-default:"8080"`
}
```

**CORS:** Allow all origins in development. Production to be configured later.

**Rate Limiting:** Not implemented in MVP. Consider adding in future iteration.

---

## 11. Deployment Considerations

### 11.1 Database Setup

The PostgreSQL database is available in the Loki Stack docker-compose:

```bash
# Start Loki Stack (includes PostgreSQL)
docker-compose -f deploy/loki/docker-compose.yml up -d
```

**PostgreSQL Connection:**
- Host: `localhost` (or docker service name)
- Port: `5432`
- Database: `investor` (to be created)
- Credentials: From environment or .env file

### 11.2 Environment Variables

```bash
# Required for CLI mode
DATABASE_URL=postgres://user:password@localhost:5432/investor?sslmode=disable

# Optional (with defaults)
DB_MAX_OPEN_CONNS=25
DB_MAX_IDLE_CONNS=10

# Existing
TELEGRAM_TOKEN=your_bot_token
LOG_LEVEL=info
```

---

## 12. File Structure

```
investor/
├── cmd/
│   └── main.go              # Application entry (daemon + HTTP)
├── internal/
│   ├── db/
│   │   ├── db.go           # PostgreSQL connection pool (pgx)
│   │   └── migrate.go      # Migration runner (goose)
│   ├── http/
│   │   ├── handler.go      # REST API handlers
│   │   ├── server.go       # HTTP server lifecycle
│   │   ├── middleware.go   # Logging, CORS, recovery
│   │   └── response.go     # JSON response helpers
│   ├── portfolio/
│   │   ├── store.go        # Database operations (Store interface)
│   │   ├── service.go      # Business logic
│   │   └── errors.go       # Portfolio sentinel errors
│   └── model/
│       ├── portfolio.go    # Portfolio domain model
│       └── quote.go        # Existing (reused)
├── migrations/
│   ├── 000001_create_portfolios.up.sql
│   ├── 000001_create_portfolios.down.sql
│   ├── 000002_create_portfolio_tickers.up.sql
│   ├── 000002_create_portfolio_tickers.down.sql
│   ├── 000003_seed_moex_tickers.up.sql
│   └── 000003_seed_moex_tickers.down.sql
├── docs/
│   ├── swagger.json        # OpenAPI spec (auto-generated by swaggo)
│   └── swagger.yaml
└── go.mod
```

---

## 13. Testing Strategy

### 13.1 Unit Tests

- Portfolio store: Mock DB, test CRUD operations
- Service: Mock store, test business logic
- HTTP handlers: httptest, test request/response
- Middleware: Test logging, CORS headers

### 13.2 Integration Tests

- Database operations with test container (testcontainers-go)
- HTTP endpoint tests (full server startup)
- End-to-end: migration → seed → API call → verify DB

---

## 14. Skills Summary

| Skill | Components | Key Patterns |
|-------|-----------|-------------|
| `golang-api` | HTTP handlers, server | REST handlers, JSON, httptest |
| `golang-database` | DB connection, store | pgx, parameterized queries, transactions |
| `golang-error-handling` | All errors | Error wrapping, sentinel errors |
| `golang-testing` | Tests | Table-driven, mocks |
| `golang-observability` | Logging | Structured logging (existing) |

---

## 15. Next Steps

1. ✅ **Architecture v2.0** — Complete (you are here)
2. ⏳ **Human Approval** — Review and approve architecture v2.0
3. 🚀 **Implementation** — Backend agent implements:
   - DB connection, pool, and migration runner
   - SQL migrations (portfolios + portfolio_tickers)
   - Seed data (20+ MOEX tickers)
   - Portfolio store and service
   - HTTP handlers and server
   - Integration tests
4. 🔍 **Review** — Reviewer validates implementation

---

## 16. Document Version

**Version:** 2.0.0  
**Status:** ✅ ACCEPTED  
**Created:** 2026-05-30 (v2.0)  

**Last Updated:** 2026-05-30