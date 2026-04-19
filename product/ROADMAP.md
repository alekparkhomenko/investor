# Roadmap: Investor

## Phase 1: MVP (Current) ✅
**Goal:** Рабочее приложение, получает котировки MOEX

### Completed
- [x] MOEX ISS API integration
- [x] Graceful shutdown
- [x] .env configuration
- [x] Console output (fmt.Println)

### In Progress
- [ ] Unit tests
- [ ] Structured logging (slog/zap)

---

## Phase 2: Observability
**Goal:** Лучшее понимание работы системы

### Tasks
- [ ] Structured logging (slog)
- [ ] Health check endpoint
- [ ] Metrics (Prometheus)
- [ ] Graceful shutdown logging

---

## Phase 3: Alerts
**Goal:** Уведомления о важных событиях

### Tasks
- [ ] Telegram Bot integration
- [ ] Price change alerts (% threshold)
- [ ] Rate limiting

---

## Phase 4: Reliability
**Goal:** Production-ready

### Tasks
- [ ] Unit tests (>80% coverage)
- [ ] Integration tests
- [ ] Retry logic with backoff
- [ ] Circuit breaker

---

## Phase 5: Storage
**Goal:** Persistence

### Tasks
- [ ] Redis for price history
- [ ] Kafka for events (optional)

---

## Future (Backlog)
- Web UI Dashboard
- WebSocket API
- Docker Compose
- Kubernetes deployment

---
**Version:** 0.1.0 | **Date:** 2026-04-19
