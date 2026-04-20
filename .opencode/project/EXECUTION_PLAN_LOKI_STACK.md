# Execution Plan: Loki Stack Deployment

## Metadata

| Field | Value |
|-------|-------|
| **Version** | 1.1 |
| **Status** | DRAFT |
| **Approval Status** | PENDING |
| **Target** | Deploy and validate Loki Stack (Loki, Promtail, Grafana) for centralized log collection |
| **Created** | 2026-04-20 |
| **Last Updated** | 2026-04-20 |

---

## Executive Summary

This plan describes the deployment of the Loki Stack infrastructure for centralized log aggregation. The application has already been migrated to use structured logging (LOKI-001–LOKI-009 completed). This plan focuses on **infrastructure deployment, configuration validation, and integration testing**.

### Scope
- **In Scope**: 
  - Consolidate YAML configuration files in single directory
  - Deploy Loki, Promtail, Grafana via Docker Compose
  - Validate configurations
  - Verify log ingestion from investor application
  - Create Grafana dashboards
  - Document operations procedures
- **Out of Scope**: 
  - Changing application code (already migrated)
  - Modifying Loki logger implementation
  - Adding new features to application

### Success Criteria
1. All YAML configs organized in single directory with Taskfile commands
2. Loki Stack runs without errors via `docker-compose -f docker-compose.loki.yml up`
3. Promtail successfully ships logs to Loki
4. Grafana can query logs from Loki
5. Application logs visible in Grafana Explore
6. Documentation complete for operations

---

## Phase 0: Project Organization (1.5 hours)

### Task 0.1: Consolidate YAML Configuration Files
**ID:** ORGANIZE-001  
**Priority:** P0  
**Estimate:** 0.5 hours  
**Dependencies:** None

**Description:**
Consolidate all scattered YAML configuration files into a single dedicated directory for improved discoverability and maintainability.

**Current State:**
- `docker-compose.loki.yml` (root)
- `loki/loki-config.yml`
- `promtail/promtail-config.yml`
- `grafana/provisioning/datasources/loki.yml`

**Deliverables:**
- [ ] New directory structure created for YAML configs
- [ ] All YAML files moved to consolidated location
- [ ] Volume paths in docker-compose updated
- [ ] No broken references after reorganization

**Acceptance Criteria:**
- [ ] All YAML configs located in single directory (e.g., `deploy/loki/` or `infrastructure/loki/`)
- [ ] Directory structure documented in README
- [ ] Docker Compose volume paths updated correctly
- [ ] `docker-compose -f <file> config` validates successfully after move
- [ ] Uses `golang-project-layout` skill for directory structure best practices

**Required Skills:**
- `golang-project-layout` — Directory structure and organization

---

### Task 0.2: Update Taskfile Commands
**ID:** ORGANIZE-002  
**Priority:** P0  
**Estimate:** 0.5 hours  
**Dependencies:** ORGANIZE-001

**Description:**
Add deployment and management commands to Taskfile.yaml for easy access to Loki Stack operations.

**Deliverables:**
- [ ] Taskfile.yaml updated with Loki Stack commands
- [ ] Commands self-documenting with descriptions
- [ ] All common operations covered

**Acceptance Criteria:**
- [ ] Command `task loki:up` — starts Loki Stack
- [ ] Command `task loki:down` — stops Loki Stack
- [ ] Command `task loki:logs` — shows logs from all services
- [ ] Command `task loki:validate` — validates configurations
- [ ] Command `task loki:status` — shows service health
- [ ] All commands work with new directory structure
- [ ] Uses `golang-cli` skill for command design patterns

**Required Skills:**
- `golang-cli` — Command structure and user experience

---

### Task 0.3: Update Documentation
**ID:** ORGANIZE-003  
**Priority:** P0  
**Estimate:** 0.5 hours  
**Dependencies:** ORGANIZE-001, ORGANIZE-002

**Description:**
Document the new directory structure and Taskfile commands for developer onboarding.

**Deliverables:**
- [ ] README section added for Loki Stack organization
- [ ] Directory structure diagram
- [ ] Quick start guide with Taskfile commands
- [ ] Migration notes (if applicable)

