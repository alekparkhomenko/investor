# 📝 Architecture Decision Log

## Overview

This document tracks architectural and technical decisions made during the project lifecycle.

**Format:** Architecture Decision Records (ADRs)  
**Status:** [PROPOSED | ACCEPTED | DEPRECATED | SUPERSEDED]

---

## Template

```markdown
## [ADR-XXX] Title

**Status:** [PROPOSED | ACCEPTED | DEPRECATED | SUPERSEDED]
**Date:** YYYY-MM-DD
**Deciders:** [Names]

### Context
What is the issue we're facing?

### Decision
What did we decide?

### Consequences
- Positive: [Benefits]
- Negative: [Trade-offs]
- Neutral: [Observations]

### Alternatives Considered
1. [Option A] - Why rejected
2. [Option B] - Why rejected

### References
- Links to related docs
```

---

## Decisions

### ADR-001: Use Go for Backend Implementation

**Status:** ACCEPTED  
**Date:** 2026-01-15  
**Deciders:** System Architect

**Context:**
Need to select primary backend language for MVP.

**Decision:**
Use Go (Golang) 1.25 for backend implementation.

**Consequences:**
- **Positive:** Excellent concurrency, fast compilation, strong typing
- **Positive:** Great for microservices and APIs
- **Positive:** Strong standard library
- **Negative:** Smaller ecosystem than Python/Node
- **Negative:** Steeper learning curve for some team members

**Alternatives Considered:**
1. **Python** - Great ecosystem, slower performance
2. **Node.js** - JavaScript everywhere, callback complexity
3. **Rust** - Excellent performance, steeper learning curve

---

### ADR-002: PostgreSQL as Primary Database

**Status:** ACCEPTED  
**Date:** 2026-01-15  
**Deciders:** System Architect

**Context:**
Need reliable, ACID-compliant database for user data.

**Decision:**
Use PostgreSQL 15 as primary relational database.

**Consequences:**
- **Positive:** ACID compliance, mature, excellent Go support
- **Positive:** JSON support for flexible schemas
- **Positive:** Strong consistency guarantees
- **Negative:** Horizontal scaling more complex than NoSQL
- **Neutral:** Need to design schema carefully

**Alternatives Considered:**
1. **MySQL** - Similar features, PostgreSQL has better Go ecosystem
2. **MongoDB** - Flexible schema, weaker consistency
3. **SQLite** - Great for dev, not for production scale

---

### ADR-003: REST API over gRPC (MVP)

**Status:** ACCEPTED  
**Date:** 2026-01-15  
**Deciders:** System Architect

**Context:**
Choose communication protocol for MVP.

**Decision:**
Use REST with JSON for MVP, consider gRPC for internal services in future.

**Consequences:**
- **Positive:** Universal client support
- **Positive:** Easy debugging with curl/browser
- **Positive:** Well-understood patterns
- **Positive:** No additional tooling needed
- **Negative:** Less efficient than gRPC
- **Negative:** No built-in type safety

**Alternatives Considered:**
1. **gRPC** - Better performance, requires HTTP/2
2. **GraphQL** - Flexible queries, adds complexity
3. **WebSocket** - Real-time, not needed for MVP

---

### ADR-004: Monolithic Architecture (MVP)

**Status:** ACCEPTED  
**Date:** 2026-01-15  
**Deciders:** System Architect

**Context:**
Decide between monolith and microservices for MVP.

**Decision:**
Build monolithic application for MVP, decompose later if needed.

**Consequences:**
- **Positive:** Simpler deployment and debugging
- **Positive:** Lower operational overhead
- **Positive:** Easier to refactor
- **Positive:** Faster development
- **Negative:** May need to split later
- **Neutral:** Clear module boundaries from start

**Alternatives Considered:**
1. **Microservices** - Better scalability, too complex for MVP
2. **Serverless** - Auto-scaling, cold start issues
3. **Modular Monolith** - Best of both worlds, added complexity

