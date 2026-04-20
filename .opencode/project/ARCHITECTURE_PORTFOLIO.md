# Feature #11: User Portfolio Selection — Architecture Design

**Version:** 1.0.0  
**Status:** 🔄 DRAFT (Architecture Phase)  
**Created:** 2026-04-20  
**Feature:** BACKLOG.md #11  
**Skills Required:** `golang-cli`, `golang-database`

---

## 1. Overview

This architecture enables users to manage their personal stock portfolio through a CLI interface. The application transitions from a long-running daemon to a dual-mode CLI application that can run interactively or in daemon mode for background price monitoring.

**Key Architectural Changes:**
- Application modes: CLI commands vs. daemon mode
- PostgreSQL storage for user portfolios
- Integration with existing MOEX ingestor
- New command structure with Cobra

---

## 2. Component Diagram

```
┌─────────────────────────────────────────────────────────────────────────────────────┐
│                         investor/cmd/main.go                             │
│  ┌─────────────────────────────────────────────────────────────────────────┐   │
│  │                    Mode Resolution                              │   │
│  │  if args.len() == 0 → RunDaemon(ctx)                         │   │
│  │  else           → CobraCLI.Execute()                         │   │
│  └─────────────────────────┬─────────────────────────────────────┘   │
└───────────────────────────┼───────────────────────────────────────────┘
                            │
          ┌─────────────────┴─────────────────┐
          │                               │
          ▼                               ▼
┌─────────────────────┐      ┌─────────────────────────────────────────┐
│   Daemon Mode      │      │           CLI Mode                   │
│  (background)    │      │  ┌─────────────────────────────┐   │
│                  │      │  │   Root Command             │   │
│  - app.Run(ctx)  │      │  │   - PersistentFlags       │   │
│  - MOEX polling │      │  │   - PreRunE            │   │
│  - PID file     │      │  └──────────┬──────────────┘   │
│                  │      │             │                    │
│  ┌────────────┐ │      │  ┌─────────┴──────┐            │
│  │App        │ │      │  │               │             │
│  │(existing)│ │      │  ▼               ▼             │
│  └──────────┘ │      │  ┌────────┐ ┌──────────┐    │
│               │      │  │ticker  │ │portfolio │    │
│  ┌────────────┐│      │  │ list   │ │ add/rm   │    │
│  │MOEXIngstor│ ���      │  └────────┘ └──────────┘    │
│  │(existing)│ │      │                     │        │
│  └──────────┘ │      │  ┌──────────────┬──────┐   │
└────────────────┘      │ portfolio   │show  │    │
                         │ show        │      │    │
                         └─────────────┴──────┘    │
                             │                    │
                             └────────┬───────────┘
                                    │
                                    ▼
                         ┌───────────────────────┐
                         │  investor/internal    │
                         │        /portfolio/    │
                         │  ┌───────────────┐   │
                         │  │ Repository  │   │
                         │  │ (DB operations)│   │
                         │  └───────────────┘   │
                         │  ┌───────────────┐   │
                         │  │  Service   │   │
                         │  │ (business logic)│   │
                         │  └───────────────┘   │
                         │  ┌───────────────┐   │
                         │  │    Store   │   │
                         │  │ (interface)  │   │
                         │  └───────────────┘   │
                         │          │          │
                         └──────────┼──────────┘
                                    │
                                    ▼
                         ┌───────────────────────┐
                         │ investor/internal     │
                         │        /db/          │
                         │  ┌───────────────┐   │
                         │  │  Connection │   │
                         │  │  Pool (pgx) │   │
                         │  └───────────────┘   │
                         │  ┌───────────────┐   │
                         │  │ Migration    │   │
                         │  │ (SQL files)  │   │
                         │  └───────────────┘   │
                         └───────────────────────┘
                                    │
                                    ▼
                         ┌───────────────────────┐
                         │  PostgreSQL         │
                         │  (Loki Stack)      │
                         │  - portfolios     │
                         │  - portfolio_tickers│
                         └───────────────────────┘
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

**Tool:** golang-migrate (CLI + library)

**Migration Files Structure:**
```
investor/
├── migrations/
│   ├── 000001_create_portfolios.down.sql
│   ├── 000001_create_portfolios.up.sql
│   ├── 000002_create_portfolio_tickers.down.sql
│   └── 000002_create_portfolio_tickers.up.sql
```

**Note:** Schema design follows best practices from `golang-database` skill. Migration SQL should be reviewed by humans before application.

---

## 4. CLI Command Structure

### 4.1 Command Tree

```
investor
├── [root]           # Default: daemon mode
├── ticker
│   └── list         # Show available MOEX tickers
├── portfolio
│   ├── show        # Show user's portfolio
│   ├── add          # Add tickers to portfolio
│   └── remove      # Remove tickers from portfolio
└── serve           # Explicit daemon mode
```

### 4.2 Command Details

| Command | Description | Arguments |
|---------|-------------|-----------|
| `investor` | Run as daemon (background mode) | None |
| `investor ticker list` | List available MOEX tickers | Optional: `--market TQBR` |
| `investor portfolio show` | Show user's portfolio with prices | None |
| `investor portfolio add <symbols...>` | Add tickers to portfolio | At least 1 symbol |
| `investor portfolio remove <symbols...>` | Remove tickers from portfolio | At least 1 symbol |
| `investor serve` | Explicit daemon mode | None |

### 4.3 Flag Structure

**Global Flags (Persistent):**
- `--config` - Config file path (default: `.investor.yaml`)
- `--log-level` - Log level (debug, info, warn, error)
- `--db-dsn` - PostgreSQL connection string

**Per-Command Flags:**
- `ticker list`: `--market` (TQBR, TQOB, etc.)
- `portfolio add`: `--dry-run` (validate without saving)

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

### 6.3 CLI Commands Interface

```go
// Commands are defined in investor/cmd/investor/