**Acceptance Criteria:**
- [ ] New team member can find configs in <30 seconds
- [ ] Commands documented with examples
- [ ] Directory structure explained
- [ ] Uses `golang-documentation` skill for technical writing

**Required Skills:**
- `golang-documentation` — Developer documentation

---

## Phase 1: Infrastructure Setup (4 hours)

### Task 1.1: Validate Docker Compose Configuration
**ID:** LOKI-DEPLOY-001  
**Priority:** P0  
**Estimate:** 1 hour  
**Dependencies:** None

**Description:**
Review and validate `docker-compose.loki.yml` for correctness, including service definitions, volumes, networks, and health checks.

**Deliverables:**
- [ ] `docker-compose.loki.yml` validated for syntax and structure
- [ ] Health checks added for all services (loki, promtail, grafana)
- [ ] Resource limits configured (memory, CPU)
- [ ] Restart policies verified
- [ ] Network configuration documented

**Acceptance Criteria:**
- [ ] `docker-compose -f docker-compose.loki.yml config` passes validation
- [ ] Health check endpoints defined for all services
- [ ] Memory limits set (Loki: 512MB, Promtail: 128MB, Grafana: 256MB)
- [ ] All volumes properly mapped
- [ ] Uses `golang-project-layout` skill for Docker best practices

**Required Skills:**
- `golang-project-layout` — Docker Compose best practices
- `golang-observability` — Health check patterns

---

### Task 1.2: Create Loki Configuration File
**ID:** LOKI-DEPLOY-002  
**Priority:** P0  
**Estimate:** 1.5 hours  
**Dependencies:** LOKI-DEPLOY-001

**Description:**
Create the missing `loki/loki-config.yml` configuration file with appropriate settings for development/production use.

**Deliverables:**
- [ ] `loki/loki-config.yml` created with complete configuration
- [ ] Schema configured for single-binary mode
- [ ] Retention policies defined (7 days for dev, 30 days for prod)
- [ ] Compactor settings configured
- [ ] Limits configured (ingestion rate, query timeout)

**Acceptance Criteria:**
- [ ] Loki starts successfully with config file
- [ ] No configuration errors in logs
- [ ] Retention policies active
- [ ] Ingestion limits prevent resource exhaustion
- [ ] Uses `golang-observability` skill for Loki configuration patterns

**Required Skills:**
- `golang-observability` — Loki configuration and tuning

---

### Task 1.3: Validate Promtail Configuration
**ID:** LOKI-DEPLOY-003  
**Priority:** P0  
**Estimate:** 1 hour  
**Dependencies:** LOKI-DEPLOY-001

**Description:**
Review and enhance `promtail/promtail-config.yml` to ensure proper log collection from investor application and Docker containers.

**Deliverables:**
- [ ] `promtail-config.yml` validated for syntax
- [ ] Docker container log scraping configured
- [ ] Investor application log path configured
- [ ] Log labels properly defined (app, env, container_name)
- [ ] Pipeline stages validated for JSON parsing

**Acceptance Criteria:**
- [ ] `promtail -config.file=promtail/promtail-config.yml -check-config` passes
- [ ] Docker container logs scraped correctly
- [ ] Labels include: `job`, `app`, `container_name`, `env`
- [ ] JSON pipeline stages extract structured fields
- [ ] Uses `golang-observability` skill for Promtail patterns

**Required Skills:**
- `golang-observability` — Promtail configuration and log scraping

---

### Task 1.4: Validate Grafana Provisioning
**ID:** LOKI-DEPLOY-004  
**Priority:** P0  
**Estimate:** 0.5 hours  
**Dependencies:** LOKI-DEPLOY-001

**Description:**
Verify Grafana datasource provisioning and create initial dashboard for log exploration.

**Deliverables:**
- [ ] `grafana/provisioning/datasources/loki.yml` validated
- [ ] Loki datasource auto-configured on startup
- [ ] Initial dashboard provisioned for log exploration
- [ ] Derived fields configured for trace/metric correlation

