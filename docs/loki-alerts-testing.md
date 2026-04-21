# Testing Grafana Alerts for Loki Stack

This document describes how to test the alert rules provisioned in Grafana.

## Prerequisites

- Loki Stack running: `docker-compose -f deploy/loki/docker-compose.yml up -d`
- Grafana accessible at http://localhost:3000
- Admin credentials: admin/admin

## Test Procedure

### 1. Verify Alert Rules Are Loaded

```bash
# Check Grafana logs for alert provisioning
docker logs grafana 2>&1 | grep -i "alert"

# Expected output should show:
# - "Loaded alert rules from file"
# - No errors about provisioning
```

In Grafana UI:
1. Navigate to **Alerting** → **Alert rules**
2. Verify you see 5 alert rules:
   - ✅ Loki Down
   - ✅ Promtail Not Sending Logs
   - ✅ High Error Rate
   - ✅ Log Ingestion Stopped
   - ✅ High Memory Usage

### 2. Test Each Alert

#### Test 1: Loki Down (Critical)

**Expected**: Alert fires after 1 minute

```bash
# Stop Loki
docker stop loki

# Wait 1-2 minutes

# Check alert status in Grafana UI
# Should show: Loki Down - Firing

# Verify notification sent (email/Telegram)

# Restart Loki
docker start loki
```

**Validation**:
- [ ] Alert fires within 2 minutes
- [ ] Notification received
- [ ] Alert resolves after Loki restarts

---

#### Test 2: Promtail Not Sending Logs (Warning)

**Expected**: Alert fires after 5 minutes

```bash
# Stop Promtail
docker stop promtail

# Wait 5-6 minutes

# Check alert status in Grafana UI
# Should show: Promtail Not Sending Logs - Firing

# Verify notification sent

# Restart Promtail
docker start promtail
```

**Validation**:
- [ ] Alert fires within 6 minutes
- [ ] Notification received
- [ ] Alert resolves after Promtail restarts

---

#### Test 3: High Error Rate (Warning)

**Expected**: Alert fires when error rate > 5% for 5 minutes

**Method A: Generate test error logs**

```bash
# Find investor log file
docker inspect investor | grep LogPath

# Add error logs (example for JSON logs)
for i in {1..50}; do
  echo '{"time":"'$$(date -Iseconds)'","level":"error","app":"investor","component":"test","message":"Test error for alert testing"}' >> /path/to/investor.log
done

# Wait 5-6 minutes

# Check alert status
```

**Method B: Use LogCLI (if available)**

```bash
# Install lokicli
# Push error logs
lokicli push --url=http://localhost:3100 --labels="{app=\"investor\",level=\"error\"}" "Test error"
```

**Validation**:
- [ ] Alert fires when error rate exceeds threshold
- [ ] Notification received
- [ ] Alert resolves when error rate decreases

---

#### Test 4: Log Ingestion Stopped (Critical)

**Expected**: Alert fires after 10 minutes of no logs

```bash
# Stop investor application
docker stop investor

# Wait 10-11 minutes

# Check alert status in Grafana UI
# Should show: Log Ingestion Stopped - Firing

# Verify notification sent

# Restart investor
docker start investor
```

**Validation**:
- [ ] Alert fires within 11 minutes
- [ ] Notification received
- [ ] Alert resolves after investor restarts

---

#### Test 5: High Memory Usage (Warning)

**Expected**: Alert fires when memory > 80% for 5 minutes

**Note**: This alert requires Prometheus datasource for container metrics. For testing, temporarily lower the threshold:

```bash
# Edit alerts.yml, change threshold from 0.8 to 0.01 (1%)
# Restart Grafana
docker-compose restart grafana

# Wait 5-6 minutes

# Check alert status

# Revert threshold back to 0.8
# Restart Grafana
```

**Validation**:
- [ ] Alert fires when threshold exceeded
- [ ] Notification received
- [ ] Alert resolves when memory usage decreases

---

### 3. Test Notification Channels

#### Test Email Notification

```bash
# Ensure SMTP is configured in .env
# In Grafana UI: Alerting → Contact points → default-receiver → Test

# Check email inbox for test message
```

**Validation**:
- [ ] Test email sent successfully
- [ ] Email contains alert details
- [ ] Email formatting correct

---

#### Test Telegram Notification

```bash
# Ensure TELEGRAM_BOT_TOKEN and TELEGRAM_CHAT_ID are set in .env
# In Grafana UI: Alerting → Contact points → default-receiver → Test

# Check Telegram channel/group for test message
```

**Validation**:
- [ ] Test message sent successfully
- [ ] Message contains alert details
- [ ] Message formatting correct (Markdown)

---

### 4. Test Alert Routing

```bash
# Trigger a critical alert (e.g., stop Loki)
docker stop loki

# Measure time to notification
# Expected: Within 30 seconds (group_wait: 10s + evaluation)

# Trigger multiple warning alerts
# Expected: Batched in single notification after 1 minute
```

**Validation**:
- [ ] Critical alerts sent immediately
- [ ] Warning alerts batched together
- [ ] Repeat intervals work correctly

---

### 5. Test Alert Silencing/Muting

```bash
# In Grafana UI: Alerting → Mute timings
# Create a mute timing
# Apply to specific alerts

# Verify alerts are muted during specified time
```

**Validation**:
- [ ] Alerts muted during specified times
- [ ] Alerts resume after mute period ends

---

## Troubleshooting

### Alerts Not Firing

1. **Check alert rule state**:
   - Go to **Alerting** → **Alert rules**
   - Check "State" column
   - Click on alert to see evaluation results

2. **Check query results**:
   - Click on alert rule
   - Go to "Query" tab
   - Run query manually to verify it returns expected results

3. **Check evaluation interval**:
   - Default: 30s
   - Can be adjusted in alert group settings

---

### Notifications Not Sent

1. **Check contact point configuration**:
   - Go to **Alerting** → **Contact points**
   - Click "Test" to verify configuration

2. **Check notification logs**:
   - Go to **Alerting** → **Notification history**
   - Look for failed notifications

3. **Check Grafana logs**:
   ```bash
   docker logs grafana 2>&1 | grep -i "notification"
   ```

---

### Provisioning Errors

1. **Check YAML syntax**:
   ```bash
   python3 -c "import yaml; yaml.safe_load(open('alerts.yml'))"
   ```

2. **Check Grafana logs**:
   ```bash
   docker logs grafana 2>&1 | grep -i "provision"
   ```

3. **Verify file paths**:
   - Ensure files are in correct directory
   - Check file permissions

---

## Success Criteria

All tests pass when:

- ✅ All 5 alert rules loaded without errors
- ✅ Each alert fires when conditions are met
- ✅ Notifications received via configured channels
- ✅ Alert routing works (critical vs warning)
- ✅ Alerts resolve when conditions clear
- ✅ No false positives during normal operation

---

## Automated Testing (Future)

Consider implementing automated alert testing:

```bash
# Example script: scripts/test-alerts.sh
#!/bin/bash

# Test Loki Down alert
echo "Testing Loki Down alert..."
docker stop loki
sleep 90
# Check alert via Grafana API
# Verify notification
docker start loki

# Add more tests...
```

---

**Last Updated**: 2026-04-20  
**Version**: 1.0
