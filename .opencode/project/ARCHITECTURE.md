# Loki Log Migration — System Architecture

**Version:** 1.1.0  
**Status:** ✅ APPROVED  
**Created:** 2026-04-19  
**Updated:** 2026-04-20  
**Immutable:** true (changes require human approval)

---

## Overview

This document describes the architecture for migrating all application logs in the `investor` project to Loki centralized logging platform. The design follows minimal-invasive principles, maintaining backward compatibility while enabling structured, searchable logging.

**Key Decisions:**
- Manual dependency injection via constructors
- Logger as concrete type (`*logger.Logger`), not interface
- Component-scoped loggers with preset fields
- Async logging via Loki library internals
- No-op behavior when Loki is disabled

---

## Architecture Diagram

```
┌─────────────────────────────────────────────────────────────────────┐
│                         investor/cmd/main.go                        │
│  ┌─────────────────────────────────────────────────────────────┐   │
│  │                    Logger Initialization                     │   │
│  │  cfg.Logger.ToPlatformLoggerConfig() → logger.New(cfg)      │   │
│  └─────────────────────────┬───────────────────────────────────┘   │
│                            │                                        │
│                            ▼                                        │
│  ┌─────────────────────────────────────────────────────────────┐   │
│  │              Dependency Injection (Manual DI)                │   │
│  │  appLogger → NewApp(cfg, ing, appLogger)                    │   │
│  │  appLogger → NewMOEXIngestor(symbols, appLogger)            │   │
│  └─────────────────────────┬───────────────────────────────────┘   │
└────────────────────────────┼───────────────────────────────────────┘
                             │
         ┌───────────────────┴───────────────────┐
         │                                       │
         ▼                                       ▼
┌───────────────────────┐            ┌───────────────────────────────┐
│   investor/internal   │            │   investor/internal           │
│        /app/          │            │      /ingestor/               │
│  ┌─────────────────┐  │            │  ┌─────────────────────────┐  │
│  │     App         │  │            │  │   MOEXIngestor          │  │
│  │─────────────────│  │            │  │─────────────────────────│  │
│  │ - log *logger   │  │            │  │ - log *logger.Logger    │  │
│  │ - cfg *Config   │  │            │  │ - client *http.Client   │  │
│  │ - ing Ingestor  │  │            │  │ - symbols map[string]   │  │
│  │ - quotesCh      │  │            │  │ - done chan struct{}    │  │
│  └─────────────────┘  │            │  └─────────────────────────┘  │
│         │             │            │           │                   │
│         ▼             │            │           ▼                   │
│  Run(ctx):            │            │  Start(ctx):                  │
│  - log.Info("start")  │            │  - log.Info("fetch")          │
│  - log.Info("stop")   │            │  - log.Error("http error")    │
│  - log.Debug("quote") │            │  - log.Warn("no quotes")      │
└───────────────────────┘            └───────────────────────────────┘
         │                                       │
         └───────────────────┬───────────────────┘
                             │
                             ▼
              ┌──────────────────────────┐
              │  platform/pkg/logger     │
              │  ┌────────────────────┐  │
              │  │   Logger struct    │  │
              │  │────────────────────│  │
              │  │ - loki *Logger     │  │
              │  │ - mu sync.Mutex    │  │
              │  │ - closed bool      │  │
              │  └────────────────────┘  │
              │         │                │
              │         ▼                │
              │  Info/Error/Warn/Debug   │
              │  (no-op if Loki disabled)│
              └──────────────────────────┘
                             │
                             ▼
              ┌──────────────────────────┐
              │   Loki Backend (opt)     │
              │   http://localhost:3100  │
              └──────────────────────────┘
```

---

## Data Flow for Logs

```
┌──────────────┐    ┌───────────────────────┐    ┌──────────────────┐
│  Component   │───▶│  logger.Logger        │───▶│  Loki (async)    │
│  (app/       │    │  - Info(ctx, msg,     │    │  - Buffered      │
│   ingestor/  │    │    fields)            │    │  - Non-blocking  │
│   metrics)   │    │  - Error(ctx, msg,    │    │  - Flush on      │
│              │    │    fields)            │    │    Close()       │
│  Fields:     │    │  - Warn/Debug/Fatal   │    │                  │
│  - component │    │                       │    │                  │
│  - duration  │    │  Internals:           │    │                  │
│  - error     │    │  - loki-logger-go     │    │                  │
│  - custom    │    │  - no-op if disabled  │    │                  │
└──────────────┘    └───────────────────────┘    └──────────────────┘
```