**Acceptance Criteria:**
- [ ] Grafana starts with Loki datasource pre-configured
- [ ] Dashboard accessible at `http://localhost:3000`
- [ ] Default credentials work (admin/admin)
- [ ] Derived fields extract `duration_ms` from logs
- [ ] Uses `golang-observability` skill for Grafana provisioning

**Required Skills:**
- `golang-observability` — Grafana provisioning and dashboards

---

## Phase 2: Integration & Testing (3 hours)

### Task 2.1: Deploy Loki Stack
**ID:** LOKI-DEPLOY-005  
**Priority:** P0  
**Estimate:** 1 hour  
**Dependencies:** LOKI-DEPLOY-001, LOKI-DEPLOY-002, LOKI-DEPLOY-003, LOKI-DEPLOY-004

**Description:**
Deploy the complete Loki Stack using Docker Compose and verify all services start correctly.

**Deliverables:**
- [ ] Deployment script created (`scripts/deploy-loki.sh`)
- [ ] All services start successfully
- [ ] Health checks pass for all services
- [ ] Services accessible on expected ports (Loki: 3100, Grafana: 3000, Promtail: 9080)

**Acceptance Criteria:**
- [ ] `docker-compose -f docker-compose.loki.yml up -d` completes without errors
- [ ] `docker-compose ps` shows all services healthy
- [ ] Loki API responds: `curl http://localhost:3100/ready`
- [ ] Grafana UI accessible: `http://localhost:3000`
- [ ] Promtail connected to Loki (verified via logs)

**Required Skills:**
- `golang-observability` — Deployment and verification

---

### Task 2.2: Verify Log Ingestion
**ID:** LOKI-DEPLOY-006  
**Priority:** P0  
**Estimate:** 1 hour  
**Dependencies:** LOKI-DEPLOY-005

**Description:**
Start the investor application and verify logs are successfully ingested into Loki.

**Deliverables:**
- [ ] Investor application started with `LOKI_ENABLED=true`
- [ ] Logs visible in Loki via API query
- [ ] Logs visible in Grafana Explore
- [ ] Structured fields properly indexed

**Acceptance Criteria:**
- [ ] Application logs appear in Loki within 5 seconds
- [ ] Query `{app="investor"}` returns logs
- [ ] Fields `component`, `duration_ms`, `level` searchable
- [ ] No errors in Promtail logs
- [ ] Uses `golang-observability` skill for log verification

**Required Skills:**
- `golang-observability` — LogQL queries and verification

---

### Task 2.3: Create Grafana Dashboard
**ID:** LOKI-DEPLOY-007  
**Priority:** P1  
**Estimate:** 1 hour  
**Dependencies:** LOKI-DEPLOY-006

**Description:**
Create a Grafana dashboard for monitoring application logs with panels for log volume, error rates, and log exploration.

**Deliverables:**
- [ ] Dashboard JSON exported (`grafana/dashboards/investor-logs.json`)
- [ ] Panel: Log volume over time
- [ ] Panel: Error rate (logs with level=error)
- [ ] Panel: Log exploration table
- [ ] Panel: Request duration histogram (from `duration_ms` field)
- [ ] Dashboard provisioned via `grafana/provisioning/dashboards/`

**Acceptance Criteria:**
- [ ] Dashboard auto-loaded on Grafana startup
- [ ] All panels display data correctly
- [ ] Time range selector works
- [ ] Log filtering by level works
- [ ] Uses `golang-observability` skill for dashboard design

**Required Skills:**
- `golang-observability` — Grafana dashboard creation

---

## Phase 3: Operations & Documentation (2 hours)

### Task 3.1: Create Operations Runbook
**ID:** LOKI-DEPLOY-008  
**Priority:** P1  
**Estimate:** 1 hour  
**Dependencies:** LOKI-DEPLOY-005

**Description:**
Document operational procedures for Loki Stack including startup, shutdown, troubleshooting, and maintenance.

**Deliverables:**
- [ ] `docs/loki-operations.md` created
- [ ] Startup procedures documented
- [ ] Shutdown procedures documented
- [ ] Troubleshooting guide (common issues)
- [ ] Log retention and cleanup procedures
- [ ] Backup/restore procedures for Grafana dashboards

