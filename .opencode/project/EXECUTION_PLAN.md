# Execution Plan: Loki Logging Migration

## Metadata

| Field | Value |
|-------|-------|
| **Version** | 1.0 |
| **Status** | DRAFT |
| **Approval Status** | PENDING |
| **Target** | Migrate all application logs to Loki |
| **Created** | 2026-04-19 |
| **Last Updated** | 2026-04-19 |

---

## Executive Summary

This plan describes the migration of the `investor` application from `fmt.Println` and `log.Printf` to the structured Loki logger available in `plantform/pkg/logger`. The migration will provide centralized log aggregation, better observability, and production-ready logging capabilities.

### Scope
- **In Scope**: Replace all `fmt.Println`, `fmt.Printf`, `log.Println`, `log.Printf` calls with structured Loki logging
- **Out of Scope**: Adding new features, changing business logic, modifying Loki logger implementation

### Success Criteria
1. All `fmt.Println` and `log.Println` calls replaced with structured logger
2. All log messages include appropriate context fields
3. Log levels (Info, Warn, Error, Debug) used correctly
4. No functionality regression
5. Application builds and runs successfully with Loki enabled/disabled

---

## Phase 1: Foundation & Configuration (4 hours)

### Task 1.1: Update Logger Configuration
**ID:** LOKI-001  
**Priority:** P0  
**Estimate:** 1 hour  
**Dependencies:** None

**Description:**
Update the logger configuration to properly initialize Loki logger in `main.go` with correct app version and environment settings.

**Deliverables:**
- [ ] `investor/cmd/main.go` updated to properly initialize logger
- [ ] Logger config includes app version from build info
- [ ] Error handling for logger initialization follows single-handling rule

**Acceptance Criteria:**
- [ ] Logger initialized before any logging occurs
- [ ] Logger properly closed on shutdown via `closer`
- [ ] No `log.Printf` calls in main.go for logger errors
- [ ] Uses `golang-error-handling` skill patterns

**Required Skills:**
- `golang-error-handling` - error wrapping and single-handling rule
- `golang-observability` - structured logging patterns

---

### Task 1.2: Inject Logger into App Component
**ID:** LOKI-002  
**Priority:** P0  
**Estimate:** 1.5 hours  
**Dependencies:** LOKI-001

**Description:**
Modify the `app.App` struct to accept and use the logger instance instead of `fmt.Println`.

**Deliverables:**
- [ ] `app.App` struct includes `log *logger.Logger` field
- [ ] `NewApp()` constructor accepts logger parameter
- [ ] All `fmt.Println` calls in `app.Run()` replaced with `log.Info()`
- [ ] All `fmt.Printf` calls replaced with structured logging

**Acceptance Criteria:**
- [ ] App compiles without errors
- [ ] All print statements replaced (verified via grep)
- [ ] Log messages include component field: `component=app`
- [ ] Context passed to all logger calls

**Required Skills:**
- `golang-observability` - structured logging with context
- `golang-project-layout` - dependency injection patterns

---

### Task 1.3: Inject Logger into Ingestor Component
**ID:** LOKI-003  
**Priority:** P0  
**Estimate:** 2 hours  
**Dependencies:** LOKI-001

**Description:**
Modify the `ingestor.MOEXIngestor` struct to accept and use the logger instance instead of `fmt.Println`.

**Deliverables:**
- [ ] `MOEXIngestor` struct includes `log *logger.Logger` field
- [ ] `NewMOEXIngestor()` accepts logger parameter
- [ ] All `fmt.Println` and `fmt.Printf` calls replaced with structured logging
- [ ] HTTP request/response logging includes relevant fields

**Acceptance Criteria:**
- [ ] Ingestor compiles without errors
- [ ] All print statements replaced (verified via grep)
- [ ] Log messages include component field: `component=moex-ingestor`
- [ ] Fetch operations logged with URL, status code, duration
- [ ] Errors logged with `log.Error()` and proper error fields

**Required Skills:**
- `golang-observability` - HTTP request logging patterns
- `golang-error-handling` - error logging (single-handling rule)
- `golang-structs-interfaces` - struct modification patterns

---

## Phase 2: Metrics & Health Logging (2 hours)

### Task 2.1: Add Logging to Metrics Package
**ID:** LOKI-004  
**Priority:** P1  
**Estimate:** 1 hour  
**Dependencies:** LOKI-002

**Description:**
Update the `metrics` package to use structured logging for PID file operations and health checks.

