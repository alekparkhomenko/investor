# AI Multi-Agent System — Agent Guidelines

This file provides context for AI coding agents operating in this repository using the multi-agent system architecture.

---

## 🤖 Multi-Agent System Overview

This project uses an AI multi-agent system with 6 specialized agents working together:

### Agent Roles

| Agent | Role | Responsibility |
|-------|------|----------------|
| **Product Owner** | 👤 | Defines WHAT to build and WHY. Isolated from engineering. |
| **Planner** | 🧠 | Converts backlog into detailed execution plan with technical tasks |
| **Orchestrator** | ⚙️ | Central coordinator managing STATE.json and phase transitions |
| **Architecture** | 🏗️ | Designs system structure, stack, and component interactions |
| **Backend** | 💻 | Senior engineer implementing code according to plan |
| **Reviewer** | 🔍 | Quality gatekeeper validating against plan and rules |

### System Workflow

```
┌─────────────┐     ┌─────────────┐     ┌─────────────────┐
│   Product   │────▶│   Planner   │────▶│ APPROVAL PENDING│
│   Owner     │     │             │     │   (Human Gate)  │
└─────────────┘     └─────────────┘     └─────────────────┘
                                              │
┌─────────────┐     ┌─────────────┐          ▼
│  Reviewer   │◀────│   Backend   │◀──┌─────────────┐
│   (Quality) │     │(Implementation)│  │ Architecture│
└──────┬──────┘     └─────────────┘   └─────────────┘
       │                                      ▲
       └──────────────────────────────────────┘
              (max 3 cycles, then escalate)
```

### Phase Lifecycle

```
INIT → PRODUCT_DEFINITION → PLANNING → APPROVAL_PENDING → 
ARCHITECTURE → IMPLEMENTATION → COMPLETED
```

### Critical Rules (Enforced by Reviewer)

1. **STATE.json is the Only Truth** — All state maintained in single file
2. **No Scope Expansion** — Implement ONLY what's in EXECUTION_PLAN.md
3. **Skills Must Be Reused** — No duplication, use SKILLS_INDEX.md
4. **Architecture is Immutable** — Changes require human approval
5. **Max 3 Review Loops** — Then escalate to human
6. **Product Owner Isolation** — No engineering contact
7. **Human Approval Gates** — Must stop and wait for approval

---

## 📁 System Artifacts

### Product Artifacts (in `.opencode/product/`)
- `PRODUCT_VISION.md` — Problem, solution, target audience, success metrics
- `ROADMAP.md` — Phase-based development plan
- `MVP_SPEC.md` — Detailed specs, user stories, acceptance criteria
- `BACKLOG.md` — Prioritized features (P0-P3)

### Project Artifacts (in `.opencode/project/`)
- `STATE.json` — System state (single source of truth)
- `EXECUTION_PLAN.md` — Technical tasks with deliverables
- `ARCHITECTURE.md` — System design and decisions
- `DECISIONS.md` — Architecture Decision Records
- `RULES.md` — System rules (7 strict rules)
- `SKILLS_INDEX.md` — Available skills catalog

### Agent Definitions (in `.opencode/agents/`)
- `product-owner.md` — Product Owner agent specification
- `planner.md` — Planner agent specification
- `orchestrator.md` — Orchestrator agent specification
- `architecture.md` — Architecture agent specification
- `backend.md` — Backend agent specification
- `reviewer.md` — Reviewer agent specification

---

## 📚 Legacy Reference

The original agent guidelines have been preserved in `AGENTS_LEGACY.md`.

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

## 📚 Additional Context

- See `init.md` for system architecture (Kafka, Redis, Docker Compose planned)
- See `mvp.md` for current MVP implementation details
- OpenCode agents in `.opencode/agents/` provide additional workflow guidance
- See `.opencode/SETUP.md` for multi-agent system setup instructions