**Acceptance Criteria:**
- [ ] Runbook covers all operational scenarios
- [ ] Troubleshooting section includes: Loki unavailable, Promtail disconnected, Grafana can't query
- [ ] Commands copy-paste ready
- [ ] Uses `golang-documentation` skill for technical writing

**Required Skills:**
- `golang-documentation` — Operations documentation

---

### Task 3.2: Add Monitoring & Alerts
**ID:** LOKI-DEPLOY-009  
**Priority:** P2  
**Estimate:** 1 hour  
**Dependencies:** LOKI-DEPLOY-006

**Description:**
Configure basic alerts for Loki Stack health and log anomalies.

**Deliverables:**
- [ ] Alert rules defined in `grafana/provisioning/alerting/`
- [ ] Alert: Loki unavailable
- [ ] Alert: Promtail disconnected
- [ ] Alert: High error rate in application logs
- [ ] Alert: Log ingestion stopped

**Acceptance Criteria:**
- [ ] Alerts provisioned on Grafana startup
- [ ] Alert notifications configured (can use console for MVP)
- [ ] Alert rules use correct LogQL syntax
- [ ] Uses `golang-observability` skill for alerting patterns

**Required Skills:**
- `golang-observability` — Alerting and monitoring

---

## Dependencies Graph

```
ORGANIZE-001 (Consolidate YAML Files)
    ├── ORGANIZE-002 (Update Taskfile)
    │       │
    │       └── ORGANIZE-003 (Update Documentation)
    │
    └── LOKI-DEPLOY-001 (Validate Docker Compose)
            ├── LOKI-DEPLOY-002 (Loki Config)
            ├── LOKI-DEPLOY-003 (Promtail Config)
            └── LOKI-DEPLOY-004 (Grafana Provisioning)
                    │
                    └── LOKI-DEPLOY-005 (Deploy Stack)
                            │
                            ├── LOKI-DEPLOY-006 (Verify Logs)
                            │       │
                            │       ├── LOKI-DEPLOY-007 (Grafana Dashboard)
                            │       └── LOKI-DEPLOY-009 (Alerts)
                            │
                            └── LOKI-DEPLOY-008 (Operations Runbook)
```

---

## Risk Assessment

| Risk | Probability | Impact | Mitigation |
|------|-------------|--------|------------|
| Loki configuration errors | Medium | High | Validate config with `loki -verify-config`, start with minimal config |
| Promtail can't read Docker logs | High | Medium | Ensure proper volume mounts, check permissions, verify Docker socket access |
| Grafana can't connect to Loki | Medium | High | Verify network configuration, use service names in URLs |
| High memory usage | Low | Medium | Configure resource limits, tune retention policies |
| Log volume too high | Medium | Low | Implement log sampling in Promtail, configure retention |
| Missing `loki/loki-config.yml` | High | High | Create config file as part of LOKI-DEPLOY-002 |

---

## Required Skills Summary

| Skill | Tasks Using It |
|-------|----------------|
| `golang-observability` | LOKI-DEPLOY-001, LOKI-DEPLOY-002, LOKI-DEPLOY-003, LOKI-DEPLOY-004, LOKI-DEPLOY-005, LOKI-DEPLOY-006, LOKI-DEPLOY-007, LOKI-DEPLOY-009 |
| `golang-project-layout` | ORGANIZE-001, LOKI-DEPLOY-001 |
| `golang-documentation` | ORGANIZE-003, LOKI-DEPLOY-008 |
| `golang-cli` | ORGANIZE-002 |

---

## Effort Estimation

### By Phase

| Phase | Tasks | Estimated Hours | Story Points |
|-------|-------|-----------------|--------------|
| Phase 0: Project Organization | ORGANIZE-001, ORGANIZE-002, ORGANIZE-003 | 1.5 | 2 |
| Phase 1: Infrastructure Setup | LOKI-DEPLOY-001, LOKI-DEPLOY-002, LOKI-DEPLOY-003, LOKI-DEPLOY-004 | 4 | 8 |
| Phase 2: Integration & Testing | LOKI-DEPLOY-005, LOKI-DEPLOY-006, LOKI-DEPLOY-007 | 3 | 5 |
| Phase 3: Operations & Documentation | LOKI-DEPLOY-008, LOKI-DEPLOY-009 | 2 | 3 |
| **Total** | **12 tasks** | **10.5 hours** | **18 SP** |

