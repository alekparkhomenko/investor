# Product Backlog

## P0 - Must Have (MVP Critical)

### 1. Project Organization
**SP:** 2 | **Deps:** None

**Description:**
Consolidate all YAML configuration files in a single dedicated directory and provide quick-access commands in Taskfile for deployment operations.

**Acceptance:**
- All YAML configs located in one directory (not scattered across project)
- Taskfile.yaml contains deployment commands
- Setup documentation exists
- Easy to discover and maintain

**Skills:** project-organization, developer-experience

---

### 2. Product Definition System
**SP:** 5 | **Deps:** None

**Description:**
Create isolated Product Owner agent generating product artifacts.

**Acceptance:**
- Generates PRODUCT_VISION.md
- Generates ROADMAP.md
- Generates MVP_SPEC.md
- Generates BACKLOG.md
- No technical details

**Skills:** product-management, requirements-analysis

---

### 3. Execution Planning System
**SP:** 3 | **Deps:** Product Definition

**Description:**
Convert backlog to executable plan with human approval.

**Acceptance:**
- Creates EXECUTION_PLAN.md
- Decomposes features
- Maps dependencies
- Stops for approval

**Skills:** project-planning, task-decomposition

---

### 4. State Management System
**SP:** 5 | **Deps:** None

**Description:**
Orchestrator manages STATE.json as single source of truth.

**Acceptance:**
- STATE.json schema defined
- Tracks phases
- Tracks tasks
- Only Orchestrator writes

**Skills:** state-management, workflow-orchestration

---

### 5. Architecture Design System
**SP:** 5 | **Deps:** Planning (approved)

**Description:**
Design system based on approved execution plan.

**Acceptance:**
- Creates ARCHITECTURE.md
- Component diagram
- Data models
- API specs
- Maps skills

**Skills:** system-design, api-design

---

### 6. Code Implementation System
**SP:** 8 | **Deps:** Architecture

**Description:**
Backend implements code following plan and architecture.

**Acceptance:**
- Implements tasks
- Follows architecture
- Uses skills
- Test coverage >= 80%
- No scope expansion

**Skills:** golang-implementation, golang-testing

---

### 7. Quality Review System
**SP:** 5 | **Deps:** Implementation

**Description:**
Reviewer validates code against plan and rules.

**Acceptance:**
- Checks compliance
- Checks skills
- Clear feedback
- Max 3 cycles
- Escalation support

**Skills:** code-review, quality-assurance

---

### 8. Skills Management System
**SP:** 3 | **Deps:** None

**Description:**
Reusable skills registry for engineering agents.

**Acceptance:**
- SKILLS_INDEX.md
- 10+ core skills
- Reuse enforcement
- No duplication

**Skills:** technical-writing, knowledge-management

---

---

## P2 - Future Enhancements

### 9. Telegram Alerting with Grafana
**SP:** 5 | **Deps:** None (See GitHub Issue #14)

**Description:**
Integrate Grafana alerting with Telegram bot for real-time notifications.

**Acceptance:**
- Telegram bot receives alerts from Grafana
- Alert rules for critical metrics (health, errors, latency)
- Notification routing by severity

**Skills:** grafana-alerting, telegram-integration, observability

---

### 10. Custom User Alerts
**SP:** 8 | **Deps:** None

**Description:**
User-defined alert rules with flexible notification channels. Users can create personal alert conditions (price thresholds, percent changes) and receive notifications via Telegram, Email, or Webhook.

**Acceptance:**
- API for CRUD alert rules
- Rule conditions: price above/below, percent change, volume spike
- Notification channels: Telegram, Email, Webhook
- Alert history logging
- Per-user rule management

**Skills:** golang-api-design, database-design, notification-systems

---

### 11. User Portfolio Selection
**SP:** 8 | **Deps:** None

**Description:**
CLI interface for selecting and managing personal ticker portfolio. Users can browse available MOEX tickers, add/remove them to their portfolio, and receive price updates for selected stocks.

**Acceptance:**
- Command: `investor ticker list` - show available MOEX tickers
- Command: `investor portfolio` - show user's portfolio
- Command: `investor portfolio add SBER GAZP` - add tickers
- Command: `investor portfolio remove TATN` - remove tickers
- PostgreSQL storage for portfolios
- Integration with MOEX ingestor

**Skills:** golang-cli, database-design, golang-db-patterns

---

## Summary

| Priority | Count | Story Points |
|----------|-------|--------------|
| P0 | 8 | 36 |
| P2 | 3 | 21 |
| **Total** | **11** | **57** |

**Estimated Duration:** 4-6 weeks (P0)

---

**Version:** 1.4 | **Last Updated:** 2026-04-20
