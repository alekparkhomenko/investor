# AGENTS.md — Agent Guidelines for Investor Project

This file provides context for AI coding agents operating in this repository.

---

## 📁 Project Structure

This is a Go monorepo using `go.work` with two modules:

- **`investor/`** — Main application (MOEX stock quote ingestor with alert system)
- **`plantform/`** — Shared platform packages (`closer`, `logger`)

```
investor/
├── cmd/main.go              # Application entry point
├── internal/
│   ├── app/                 # Application orchestration
│   ├── config/              # Configuration loading
│   ├── ingestor/            # MOEX data ingestion (WebSocket/HTTP)
│   ├── model/               # Data models
│   └── metrics/             # Health and PID management
└── go.mod

plantform/
├── pkg/
│   ├── closer/              # Graceful shutdown handler
│   └── logger/              # Structured zap logger
└── go.mod
```

---

## 🛠 Build, Lint & Test Commands

### Running the Application

```bash
# From root (go.work context)
go run investor/cmd/main.go

# Or from module directory
cd investor && go run cmd/main.go
```

### Linting (golangci-lint)

The project uses **golangci-lint** with strict rules. Always run lint before committing:

```bash
golangci-lint run ./...
```

For faster local runs during development:

```bash
golangci-lint run --fast ./...
```

### Running Tests

```bash
# All tests across both modules
go test ./...

# Single test (run specific test function)
go test -v -run TestFunctionName ./investor/...

# With coverage
go test -cover ./...
```

### Formatting

```bash
# Format code (gofumpt)
gofmt -w investor/ plantform/

# Sort imports (gci)
gci write investor/... plantform/...
```

---

## 📏 Code Style Guidelines

### General Principles

- **Use Go 1.25** (minimum required version in `go.mod`)
- **Structured logging** with `go.uber.org/zap` — never use `fmt.Print*`
- **Graceful shutdown** via `context.Context` + signal handling
- **Defensive coding** — validate inputs, handle errors explicitly

### Import Organization (gci)

Imports must follow this order:

1. **Standard library** (`context`, `fmt`, `time`, etc.)
2. **Third-party dependencies** (`github.com/...`, `go.uber.org/...`)
3. **Local project imports** (`github.com/alekparkhomenko/investor/...`)

Example:

```go
import (
    "context"
    "errors"
    "fmt"
    "time"

    "github.com/gorilla/websocket"
    "go.uber.org/zap"

    "github.com/alekparkhomenko/investor/investor/internal/config"
    "github.com/alekparkhomenko/investor/plantform/pkg/logger"
)
```

### Naming Conventions

- **Variables**: `camelCase` (e.g., `client`, `requiredSymbols`)
- **Constants**: `CamelCase` or `UPPER_SNAKE_CASE` for grouped consts (e.g., `BaseURL`)
- **Structs/Types**: `PascalCase` (e.g., `MOEXIngestor`, `AppConfig`)
- **Interfaces**: `PascalCase` with `er` suffix (e.g., `Ingestor`, `Reader`)
- **Error variables**: Must start with `Err` (e.g., `ErrMOEXUnavailable`)

### Error Handling

- Use `errors.Join()` to combine multiple errors
- Use `errors.Is()` / `errors.As()` for error inspection — never use direct comparison
- Wrap errors with context: `fmt.Errorf("%w: ...", ErrSomething, err)`
- Never ignore errors — use `_` only when explicitly acceptable

Example:

```go
resp, err := m.client.Do(req)
if err != nil {
    if ctx.Err() != nil {
        return nil, errors.Join(ErrTimeout, ctx.Err())
    }
    return nil, errors.Join(ErrMOEXUnavailable, err)
}
```

### Context Usage

- Always pass `context.Context` as the first parameter to functions that perform I/O
- Use `contextcheck` linter — it enforces this rule
- Never store `context.Context` in structs (`containedctx` linter)

### HTTP Clients

- **Never use `http.DefaultClient`** — create custom clients with timeouts
- Always close response bodies: `defer resp.Body.Close()`

### Logging

- Use structured logger (`zap`) with component field for context
- Use appropriate log levels: `Debug` (dev), `Info` (normal), `Warn` (recoverable), `Error` (failure)

```go
log := logger.With(zap.String("component", "moex-ingestor"))
log.Info(ctx, "fetching quotes", zap.String("url", url))
```

### Cyclomatic Complexity

- Maximum complexity per function: **20** (enforced by `cyclop` linter)
- Keep functions small and focused

### Prohibited Patterns (enforced by forbidigo)

- `fmt.Print*` — use structured logger instead
- `time.Sleep` in production — use timers/context
- `http.DefaultClient` — use custom client with timeouts

---

## 🔧 Configuration

Configuration is loaded from environment variables via `.env` file:

```bash
# Required
TELEGRAM_TOKEN=your_bot_token

# Optional (with defaults)
PID_FILE=/tmp/investor.pid
LOG_LEVEL=info
LOG_JSON=false
SYMBOLS=SBER,GAZP,TATN
POLL_INTERVAL=10s
```

---

## 🔄 Graceful Shutdown Pattern

Use the `closer` package from `plantform/pkg/closer`:

```go
closer.AddNamed("app", func(ctx context.Context) error {
    return a.Stop()
})

closer.Configure(syscall.SIGINT, syscall.SIGTERM)
closer.SetLogger(log)

// Your main logic runs here

// On shutdown, defer cleaner will be called automatically
```

---

## 🤖 Agent Workflow

This project uses OpenCode agent system:

1. **coordinator** — receives tasks, delegates to planner/executor
2. **planner** — analyzes task, creates execution plan
3. **executor** — implements the plan

When making changes:
1. Run `golangci-lint run ./...` before submitting
2. Ensure tests pass (`go test ./...`)
3. Follow import ordering (gci sections)
4. Add structured logging to new components

---

## 📚 Additional Context

- See `init.md` for system architecture (Kafka, Redis, Docker Compose planned)
- See `mvp.md` for current MVP implementation details
- OpenCode agents in `.opencode/agents/` provide additional workflow guidance