**Deliverables:**
- [ ] `metrics.WritePID()` accepts optional logger parameter
- [ ] PID write failures logged with `log.Warn()`
- [ ] Health check failures logged appropriately

**Acceptance Criteria:**
- [ ] No `fmt.Println` in metrics package
- [ ] Warning-level logs for non-critical failures
- [ ] Log messages include component field: `component=metrics`

**Required Skills:**
- `golang-observability` - log level selection
- `golang-error-handling` - error handling patterns

---

### Task 2.2: Add Graceful Shutdown Logging
**ID:** LOKI-005  
**Priority:** P1  
**Estimate:** 1 hour  
**Dependencies:** LOKI-001, LOKI-002

**Description:**
Replace `fmt.Println` calls in `gracefulShutdown()` with structured logging and add detailed shutdown sequence logging.

**Deliverables:**
- [ ] `gracefulShutdown()` uses logger from main
- [ ] Shutdown sequence logged step-by-step
- [ ] Errors during shutdown logged with `log.Error()`
- [ ] Shutdown completion logged with `log.Info()`

**Acceptance Criteria:**
- [ ] No `fmt.Println` in shutdown function
- [ ] Shutdown logs include timing information
- [ ] Errors properly wrapped and logged

**Required Skills:**
- `golang-observability` - structured logging
- `golang-error-handling` - error wrapping

---

## Phase 3: Enhanced Observability (3 hours)

### Task 3.1: Add Request Duration Tracking
**ID:** LOKI-006  
**Priority:** P1  
**Estimate:** 1.5 hours  
**Dependencies:** LOKI-003

**Description:**
Add duration tracking to MOEX HTTP requests and log fetch duration with each request.

**Deliverables:**
- [ ] `fetchQuotes()` tracks request duration
- [ ] Duration logged in milliseconds
- [ ] Slow requests (>1s) logged as warnings
- [ ] Duration included in log fields

**Acceptance Criteria:**
- [ ] All fetch operations include `duration_ms` field
- [ ] Slow requests trigger warning-level logs
- [ ] No performance regression introduced

**Required Skills:**
- `golang-observability` - metrics and logging correlation
- `golang-context` - context propagation for tracing

---

### Task 3.2: Add Log Sampling for High-Frequency Logs
**ID:** LOKI-007  
**Priority:** P2  
**Estimate:** 1.5 hours  
**Dependencies:** LOKI-003

**Description:**
Implement log sampling for high-frequency "no quotes fetched" messages to prevent log spam.

**Deliverables:**
- [ ] Sampling logic for repetitive log messages
- [ ] Configurable sample rate (e.g., log 1 in 10)
- [ ] Sampled logs include indicator field

**Acceptance Criteria:**
- [ ] High-frequency logs reduced by sampling
- [ ] Sample rate configurable via environment
- [ ] Sampled logs identifiable in Loki

**Required Skills:**
- `golang-observability` - log sampling patterns
- `golang-concurrency` - thread-safe sampling counters

---

## Phase 4: Testing & Validation (3 hours)

### Task 4.1: Update Unit Tests
**ID:** LOKI-008  
**Priority:** P0  
**Estimate:** 1.5 hours  
**Dependencies:** LOKI-002, LOKI-003

**Description:**
Update existing unit tests to work with the new logger injection pattern.

**Deliverables:**
- [ ] Tests updated to pass logger instances
- [ ] Mock logger for unit tests (no-op or test logger)
- [ ] All tests pass with new logger pattern

**Acceptance Criteria:**
- [ ] `go test ./...` passes
- [ ] Test coverage maintained or improved
- [ ] No test uses real Loki logger

**Required Skills:**
- `golang-testing` - test patterns with dependencies
- `golang-structs-interfaces` - interface-based testing

---

### Task 4.2: Integration Testing with Loki
**ID:** LOKI-009  
**Priority:** P1  
**Estimate:** 1.5 hours  
**Dependencies:** LOKI-001, LOKI-006

**Description:**
Verify logs are correctly sent to Loki when enabled and application works with Loki disabled.

**Deliverables:**
- [ ] Test script to verify Loki log ingestion
- [ ] Verification that app works with `LOKI_ENABLED=false`
- [ ] Log correlation verified (trace IDs if applicable)

**Acceptance Criteria:**
- [ ] Logs visible in Loki when enabled
- [ ] App runs without errors when Loki disabled
- [ ] Log format matches expected schema

**Required Skills:**
- `golang-observability` - Loki verification

---

## Dependencies Graph

