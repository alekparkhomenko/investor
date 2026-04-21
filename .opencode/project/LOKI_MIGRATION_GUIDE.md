# Loki Stack Migration Guide

**For:** Backend Agent  
**Task:** ORGANIZE-001 — Consolidate YAML Configuration Files  
**Time Estimate:** 0.5 hours

---

## Quick Summary

**Decision:** Move all Loki Stack YAML files to `deploy/loki/`

**Why:** Single location for all configs, improved discoverability, follows golang-project-layout best practices.

---

## File Migration Checklist

### Step 1: Create Directory Structure

```bash
mkdir -p deploy/loki/grafana/provisioning/datasources
mkdir -p deploy/loki/grafana/provisioning/dashboards
```

### Step 2: Move Files

```bash
# Move docker-compose file (rename)
mv docker-compose.loki.yml deploy/loki/docker-compose.yml

# Move Loki config
mv loki/loki-config.yml deploy/loki/loki-config.yml

# Move Promtail config
mv promtail/promtail-config.yml deploy/loki/promtail-config.yml

# Move Grafana datasource
mv grafana/provisioning/datasources/loki.yml deploy/loki/grafana/provisioning/datasources/loki.yml

# Remove empty directories (optional, verify first)
rmdir loki promtail
```

### Step 3: Update Volume Paths in docker-compose.yml

**File:** `deploy/loki/docker-compose.yml`

**Change this:**
```yaml
volumes:
  - ./loki:/etc/loki
  - ./promtail:/etc/promtail
  - ./grafana/provisioning:/etc/grafana/provisioning
```

**To this:**
```yaml
volumes:
  - ./loki-config.yml:/etc/loki/loki-config.yml:ro
  - ./promtail-config.yml:/etc/promtail/promtail-config.yml:ro
  - ./grafana/provisioning:/etc/grafana/provisioning:ro
```

### Step 4: Validate

```bash
# From project root
docker-compose -f deploy/loki/docker-compose.yml config

# Should output valid Docker Compose config without errors
```

---

## New Directory Structure

```
deploy/
  loki/
    docker-compose.yml          # Main entry point
    loki-config.yml             # Loki server config
    promtail-config.yml         # Promtail collector config
    grafana/
      provisioning/
        datasources/
          loki.yml              # Loki datasource for Grafana
        dashboards/             # For future dashboards
```

---

## Deployment Commands

**Start Loki Stack:**
```bash
docker-compose -f deploy/loki/docker-compose.yml up -d
```

**Stop Loki Stack:**
```bash
docker-compose -f deploy/loki/docker-compose.yml down
```

**View Logs:**
```bash
docker-compose -f deploy/loki/docker-compose.yml logs -f
```

---

## Acceptance Criteria (from ORGANIZE-001)

- [x] All YAML configs located in single directory (`deploy/loki/`)
- [ ] Directory structure documented in README (Task ORGANIZE-003)
- [x] Docker Compose volume paths updated correctly
- [ ] `docker-compose -f deploy/loki/docker-compose.yml config` validates successfully
- [x] Uses `golang-project-layout` skill for directory structure best practices

---

## Architecture References

- **ARCHITECTURE.md** — Section "Loki Stack Organization" (updated v1.1.0)
- **DECISIONS.md** — ADR-008: Loki Stack Directory Organization
- **EXECUTION_PLAN_LOKI_STACK.md** — Task ORGANIZE-001

---

## Next Steps

After completing this migration:

1. **ORGANIZE-002** — Update Taskfile.yaml with Loki commands
2. **ORGANIZE-003** — Update documentation (README)
3. **LOKI-DEPLOY-001+** — Continue with infrastructure setup

---

**Created:** 2026-04-20  
**Version:** 1.0  
**Status:** Ready for Implementation