---

## Components

### 1. Logger Initialization (cmd/main.go)

**Responsibility:** Initialize logger and inject into components

**Structure:**
```go
func main() {
    // 1. Load config
    err := config.Load()
    
    // 2. Create logger
    loggerCfg := cfg.Logger.ToPlatformLoggerConfig()
    appLogger, err := logger.New(&loggerCfg)
    defer appLogger.Close()
    
    // 3. Inject into components
    ing := ingestor.NewMOEXIngestor(cfg.App.Symbols(), appLogger)
    a := app.NewApp(cfg, ing, appLogger)
    
    // 4. Register in closer
    closer.AddNamed("logger", func(ctx context.Context) error {
        return appLogger.Close()
    })
}
```

**Changes:** LOKI-001, LOKI-005

---

### 2. App Component (internal/app/app.go)

**Responsibility:** Application orchestration with logging

**Structure:**
```go
type App struct {
    cfg      *config.Config
    ing      ingestor.Ingestor
    quotesCh chan []model.Quote
    pidFile  string
    log      *logger.Logger  // NEW
}

func NewApp(cfg *config.Config, ing ingestor.Ingestor, log *logger.Logger) *App {
    return &App{
        cfg:      cfg,
        ing:      ing,
        quotesCh: make(chan []model.Quote, 100),
        pidFile:  "/tmp/investor.pid",
        log:      log,
    }
}

func (a *App) Run(ctx context.Context) error {
    a.log.Info(ctx, "starting investor", logger.Fields{
        "component": "app",
        "symbols":   a.cfg.App.Symbols(),
    })
    
    // ... existing logic
    
    a.log.Info(ctx, "investor stopped", logger.Fields{
        "component": "app",
    })
    
    return nil
}
```

**Changes:** LOKI-002

---

### 3. MOEXIngestor Component (internal/ingestor/moex.go)

**Responsibility:** Fetch quotes from MOEX with logging

**Structure:**
```go
type MOEXIngestor struct {
    symbols map[string]struct{}
    client  *http.Client
    done    chan struct{}
    log     *logger.Logger  // NEW
}

func NewMOEXIngestor(symbols string, log *logger.Logger) *MOEXIngestor {
    return &MOEXIngestor{
        symbols: parseSymbols(symbols),
        client:  &http.Client{Timeout: 10 * time.Second},
        done:    make(chan struct{}),
        log:     log,
    }
}

func (m *MOEXIngestor) fetchQuotes(ctx context.Context) ([]model.Quote, error) {
    start := time.Now()
    
    m.log.Info(ctx, "fetching quotes", logger.Fields{
        "component": "moex-ingestor",
        "symbols":   len(m.symbols),
    })
    
    // ... fetch logic
    
    duration := time.Since(start).Milliseconds()
    m.log.Info(ctx, "quotes fetched", logger.Fields{
        "component":  "moex-ingestor",
        "duration_ms": duration,
        "count":      len(quotes),
    })
    
    return quotes, nil
}
```

**Changes:** LOKI-003, LOKI-006, LOKI-007

---

### 4. Metrics Component (internal/metrics/health.go)

**Responsibility:** Health checks and PID management

**Structure:**
```go
func WritePID(pidFile string, log *logger.Logger) error {
    // ... existing logic
    
    if log != nil {
        log.Info(context.Background(), "PID file written", logger.Fields{
            "component": "metrics",
            "pid":       pid,
            "file":      pidFile,
        })
    }
    
    return nil
}
```

**Changes:** LOKI-004

---

## Configuration Structure

```go
// investor/internal/config/env/logger.go
type LoggerConfig struct {
    LokiEnabled bool   `env:"LOKI_ENABLED" env-default:"false"`
    LokiHost    string `env:"LOKI_HOST" env-default:"http://localhost:3100"`
    LokiEnv     string `env:"LOKI_ENV" env-default:"development"`
    AppName     string `env:"APP_NAME" env-default:"investor"`
    AppVersion  string // From build info
}

// Conversion to platform logger config
func (c *LoggerConfig) ToPlatformLoggerConfig() logger.Config {
    return logger.Config{
        LokiEnabled: c.LokiEnabled,
        LokiHost:    c.LokiHost,
        LokiEnv:     c.LokiEnv,
        AppName:     c.AppName,
        AppVersion:  c.AppVersion,
    }
}
```

---

## Testing Strategy

### Unit Tests (LOKI-008)