---

### ADR-005: JWT for Authentication

**Status:** ACCEPTED  
**Date:** 2026-01-15  
**Deciders:** System Architect

**Context:**
Select authentication mechanism.

**Decision:**
Use JWT (JSON Web Tokens) for stateless authentication.

**Consequences:**
- **Positive:** Stateless, scalable
- **Positive:** No server-side session storage needed
- **Positive:** Cross-domain compatible
- **Negative:** Token size (adds overhead)
- **Negative:** Cannot revoke immediately (need short expiry)

**Alternatives Considered:**
1. **Session Cookies** - Server-side state, harder to scale
2. **OAuth 2.0** - Complex, needed for third-party auth later
3. **API Keys** - Simpler, less secure for user auth

---

### ADR-006: Multi-Agent System Architecture

**Status:** ACCEPTED  
**Date:** 2026-01-15  
**Deciders:** Product Owner, System Architect

**Context:**
Design approach for AI-assisted software development.

**Decision:**
Use multi-agent system with clear separation:
- Product Owner (isolated)
- Planner
- Orchestrator (state management)
- Architecture
- Backend
- Reviewer

**Consequences:**
- **Positive:** Clear separation of concerns
- **Positive:** Human oversight at critical points
- **Positive:** Scalable agent specialization
- **Positive:** Reusable skills system
- **Negative:** Coordination complexity
- **Negative:** More moving parts than single agent

**Alternatives Considered:**
1. **Single Agent** - Simpler, less separation
2. **Pair Programming AI** - Collaborative but less structured
3. **Code Generation Only** - No governance or process

---

### ADR-007: STATE.json as Single Source of Truth

**Status:** ACCEPTED  
**Date:** 2026-01-15  
**Deciders:** System Architect

**Context:**
How to coordinate multiple agents and track progress.

**Decision:**
Use STATE.json file as single source of truth, with Orchestrator as sole writer.

**Consequences:**
- **Positive:** Clear coordination mechanism
- **Positive:** Human-readable state
- **Positive:** Easy to audit and debug
- **Positive:** Version controllable
- **Negative:** File I/O overhead
- **Neutral:** Need careful concurrency handling

**Alternatives Considered:**
1. **Database** - More robust, adds complexity
2. **In-Memory State** - Faster, lost on restart
3. **Event Stream** - Better for real-time, more complex

---

### ADR-005: JWT for Authentication

**Status:** ACCEPTED  
**Date:** 2026-01-15  
**Deciders:** System Architect

**Context:**
Select authentication mechanism.

**Decision:**
Use JWT (JSON Web Tokens) for stateless authentication.

**Consequences:**
- **Positive:** Stateless, scalable
- **Positive:** No server-side session storage needed
- **Positive:** Cross-domain compatible
- **Negative:** Token size (adds overhead)
- **Negative:** Cannot revoke immediately (need short expiry)

**Alternatives Considered:**
1. **Session Cookies** - Server-side state, harder to scale
2. **OAuth 2.0** - Complex, needed for third-party auth later
3. **API Keys** - Simpler, less secure for user auth

---

### ADR-006: Multi-Agent System Architecture

**Status:** ACCEPTED  
**Date:** 2026-01-15  
**Deciders:** Product Owner, System Architect

**Context:**
Design approach for AI-assisted software development.

**Decision:**
Use multi-agent system with clear separation:
- Product Owner (isolated)
- Planner
- Orchestrator (state management)
- Architecture
- Backend
- Reviewer

**Consequences:**
- **Positive:** Clear separation of concerns
- **Positive:** Human oversight at critical points
- **Positive:** Scalable agent specialization
- **Positive:** Reusable skills system
- **Negative:** Coordination complexity
- **Negative:** More moving parts than single agent

**Alternatives Considered:**
1. **Single Agent** - Simpler, less separation
2. **Pair Programming AI** - Collaborative but less structured
3. **Code Generation Only** - No governance or process