```
LOKI-001 (Logger Config)
    ├── LOKI-002 (App Logger)
    │       ├── LOKI-004 (Metrics Logging)
    │       └── LOKI-005 (Shutdown Logging)
    │               └── LOKI-008 (Unit Tests)
    │
    ├── LOKI-003 (Ingestor Logger)
    │       ├── LOKI-006 (Duration Tracking)
    │       │       └── LOKI-009 (Integration Test)
    │       └── LOKI-007 (Log Sampling)
    │
    └── LOKI-008 (Unit Tests)
```

---

## Risk Assessment

| Risk | Probability | Impact | Mitigation |
|------|-------------|--------|------------|
| Logger initialization fails silently | Medium | High | Add explicit error logging, fail-fast on critical errors |
| Performance degradation from logging | Low | Medium | Use async logging, implement sampling (LOKI-007) |
| Log volume too high for Loki | Medium | Medium | Implement sampling, configure retention policies |
| Breaking changes to existing tests | High | Low | Update tests systematically, use mock logger |
| Context propagation missed | Medium | Low | Code review checklist, linting rules |

---

## Required Skills Summary

| Skill | Tasks Using It |
|-------|----------------|
| `golang-observability` | LOKI-001, LOKI-002, LOKI-003, LOKI-004, LOKI-005, LOKI-006, LOKI-007, LOKI-009 |
| `golang-error-handling` | LOKI-001, LOKI-003, LOKI-004, LOKI-005 |
| `golang-project-layout` | LOKI-002 |
| `golang-structs-interfaces` | LOKI-003, LOKI-008 |
| `golang-context` | LOKI-006 |
| `golang-concurrency` | LOKI-007 |
| `golang-testing` | LOKI-008 |

---

## Effort Estimation

### By Phase

| Phase | Tasks | Estimated Hours | Story Points |
|-------|-------|-----------------|--------------|
| Phase 1: Foundation | LOKI-001, LOKI-002, LOKI-003 | 4.5 | 8 |
| Phase 2: Metrics & Health | LOKI-004, LOKI-005 | 2 | 3 |
| Phase 3: Enhanced Observability | LOKI-006, LOKI-007 | 3 | 5 |
| Phase 4: Testing & Validation | LOKI-008, LOKI-009 | 3 | 5 |
| **Total** | **9 tasks** | **12.5 hours** | **21 SP** |

### By Priority

| Priority | Tasks | Estimated Hours |
|----------|-------|-----------------|
| P0 | LOKI-001, LOKI-002, LOKI-003, LOKI-008 | 6 |
| P1 | LOKI-004, LOKI-005, LOKI-006, LOKI-009 | 5.5 |
| P2 | LOKI-007 | 1.5 |

---

## Implementation Notes

### Code Changes Required

**Files to modify:**
1. `investor/cmd/main.go` - Logger initialization and injection
2. `investor/internal/app/app.go` - Replace fmt.Println with logger
3. `investor/internal/ingestor/moex.go` - Replace fmt.Println with logger
4. `investor/internal/metrics/health.go` - Add optional logging
5. Test files - Update to pass logger instances

### Files NOT to modify:
- `plantform/pkg/logger/logger.go` - Logger implementation is complete
- `investor/internal/config/env/logger.go` - Config is complete

### Logging Standards

**Log levels:**
- `Debug` - Development details (URL, raw data)
- `Info` - Normal operations (started, stopped, fetched N quotes)
- `Warn` - Recoverable issues (slow response, PID write failure)
- `Error` - Failures (HTTP errors, context cancellation)

**Required fields:**
- `component` - Source component (app, moex-ingestor, metrics)
- `duration_ms` - For timed operations
- `error` - For error logs

**Context propagation:**
- Always pass `context.Context` to logger methods
- Use `logger.InfoContext()` when available

---

## Approval Checklist

Before implementation can begin, verify:

- [ ] All P0 tasks clearly defined with deliverables
- [ ] Dependencies mapped correctly
- [ ] Required skills identified
- [ ] Effort estimates reasonable
- [ ] Risks identified with mitigations
- [ ] No scope creep beyond logging migration
- [ ] Architecture unchanged (logger injection only)

---

## Next Steps

**After human approval:**
1. Orchestrator updates STATE.json to `ARCHITECTURE` phase
2. Architecture Agent reviews plan and creates ARCHITECTURE.md if needed
3. Backend Agent implements tasks phase by phase
4. Reviewer validates each phase against acceptance criteria
5. Maximum 3 review cycles before escalation

---

**Status:** APPROVAL_PENDING  
**Human Review Required:** YES  
**Cannot proceed without approval**