```go
func TestApp_Run(t *testing.T) {
    cfg := &config.Config{...}
    
    // No-op logger (Loki disabled)
    loggerCfg := logger.Config{LokiEnabled: false}
    log, _ := logger.New(&loggerCfg)
    defer log.Close()
    
    ing := ingestor.NewMOEXIngestor("SBER", log)
    app := app.NewApp(cfg, ing, log)
    
    // Test works without real Loki
}
```

### Integration Tests (LOKI-009)

```bash
#!/bin/bash
# scripts/verify-logs.sh

# 1. Start app with Loki enabled
LOKI_ENABLED=true go run investor/cmd/main.go &
APP_PID=$!

# 2. Generate test logs
sleep 5

# 3. Query Loki API
curl "http://localhost:3100/loki/api/v1/query_range?query={app=\"investor\"}"

# 4. Verify log presence
# 5. Cleanup
kill $APP_PID
```

---

## Architecture Decision Records (ADRs)

### ADR-001: Manual Dependency Injection

**Status:** ✅ Accepted  
**Date:** 2026-04-19

**Context:** Need to inject logger into multiple components

**Decision:** Use manual constructor injection (no DI library)

**Consequences:**
- ✅ Minimal code changes
- ✅ No external dependencies
- ✅ Easy to understand and test
- ⚠️ Manual wiring in main.go (acceptable for small project)

**Skills Applied:** `golang-dependency-injection`

---

### ADR-002: Logger as Concrete Type

**Status:** ✅ Accepted  
**Date:** 2026-04-19

**Context:** Whether to create Logger interface for mocking

**Decision:** Use `*logger.Logger` directly, no interface

**Consequences:**
- ✅ No premature abstraction
- ✅ Logger already no-op safe
- ✅ Simpler code
- ⚠️ Tests use real logger (no-op when disabled)

**Skills Applied:** `golang-structs-interfaces`

---

### ADR-003: Component-Scoped Loggers

**Status:** ✅ Accepted  
**Date:** 2026-04-19

**Context:** How to add component context to logs

**Decision:** Each component stores logger with preset `component` field

**Consequences:**
- ✅ Consistent log structure
- ✅ Easy filtering in Loki
- ⚠️ Slight memory overhead (negligible)

**Skills Applied:** `golang-observability`

---

### ADR-004: Async Logging via Loki Library

**Status:** ✅ Accepted  
**Date:** 2026-04-19

**Context:** How to handle log flushing without blocking

**Decision:** Rely on loki-logger-go internal async buffering

**Consequences:**
- ✅ Non-blocking logging
- ✅ Graceful shutdown via Close()
- ⚠️ Potential log loss on crash (acceptable trade-off)

**Skills Applied:** `golang-observability`

---

## Task to Component Mapping

| Task ID | Component | Changes Required |
|---------|-----------|------------------|
| **LOKI-001** | `cmd/main.go` | Logger init, closer registration |
| **LOKI-002** | `internal/app/app.go` | Add `log *logger.Logger`, update `NewApp()`, replace `fmt.Println` |
| **LOKI-003** | `internal/ingestor/moex.go` | Add `log *logger.Logger`, update `NewMOEXIngestor()`, replace `fmt.Println` |
| **LOKI-004** | `internal/metrics/health.go` | Optional logger param for `WritePID()` |
| **LOKI-005** | `cmd/main.go` | Update `gracefulShutdown()` to use logger |
| **LOKI-006** | `internal/ingestor/moex.go` | Add duration tracking in `fetchQuotes()` |
| **LOKI-007** | `internal/ingestor/moex.go` | Add sampling for "no quotes" logs |
| **LOKI-008** | Test files | Pass logger to constructors |
| **LOKI-009** | Scripts | Add Loki verification script |

---

## Implementation Principles

### 1. KISS (Keep It Simple, Stupid)
- Minimal changes to existing code
- No over-engineering
- Direct logger injection

### 2. YAGNI (You Ain't Gonna Need It)
- No Logger interface (not needed yet)
- No DI library (project is small)
- No complex log aggregation (Loki handles it)

### 3. Single Responsibility
- Each component owns its logging
- Logger package only handles logging
- Config package only handles configuration

### 4. Loose Coupling
- Logger injected, not hardcoded
- Components don't know about Loki internals
- Easy to swap logger backend if needed

### 5. Fail Fast
- Logger initialization errors logged immediately
- Application starts even if Loki unavailable (no-op mode)
- Explicit error handling throughout

---

## Deployment Architecture