---

### ADR-007: STATE.json as Single Source of Truth

**Status:** ACCEPTED  
**Date:** 2026-01-15  
**Deciders:** System Architect

**Context:**
How to coordinate multiple agents and track progress.

**Decision:**
Use STATE.json file as single source of truth, with Orchestrator as sole writer.

**Consequences:**
- **Positive:** Clear coordination mechanism
- **Positive:** Human-readable state
- **Positive:** Easy to audit and debug
- **Positive:** Version controllable
- **Negative:** File I/O overhead
- **Neutral:** Need careful concurrency handling

**Alternatives Considered:**
1. **Database** - More robust, adds complexity
2. **In-Memory State** - Faster, lost on restart
3. **Event Stream** - Better for real-time, more complex

---

### ADR-008: Loki Stack Directory Organization

**Status:** ACCEPTED  
**Date:** 2026-04-20  
**Deciders:** Architecture Agent

**Context:**
Loki Stack YAML configuration files are scattered across the repository:
- `docker-compose.loki.yml` (root)
- `loki/loki-config.yml`
- `promtail/promtail-config.yml`
- `grafana/provisioning/datasources/loki.yml`

This makes it difficult to:
- Find all configuration files quickly
- Understand the complete Loki Stack setup
- Maintain and update configurations
- Onboard new team members

**Decision:**
Consolidate all Loki Stack YAML files into `deploy/loki/` directory with the following structure:

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

**Consequences:**
- **Positive:** All configs in single, discoverable location
- **Positive:** Self-contained deployment unit
- **Positive:** Easy to deploy: `docker-compose -f deploy/loki/docker-compose.yml up`
- **Positive:** Scalable pattern for future deployments (k8s, local, etc.)
- **Positive:** Follows golang-project-layout conventions
- **Negative:** Requires updating volume paths in docker-compose.yml
- **Negative:** One-time migration effort

**Alternatives Considered:**

1. **`loki-stack/` (root-level directory)**
   - Why rejected: Creates too many root-level directories, doesn't scale for multiple deployment targets

2. **`.infra/loki/` (hidden directory)**
   - Why rejected: Hidden directories reduce discoverability, anti-pattern for important configs

3. **Keep current structure (scattered files)**
   - Why rejected: Poor discoverability, hard to maintain, violates "single location" principle

4. **`infrastructure/loki/` (verbose naming)**
   - Why rejected: Too verbose, `deploy/` is more concise and commonly used

**Implementation Details:**

| Current Path | New Path |
|--------------|----------|
| `docker-compose.loki.yml` | `deploy/loki/docker-compose.yml` |
| `loki/loki-config.yml` | `deploy/loki/loki-config.yml` |
| `promtail/promtail-config.yml` | `deploy/loki/promtail-config.yml` |
| `grafana/provisioning/datasources/loki.yml` | `deploy/loki/grafana/provisioning/datasources/loki.yml` |

**Volume Path Changes:**
```yaml
# Before
volumes:
  - ./loki:/etc/loki
  - ./promtail:/etc/promtail
  - ./grafana/provisioning:/etc/grafana/provisioning

# After
volumes:
  - ./loki-config.yml:/etc/loki/loki-config.yml:ro
  - ./promtail-config.yml:/etc/promtail/promtail-config.yml:ro
  - ./grafana/provisioning:/etc/grafana/provisioning:ro
```

**References:**
- EXECUTION_PLAN_LOKI_STACK.md — Task ORGANIZE-001
- ARCHITECTURE.md — Section "Loki Stack Organization"
- golang-project-layout skill — Directory structure best practices

---

## Stats

**Total Decisions:** 8  
**Accepted:** 8  
**Proposed:** 0  
**Deprecated:** 0  
**Superseded:** 0

---

**Document Version:** 1.1  
**Last Updated:** 2026-04-20