// Each command file provides:
// - Use: string (command name)
// - Short: string (description)
// - RunE: func(cmd *cobra.Command, args []string) error

// All commands follow cobra conventions from golang-cli skill
```

---

## 7. Skill Mapping

| Component | Skill | Usage |
|-----------|-------|-------|
| CLI commands | `golang-cli` | Command structure, Cobra + Viper, flag handling |
| DB operations | `golang-database` | pgx connection, parameterized queries, error handling |
| Logging | `golang-observability` | Structured logging (existing) |
| Error handling | `golang-error-handling` | Error wrapping, sentinel errors |
| Testing | `golang-testing` | Unit tests, integration tests |

### Skill-Specific Patterns

**golang-cli patterns to follow:**
- Root command with PersistentPreRunE for config init
- Subcommands in separate files
- Bind all flags to Viper
- stdout for data, stderr for logs
- Exit codes on errors

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
| CLI Commands | `investor/cmd/investor/` | Cobra commands |
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

**Status:** 📋 PROPOSED  
**Date:** 2026-04-20

**Context:**
The investor application currently runs only as a long-running daemon. Feature #11 requires CLI commands for portfolio management while maintaining daemon capability.

**Decision:**
Implement dual-mode architecture:
- No arguments → daemon mode (existing behavior)
- With arguments → CLI mode (cobra commands)

**Consequences:**
- ✅ Backward compatible with existing deployments
- ✅ Single binary serves both use cases
- ⚠️ Requires careful signal handling between modes

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

**Status:** 📋 PROPOSED  
**Date:** 2026-04-20

**Context:**
Need CLI interface for portfolio management commands.

**Decision:**
Use Cobra for command structure and Viper for configuration layering.

**Consequences:**
- ✅ Industry standard (kubectl, docker, etc.)
- ✅ Automatic completions
- ✅ Config file + env var + flag layering
- ⚠️ Additional dependency

**Skills Applied:** `golang-cli`

---

### ADR-009: Manual Dependency Injection

**Status:** 📋 PROPOSED  
**Date:** 2026-04-20

**Context:**
Need to inject dependencies into CLI commands and services.

**Decision:**
Continue using manual constructor injection (existing pattern from ADR-001).

**Consequences:**
- ✅ Consistent with existing code
- ✅ No DI library required
- ✅ Easy to test
- ⚠️ Manual wiring in main.go

**Skills Applied:** `golang-dependency-injection`

---

## 10. Implementation Notes

### 10.1 Mode Resolution Flow

```go
// investor/cmd/main.go

