#!/bin/bash
#
# verify-logs.sh - Integration test script for Loki log verification
#
# This script:
# 1. Starts the investor application with Loki enabled
# 2. Waits for logs to be generated
# 3. Queries Loki API to verify log presence
# 4. Verifies required log fields are present
# 5. Cleans up resources
#
# Usage: ./scripts/verify-logs.sh
#
# Requirements:
# - Loki running at http://localhost:3100
# - Go installed for building the application
#

set -euo pipefail

# Configuration
LOKI_HOST="${LOKI_HOST:-http://localhost:3100}"
APP_NAME="${APP_NAME:-investor}"
LOKI_ENV="${LOKI_ENV:-development}"
TEST_DURATION="${TEST_DURATION:-10}"  # seconds to run the app

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

log_info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

log_warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Check if Loki is available
check_loki() {
    log_info "Checking Loki availability at ${LOKI_HOST}..."
    if ! curl -s -o /dev/null -w "%{http_code}" "${LOKI_HOST}/loki/api/v1/labels" | grep -q "200"; then
        log_error "Loki is not available at ${LOKI_HOST}"
        log_error "Please start Loki first or set LOKI_HOST environment variable"
        exit 1
    fi
    log_info "Loki is available"
}

# Start the application
start_app() {
    log_info "Starting investor application with Loki enabled..."
    
    # Build the application
    cd "$(dirname "$0")/.."
    go build -o /tmp/investor-test ./investor/cmd/main.go
    
    # Start the app in background
    LOKI_ENABLED=true \
    LOKI_HOST="${LOKI_HOST}" \
    LOKI_ENV="${LOKI_ENV}" \
    SYMBOLS="SBER" \
    POLL_INTERVAL="2s" \
    /tmp/investor-test &
    
    APP_PID=$!
    log_info "Application started with PID: ${APP_PID}"
    
    # Give it time to generate logs
    sleep "${TEST_DURATION}"
    
    return 0
}

# Stop the application
stop_app() {
    if [ -n "${APP_PID:-}" ] && kill -0 "${APP_PID}" 2>/dev/null; then
        log_info "Stopping application (PID: ${APP_PID})..."
        kill "${APP_PID}" || true
        wait "${APP_PID}" 2>/dev/null || true
        log_info "Application stopped"
    fi
    
    # Cleanup binary
    rm -f /tmp/investor-test
}

# Query Loki for logs
query_loki() {
    local query="$1"
    local description="$2"
    
    log_info "Querying Loki: ${description}"
    log_info "Query: ${query}"
    
    local response
    response=$(curl -s -G "${LOKI_HOST}/loki/api/v1/query_range" \
        --data-urlencode "query=${query}" \
        --data-urlencode "start=$(date -d '5 minutes ago' +%s)000000000" \
        --data-urlencode "end=$(date +%s)000000000")
    
    echo "${response}" | jq '.'
    
    # Check if we got results
    local status
    status=$(echo "${response}" | jq -r '.status')
    
    if [ "${status}" != "success" ]; then
        log_error "Query failed"
        return 1
    fi
    
    local result_type
    result_type=$(echo "${response}" | jq -r '.data.resultType')
    
    if [ "${result_type}" = "streams" ]; then
        local streams_count
        streams_count=$(echo "${response}" | jq '.data.result | length')
        
        if [ "${streams_count}" -gt 0 ]; then
            log_info "✓ Found ${streams_count} log stream(s)"
            return 0
        else
            log_warn "✗ No logs found"
            return 1
        fi
    fi
    
    return 0
}

# Verify log fields
verify_log_fields() {
    log_info "Verifying log field presence..."
    
    # Query for logs with component field
    local query="{app=\"${APP_NAME}\", env=\"${LOKI_ENV}\"}"
    
    local response
    response=$(curl -s -G "${LOKI_HOST}/loki/api/v1/query_range" \
        --data-urlencode "query=${query}" \
        --data-urlencode "start=$(date -d '5 minutes ago' +%s)000000000" \
        --data-urlencode "end=$(date +%s)000000000")
    
    # Check for component field in logs
    if echo "${response}" | jq -e '.data.result[].stream.component' > /dev/null 2>&1; then
        log_info "✓ 'component' field found in logs"
    else
        log_warn "✗ 'component' field not found in logs"
    fi
    
    # Check for duration_ms field (from LOKI-006)
    if echo "${response}" | jq -e '.data.result[].values[] | .[1] | contains("duration_ms")' > /dev/null 2>&1; then
        log_info "✓ 'duration_ms' field found in logs"
    else
        log_warn "✗ 'duration_ms' field not found (may be ok if no fetches completed)"
    fi
    
    # Check for sampled field (from LOKI-007)
    if echo "${response}" | jq -e '.data.result[].values[] | .[1] | contains("sampled")' > /dev/null 2>&1; then
        log_info "✓ 'sampled' field found in logs (log sampling working)"
    else
        log_warn "✗ 'sampled' field not found (may be ok if no 'no quotes' scenario)"
    fi
}

# Main execution
main() {
    log_info "=== Loki Log Verification Script ==="
    log_info ""
    
    # Setup trap for cleanup
    trap stop_app EXIT
    
    # Step 1: Check Loki
    check_loki
    
    # Step 2: Start application
    start_app
    
    # Step 3: Query Loki for application logs
    log_info ""
    log_info "=== Querying Loki for Application Logs ==="
    
    # Query by app name
    query_loki "{app=\"${APP_NAME}\"}" "Logs by app name"
    
    # Query by environment
    query_loki "{app=\"${APP_NAME}\", env=\"${LOKI_ENV}\"}" "Logs by app and environment"
    
    # Query by component
    query_loki "{app=\"${APP_NAME}\", component=\"app\"}" "App component logs"
    query_loki "{app=\"${APP_NAME}\", component=\"moex-ingestor\"}" "MOEX ingestor logs"
    query_loki "{app=\"${APP_NAME}\", component=\"main\"}" "Main component logs"
    
    # Step 4: Verify log fields
    log_info ""
    log_info "=== Verifying Log Fields ==="
    verify_log_fields
    
    # Step 5: Summary
    log_info ""
    log_info "=== Verification Complete ==="
    log_info "Check Grafana or Loki UI for detailed log inspection"
    log_info "Grafana URL: http://localhost:3000/explore"
}

# Run main function
main "$@"