```
┌──────────────────────────────────────────────────────────┐
│  Docker Container (investor)                             │
│  ┌────────────────────────────────────────────────────┐ │
│  │  Application                                       │ │
│  │  - cmd/main.go                                     │ │
│  │  - internal/app                                    │ │
│  │  - internal/ingestor                               │ │
│  │  - internal/metrics                                │ │
│  └────────────────────────────────────────────────────┘ │
│                          │                               │
│                          │ Structured logs               │
│                          ▼                               │
│  ┌────────────────────────────────────────────────────┐ │
│  │  Promtail (sidecar)                                │ │
│  │  - Reads from stdout/stderr                        │ │
│  │  - Adds labels (app, env, container_id)            │ │
│  │  - Pushes to Loki                                  │ │
│  └────────────────────────────────────────────────────┘ │
└──────────────────────────────────────────────────────────┘
                          │
                          │ HTTP/Push
                          ▼
┌──────────────────────────────────────────────────────────┐
│  Loki Service                                            │
│  - Stores logs                                           │
│  - Indexes by labels                                     │
│  - Query API (LogQL)                                     │
└──────────────────────────────────────────────────────────┘
                          │
                          │ Query
                          ▼
┌──────────────────────────────────────────────────────────┐
│  Grafana Dashboard                                       │
│  - Log exploration                                       │
│  - Alerts                                                │
│  - Metrics correlation                                   │
└──────────────────────────────────────────────────────────┘
```

---

## Skills Applied

| Skill | Usage |
|-------|-------|
| `golang-observability` | Structured logging patterns, async logging, component-scoped loggers |
| `golang-dependency-injection` | Manual constructor injection pattern |
| `golang-structs-interfaces` | Interface design decisions (when NOT to use interfaces) |

---

## Next Steps

1. ✅ **Architecture Complete** — This document is ready
2. ⏳ **Human Approval** — Review and approve architecture
3. 🚀 **Backend Implementation** — Backend agent implements tasks LOKI-001 through LOKI-009
4. 🔍 **Review** — Reviewer validates implementation against this architecture

---

## Loki Stack Organization

### Directory Structure

All Loki Stack YAML configuration files are consolidated in a single directory for improved discoverability and maintainability.

```
deploy/
  loki/
    docker-compose.yml          # Main Docker Compose file
    loki-config.yml             # Loki server configuration
    promtail-config.yml         # Promtail log collector configuration
    grafana/
      provisioning/
        datasources/
          loki.yml              # Grafana Loki datasource
        dashboards/
          investor-logs.json    # Pre-configured dashboard
          dashboard.yml         # Dashboard provisioning config
```

### File Migration Plan

| Current Path | New Path | Action |
|--------------|----------|--------|
| `docker-compose.loki.yml` (root) | `deploy/loki/docker-compose.yml` | Move + rename |
| `loki/loki-config.yml` | `deploy/loki/loki-config.yml` | Move |
| `promtail/promtail-config.yml` | `deploy/loki/promtail-config.yml` | Move |
| `grafana/provisioning/datasources/loki.yml` | `deploy/loki/grafana/provisioning/datasources/loki.yml` | Move |

### Volume Path Updates

After migration, `docker-compose.yml` volume paths must be updated:

**Before:**
```yaml
volumes:
  - ./loki:/etc/loki
  - ./promtail:/etc/promtail
  - ./grafana/provisioning:/etc/grafana/provisioning
```

**After:**
```yaml
volumes:
  - ./loki-config.yml:/etc/loki/loki-config.yml:ro
  - ./promtail-config.yml:/etc/promtail/promtail-config.yml:ro
  - ./grafana/provisioning:/etc/grafana/provisioning:ro
```

### Deployment Commands

```bash
# From project root
docker-compose -f deploy/loki/docker-compose.yml up -d

# Or from deploy/loki directory
cd deploy/loki
docker-compose up -d
```

### Design Principles

1. **Single Location** — All Loki Stack configs in one place
2. **Self-Contained** — Can be deployed independently from application
3. **Clear Structure** — Config files at root, Grafana provisioning in subdirectory
4. **Read-Only Mounts** — Config files mounted as `:ro` for security
5. **Scalable** — `deploy/` can host other deployment configurations (k8s, local, etc.)

### Architecture Decision Record

This organization follows **ADR-005: Loki Stack Organization** (see DECISIONS.md).

---

**Document Version:** 1.1.0  
**Last Updated:** 2026-04-20  
**Status:** ✅ APPROVED (Ready for Implementation)