func main() {
    // Load configuration
    if err := config.Load(); err != nil {
        fmt.Fprintf(os.Stderr, "failed to load config: %v\n", err)
        os.Exit(1)
    }
    
    cfg := config.AppConfig()
    
    // Initialize logger (existing)
    appLogger, err := logger.New(&cfg.Logger)
    if err != nil {
        fmt.Fprintf(os.Stderr, "failed to initialize logger: %v\n", err)
        os.Exit(1)
    }
    
    // Mode resolution
    if len(os.Args) == 1 || os.Args[1] == "serve" {
        // Daemon mode (existing behavior)
        runDaemon(cfg, appLogger)
    } else {
        // CLI mode
        if err := runCLI(cfg, appLogger); err != nil {
            // Exit code handled by cobra
            os.Exit(1)
        }
    }
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

### 10.3 Error Codes for CLI

Following Unix conventions from `golang-cli` skill:

| Exit Code | Meaning | Usage |
|----------|---------|-------|
| 0 | Success | Normal completion |
| 1 | General error | Runtime failures |
| 2 | Usage error | Invalid flags/arguments |

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
│   ├── main.go              # Mode resolution
│   └── investor/
│       ├── root.go          # Root command
│       ├── ticker/
│       │   └── list.go     # ticker list command
│       └── portfolio/
│           ├── show.go     # portfolio show command
│           ├── add.go      # portfolio add command
│           └── remove.go   # portfolio remove command
├── internal/
│   ├── db/
│   │   └── db.go           # PostgreSQL connection
│   ├── portfolio/
│   │   ├── store.go        # Database operations
│   │   ├── service.go     # Business logic
│   │   └── errors.go      # Portfolio errors
│   └── model/
│       ├── portfolio.go   # Portfolio model (new)
│       └── quote.go       # Existing (reused)
├── migrations/
│   ├── 000001_create_portfolios.up.sql
│   ├── 000001_create_portfolios.down.sql
│   ├── 000002_create_portfolio_tickers.up.sql
│   └── 000002_create_portfolio_tickers.down.sql
└── go.mod
```

---

## 13. Testing Strategy

### 13.1 Unit Tests

- Portfolio store: Mock DB, test CRUD operations
- Service: Mock store, test business logic
- Commands: Execute and capture output

### 13.2 Integration Tests

- Database operations with test container
- CLI command execution
- End-to-end portfolio workflow

---

## 14. Skills Summary

| Skill | Components | Key Patterns |
|-------|-----------|-------------|
| `golang-cli` | CLI commands | Cobra + Viper, flag binding, exit codes |
| `golang-database` | DB connection, store | pgx, parameterized queries, transactions |
| `golang-error-handling` | All errors | Error wrapping, sentinel errors |
| `golang-testing` | Tests | Table-driven, mocks |
| `golang-observability` | Logging | Structured logging (existing) |

---

## 15. Next Steps

1. 🔄 **This Architecture** — Architecture Agent completes design
2. ⏳ **Human Approval** — Review and approve architecture
3. 🚀 **Implementation** — Backend agent implements:
   - DB connection and migrations
   - Portfolio store and service
   - CLI commands
   - Integration tests
4. 🔍 **Review** — Reviewer validates implementation

---

## 16. Document Version

**Version:** 1.0.0  
**Status:** 🔄 DRAFT  
**Created:** 2026-04-20  

**Last Updated:** 2026-04-20