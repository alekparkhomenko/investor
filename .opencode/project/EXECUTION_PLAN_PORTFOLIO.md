# Execution Plan — Feature #16: User Portfolio Selection

**Version:** 1.0.0  
**Status:** DRAFT  
**Created:** 2026-04-20  
**Feature:** GitHub Issue #16  
**Architecture:** ARCHITECTURE_PORTFOLIO.md v2.0  

---

## Overview

HTTP API with Swagger UI for managing personal stock ticker portfolio.

## Implementation Tasks

### PORT-001: Database Migrations + Seed Data
**SP:** 2 | **Priority:** P0 | **Status:** TODO

**Deliverables:**
- SQL migration: `investor/migrations/001_portfolio.up.sql`
- SQL rollback: `investor/migrations/001_portfolio.down.sql`
- Seed data: 20+ MOEX tickers

**Verification:**
- Migration runs without errors
- Tables created: `moex_tickers`, `portfolio`
- Seed data populated

---

### PORT-002: Data Models + Storage (PostgreSQL)
**SP:** 3 | **Priority:** P0 | **Status:** TODO | **Depends:** PORT-001

**Deliverables:**
- Go models: `investor/internal/model/ticker.go`, `investor/internal/model/portfolio.go`
- Storage: `investor/internal/storage/portfolio.go` (PostgreSQL using pgx)

**Verification:**
- Store interface implemented
- CRUD operations work
- Unit tests pass

---

### PORT-003: HTTP Server + Handlers + Swagger
**SP:** 3 | **Priority:** P0 | **Status:** TODO | **Depends:** PORT-002

**Deliverables:**
- HTTP handler: `investor/internal/http/handler.go`
- HTTP server: `investor/internal/http/server.go`
- Swagger docs (using swaggo)

**Endpoints:**
| Method | Path | Description |
|--------|------|-------------|
| GET | `/api/v1/tickers` | List available MOEX tickers |
| GET | `/api/v1/portfolio` | Get user's portfolio |
| POST | `/api/v1/portfolio` | Add tickers to portfolio |
| DELETE | `/api/v1/portfolio/{ticker}` | Remove ticker from portfolio |

**Swagger:** `/swagger/index.html`

---

### PORT-004: Integration in main.go + Config
**SP:** 2 | **Priority:** P0 | **Status:** TODO | **Depends:** PORT-003

**Deliverables:**
- Config updates: HTTP host/port, DB connection
- Integration in `investor/cmd/main.go`
- Graceful shutdown

**Verification:**
- Application starts without errors
- HTTP endpoints respond correctly

---

## Dependencies

| Task | Dependency |
|------|-----------|
| PORT-001 | None |
| PORT-002 | PORT-001 |
| PORT-003 | PORT-002 |
| PORT-004 | PORT-003 |

---

## Skills Required

| Task | Skills |
|------|--------|
| PORT-001 | golang-database |
| PORT-002 | golang-database, golang-structs-interfaces |
| PORT-003 | golang-api, swagger |
| PORT-004 | golang-cli, golang-config |

---

## Verification Checklist

- [ ] PORT-001: Migration executes without errors
- [ ] PORT-002: Storage CRUD operations work
- [ ] PORT-003: HTTP endpoints respond correctly
- [ ] PORT-004: Application builds and runs
- [ ] Swagger UI accessible at `/swagger/index.html`

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