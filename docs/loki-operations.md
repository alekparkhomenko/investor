# Loki Operations Runbook

This runbook provides operational procedures for the Loki Stack including alerting configuration, troubleshooting, and maintenance.

## Table of Contents

- [Quick Start](#quick-start)
- [Alert Configuration](#alert-configuration)
- [Alert Rules](#alert-rules)
- [Notification Channels](#notification-channels)
- [Testing Alerts](#testing-alerts)
- [Troubleshooting](#troubleshooting)
- [Maintenance](#maintenance)

---

## Quick Start

### Start Loki Stack with Alerting

```bash
# From project root
docker-compose -f deploy/loki/docker-compose.yml up -d

# Verify all services are healthy
docker-compose -f deploy/loki/docker-compose.yml ps

# Check Grafana logs for alert provisioning
docker logs grafana | grep -i alert
```

### Access Grafana

- **URL**: http://localhost:3000
- **Username**: admin
- **Password**: admin

Navigate to **Alerting** → **Alert rules** to view provisioned alerts.

---

## Alert Configuration

### File Structure

```
deploy/loki/grafana/provisioning/
├── alerting/
│   ├── alerts.yml                    # Alert rules
│   └── notification-channels.yml     # Contact points & policies
├── dashboards/
│   ├── dashboard.yml                 # Dashboard provisioning
│   └── investor-logs.json            # Logs dashboard
└── datasources/
    └── loki.yml                      # Loki datasource
```

### Enable/Disable Alerts

To temporarily disable all alerts:

```bash
# Add to docker-compose.yml environment
- GF_ALERTING_ENABLED=false
```

To disable specific alerts, comment them out in `alerts.yml` and restart Grafana.

---

## Alert Rules

### 1. Loki Down (Critical)

**Purpose**: Detect when Loki service is unavailable

**Condition**: `up{job="loki"} == 0` for 1 minute

**Severity**: Critical

**Actions**:
- Immediate notification
- Repeat every 1 hour until resolved

**Runbook**: See [Loki Down](#loki-down)

---

### 2. Promtail Not Sending Logs (Warning)

**Purpose**: Detect when Promtail stops shipping logs

**Condition**: No logs from Promtail for 5 minutes

**Severity**: Warning

**Actions**:
- Batch with other warnings
- Repeat every 4 hours

**Runbook**: See [Promtail Disconnected](#promtail-disconnected)

---

### 3. High Error Rate (Warning)

**Purpose**: Detect elevated error rates in application logs

**Condition**: Error rate > 5% of total logs for 5 minutes

**Severity**: Warning

**Actions**:
- Batch with other warnings
- Repeat every 4 hours

**Runbook**: See [High Error Rate](#high-error-rate)

---

### 4. Log Ingestion Stopped (Critical)

**Purpose**: Detect when application stops sending logs

**Condition**: No logs from investor app for 10 minutes

**Severity**: Critical

**Actions**:
- Immediate notification
- Repeat every 1 hour until resolved

**Runbook**: See [Log Ingestion Stopped](#log-ingestion-stopped)

---

### 5. High Memory Usage (Warning)

**Purpose**: Detect when Loki approaches memory limits

**Condition**: Memory usage > 80% of limit for 5 minutes

**Severity**: Warning

**Actions**:
- Batch with other warnings
- Repeat every 4 hours

**Runbook**: See [High Memory Usage](#high-memory-usage)

---

## Notification Channels

### Email Configuration

To enable email notifications, configure SMTP in `.env`:

```bash
# Email notifications
GF_SMTP_ENABLED=true
GF_SMTP_HOST=smtp.example.com:587
GF_SMTP_USER=grafana@example.com
GF_SMTP_PASSWORD=your_password
GF_SMTP_FROM_ADDRESS=grafana@example.com
```

### Telegram Configuration

To enable Telegram notifications:

1. **Create a bot** via [@BotFather](https://t.me/botfather)
   - Send `/newbot` and follow instructions
   - Save the bot token

2. **Get chat ID**:
   - Add bot to a group/channel
   - Send a message in the chat
   - Get chat ID via: `curl https://api.telegram.org/bot<YOUR_TOKEN>/getUpdates`

3. **Add to .env**:
   ```bash
   TELEGRAM_BOT_TOKEN=123456:ABC-DEF1234ghIkl-zyx57W2v1u123ew11
   TELEGRAM_CHAT_ID=-1001234567890
   ```

4. **Restart Grafana**:
   ```bash
   docker-compose -f deploy/loki/docker-compose.yml restart grafana
   ```

### Notification Policy

- **Critical alerts**: Immediate (10s wait), repeat every 1 hour
- **Warning alerts**: Batched (1m wait), repeat every 4 hours
- **Grouping**: By alertname, severity, and component

---

## Testing Alerts

### Test Loki Down Alert

```bash
# Stop Loki service
docker stop loki

# Wait 1-2 minutes for alert to fire
# Check Grafana UI: Alerting → Alert rules → Loki Down

# Restart Loki
docker start loki
```

### Test Promtail Not Sending Alert

```bash
# Stop Promtail service
docker stop promtail

# Wait 5 minutes for alert to fire
# Check Grafana UI: Alerting → Alert rules → Promtail Not Sending Logs

# Restart Promtail
docker start promtail
```

### Test High Error Rate Alert

```bash
# Generate error logs in investor application
# Method 1: Trigger actual errors in the app
# Method 2: Manually write error logs to log file

# Example: Add to investor app temporarily
log.Error("test error for alert testing", zap.String("component", "test"))

# Wait 5 minutes for alert to fire
# Check Grafana UI: Alerting → Alert rules → High Error Rate
```

### Test Log Ingestion Stopped Alert

```bash
# Stop investor application
docker stop investor

# Wait 10 minutes for alert to fire
# Check Grafana UI: Alerting → Alert rules → Log Ingestion Stopped

# Restart investor application
docker start investor
```

### Test High Memory Alert

```bash
# Note: This is difficult to test in development
# Use Grafana UI to manually trigger test:
# Alerting → Alert rules → High Memory Usage → Test rule

# Or temporarily lower the threshold in alerts.yml to 0.01 (1%)
# for testing purposes
```

### Verify Alert Notifications

1. Check Grafana UI: **Alerting** → **Alert instances**
2. View notification logs: **Alerting** → **Notification history**
3. Check email/Telegram for received notifications

---

## Troubleshooting

### Loki Down

**Symptoms**:
- Alert "Loki Down" is firing
- Cannot query logs in Grafana
- Promtail logs show connection errors

**Diagnosis**:
```bash
# Check Loki service status
docker-compose -f deploy/loki/docker-compose.yml ps loki

# Check Loki logs
docker logs loki

# Test Loki health endpoint
curl http://localhost:3100/ready

# Check resource usage
docker stats loki
```

**Resolution**:
1. **Restart Loki**: `docker-compose restart loki`
2. **Check disk space**: `df -h` (Loki needs space for chunks/index)
3. **Check memory**: Increase memory limit if OOM
4. **Check logs**: Look for errors in Loki logs

**Prevention**:
- Monitor memory usage
- Configure appropriate retention (7 days default)
- Set resource limits in Docker Compose

---

### Promtail Disconnected

**Symptoms**:
- Alert "Promtail Not Sending Logs" is firing
- No new logs appearing in Grafana
- Promtail logs show errors

**Diagnosis**:
```bash
# Check Promtail service status
docker-compose -f deploy/loki/docker-compose.yml ps promtail

# Check Promtail logs
docker logs promtail

# Verify Promtail can reach Loki
docker exec promtail wget -qO- http://loki:3100/ready

# Check log file permissions
ls -la /var/log
```

**Resolution**:
1. **Restart Promtail**: `docker-compose restart promtail`
2. **Check Loki connectivity**: Ensure Loki is running
3. **Verify config**: Check promtail-config.yml for errors
4. **Check permissions**: Ensure Promtail can read log files

**Prevention**:
- Configure health checks
- Set restart policies
- Monitor Promtail metrics

---

### High Error Rate

**Symptoms**:
- Alert "High Error Rate" is firing
- Increased error logs in application

**Diagnosis**:
```bash
# Query error logs in Grafana LogQL
{app="investor", level="error"} | json

# Check error rate trend
sum(rate({app="investor", level="error"}[5m]))

# Compare with total log volume
sum(rate({app="investor"}[5m]))
```

**Resolution**:
1. **Investigate errors**: Check error log messages
2. **Check dependencies**: Database, external APIs
3. **Review recent changes**: Deployments, config changes
4. **Scale if needed**: Add resources if overloaded

**Prevention**:
- Implement proper error handling
- Add circuit breakers
- Monitor error rates continuously

---

### Log Ingestion Stopped

**Symptoms**:
- Alert "Log Ingestion Stopped" is firing
- No logs from investor application

**Diagnosis**:
```bash
# Check if application is running
docker-compose ps investor

# Check application logs
docker logs investor

# Verify log file exists
ls -la /path/to/investor/logs

# Check Promtail is scraping the file
docker logs promtail | grep investor
```

**Resolution**:
1. **Restart application**: `docker-compose restart investor`
2. **Check logging config**: Ensure LOKI_ENABLED=true
3. **Verify log path**: Check promtail-config.yml paths
4. **Check disk space**: Ensure log file can be written

**Prevention**:
- Monitor application health
- Configure log rotation
- Set up application health checks

---

### High Memory Usage

**Symptoms**:
- Alert "High Memory Usage" is firing
- Loki using >80% of memory limit

**Diagnosis**:
```bash
# Check Loki memory usage
docker stats loki

# Query Loki metrics
curl http://localhost:3100/metrics | grep go_memstats

# Check chunk/index size
du -sh /var/lib/loki/*
```

**Resolution**:
1. **Increase memory limit**: Edit docker-compose.yml
2. **Reduce retention**: Lower retention period in loki-config.yml
3. **Compact chunks**: Run compaction manually
4. **Restart Loki**: Clear memory (temporary fix)

**Prevention**:
- Monitor memory trends
- Configure appropriate retention
- Set up compaction
- Use appropriate schema config

---

## Maintenance

### Backup Grafana Dashboards

```bash
# Export dashboards via API
curl -H "Authorization: Bearer <API_KEY>" \
  http://localhost:3000/api/search \
  > grafana-dashboards-backup.json

# Or manually export from UI:
# Dashboard → Settings → Export → Save to file
```

### Update Alert Rules

1. Edit `deploy/loki/grafana/provisioning/alerting/alerts.yml`
2. Validate YAML syntax: `yamllint alerts.yml`
3. Restart Grafana: `docker-compose restart grafana`
4. Verify in UI: **Alerting** → **Alert rules**

### Rotate Logs

Loki automatically handles log retention based on configuration. To manually clean old data:

```bash
# Stop Loki
docker stop loki

# Remove old chunks (be careful!)
# Only do this if you understand the implications
rm -rf /var/lib/loki/chunks/*

# Restart Loki
docker start loki
```

### Update Grafana

```bash
# Update docker-compose.yml with new version
# Pull new image
docker-compose pull grafana

# Restart Grafana
docker-compose restart grafana

# Verify provisioning still works
docker logs grafana | grep -i provision
```

---

## Monitoring the Monitors

### Grafana Health Check

```bash
curl http://localhost:3000/api/health
# Expected: {"commit":"...","database":"ok","version":"10.2.0"}
```

### Alert Rule Validation

```bash
# Check if alert rules are loaded
docker logs grafana | grep -i "loaded alert"

# Check for provisioning errors
docker logs grafana | grep -i "error.*alert"
```

### Notification Channel Test

In Grafana UI:
1. Go to **Alerting** → **Contact points**
2. Select a contact point
3. Click **Test**
4. Verify notification received

---

## Appendix: LogQL Quick Reference

### Basic Queries

```logql
# All logs from investor app
{app="investor"}

# Error logs only
{app="investor", level="error"}

# Logs from specific component
{app="investor", component="ingestor"}

# Logs containing specific text
{app="investor"} |= "timeout"

# Logs with JSON field extraction
{app="investor"} | json | duration_ms > 100
```

### Metric Queries

```logql
# Log rate (logs per second)
rate({app="investor"}[5m])

# Error rate
sum(rate({app="investor", level="error"}[5m]))

# Error ratio
sum(rate({app="investor", level="error"}[5m])) / sum(rate({app="investor"}[5m]))

# Log volume histogram
sum by (level) (rate({app="investor"}[5m]))
```

---

## Support

For issues not covered in this runbook:

1. Check Grafana logs: `docker logs grafana`
2. Check Loki logs: `docker logs loki`
3. Check Promtail logs: `docker logs promtail`
4. Review Grafana documentation: https://grafana.com/docs/
5. Check Loki documentation: https://grafana.com/docs/loki/

---

**Last Updated**: 2026-04-20  
**Version**: 1.0  
**Maintained By**: Backend Team