### By Priority

| Priority | Tasks | Estimated Hours |
|----------|-------|-----------------|
| P0 | ORGANIZE-001, ORGANIZE-002, ORGANIZE-003, LOKI-DEPLOY-001, LOKI-DEPLOY-002, LOKI-DEPLOY-003, LOKI-DEPLOY-004, LOKI-DEPLOY-005, LOKI-DEPLOY-006 | 7.5 |
| P1 | LOKI-DEPLOY-007, LOKI-DEPLOY-008 | 2 |
| P2 | LOKI-DEPLOY-009 | 1 |

---

## Implementation Notes

### Files to Create/Modify

**Files to create:**
1. New directory structure for YAML configs (ORGANIZE-001) — exact structure determined by Architecture Agent
2. `grafana/provisioning/dashboards/investor-logs.json` — Dashboard (LOKI-DEPLOY-007)
3. `grafana/provisioning/dashboards/loki-dashboard.yml` — Dashboard provisioning (LOKI-DEPLOY-007)
4. `grafana/provisioning/alerting/alerts.yml` — Alert rules (LOKI-DEPLOY-009)
5. `scripts/deploy-loki.sh` — Deployment script (LOKI-DEPLOY-005)
6. `docs/loki-operations.md` — Operations runbook (LOKI-DEPLOY-008)

**Files to modify:**
1. `Taskfile.yaml` — Add Loki Stack commands (ORGANIZE-002)
2. `docker-compose.loki.yml` — Update volume paths after reorganization (ORGANIZE-001, LOKI-DEPLOY-001)
3. Files will be moved to new directory structure (ORGANIZE-001)

**Files to validate (already exist):**
1. `promtail/promtail-config.yml` — Validate (LOKI-DEPLOY-003)
2. `grafana/provisioning/datasources/loki.yml` — Validate (LOKI-DEPLOY-004)

**Files NOT to modify:**
- Application code (already migrated)
- `plantform/pkg/logger/logger.go` — Logger implementation complete
- `investor/internal/config/env/logger.go` — Config complete

---

## Validation Checklist

Before marking deployment complete:

### Infrastructure
- [ ] All Docker Compose services start without errors
- [ ] Health checks pass for Loki, Promtail, Grafana
- [ ] Services accessible on expected ports

### Log Ingestion
- [ ] Application logs visible in Loki
- [ ] Promtail successfully shipping logs
- [ ] No errors in Promtail logs

### Grafana
- [ ] Loki datasource configured
- [ ] Dashboard displays logs
- [ ] LogQL queries work correctly
- [ ] Derived fields extract `duration_ms`

### Operations
- [ ] Runbook documented
- [ ] Alerts configured
- [ ] Team trained on basic operations

---

## Approval Checklist

Before implementation can begin, verify:

- [ ] All P0 tasks clearly defined with deliverables
- [ ] Dependencies mapped correctly (Phase 0 → Phase 1)
- [ ] Required skills identified
- [ ] Effort estimates reasonable (1.5h for org, 9h for deployment)
- [ ] Risks identified with mitigations
- [ ] No scope creep beyond deployment
- [ ] Architecture unchanged (infrastructure only)
- [ ] Application code not modified
- [ ] YAML organization improves developer experience

---

## Next Steps

**After human approval:**
1. Orchestrator updates STATE.json to `ARCHITECTURE` phase
2. Architecture Agent reviews plan (architecture already complete for logging)
3. Backend Agent implements tasks phase by phase
4. Reviewer validates each phase against acceptance criteria
5. Maximum 3 review cycles before escalation

---

**Status:** APPROVAL_PENDING  
**Human Review Required:** YES  
**Cannot proceed without approval**
