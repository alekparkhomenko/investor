# Loki Stack Operations Runbook

**Версия:** 1.0  
**Последнее обновление:** 2026-04-20  
**Статус:** Production Ready  
**Ответственная команда:** Platform Engineering

---

## 📋 Содержание

1. [Quick Start](#1-quick-start)
2. [Common Operations](#2-common-operations)
3. [Troubleshooting](#3-troubleshooting)
4. [Monitoring & Alerts](#4-monitoring--alerts)
5. [Backup & Recovery](#5-backup--recovery)
6. [Performance Tuning](#6-performance-tuning)

---

## 1. Quick Start

### 1.1. Как запустить Loki Stack

#### Быстрый запуск (рекомендуется)

```bash
cd /Users/alekparkhomenko/Projects/investor

# Запуск всех сервисов
task loki:up
```

#### Ручной запуск через Docker Compose

```bash
# Запуск в фоновом режиме
docker-compose -f deploy/loki/docker-compose.yml up -d

# Запуск в режиме foreground (для отладки)
docker-compose -f deploy/loki/docker-compose.yml up
```

#### Проверка успешного запуска

```bash
# Проверка статуса сервисов
task loki:status

# Ожидаемый вывод:
# NAME       IMAGE                        STATUS
# loki       grafana/loki:2.9.0           Up
# promtail   grafana/promtail:2.9.0       Up
# grafana    grafana/grafana:10.2.0       Up
```

#### Проверка доступности сервисов

```bash
# Loki health check
curl http://localhost:3100/ready
# Ответ: "ready"

# Loki status API
curl http://localhost:3100/loki/api/v1/status

# Grafana health check
curl http://localhost:3000/api/health
# Ответ: {"commit":"...","database":"ok","version":"10.2.0"}

# Promtail не имеет HTTP API, проверяем через логи
docker-compose -f deploy/loki/docker-compose.yml logs promtail | tail -20
```

### 1.2. Как остановить Loki Stack

#### Штатная остановка

```bash
# Мягкая остановка (сохранение данных)
task loki:down

# Или вручную
docker-compose -f deploy/loki/docker-compose.yml down
```

#### Аварийная остановка

```bash
# Принудительная остановка
docker-compose -f deploy/loki/docker-compose.yml down --force

# Остановить конкретный сервис
docker-compose -f deploy/loki/docker-compose.yml stop loki
```

#### Остановка с очисткой данных

```bash
# ⚠️ ВНИМАНИЕ: Удаляет все данные Loki!
task loki:prune

# Или вручную
docker-compose -f deploy/loki/docker-compose.yml down -v
```

### 1.3. Как проверить статус

#### Комплексная проверка статуса

```bash
#!/bin/bash
# Скрипт проверки статуса Loki Stack

echo "🔍 Проверка статуса Loki Stack..."
echo ""

# Проверка Docker контейнеров
echo "📦 Docker контейнеры:"
docker-compose -f deploy/loki/docker-compose.yml ps
echo ""

# Проверка Loki
echo "🟢 Loki API:"
if curl -s http://localhost:3100/ready | grep -q "ready"; then
    echo "   ✅ Loki готов"
else
    echo "   ❌ Loki не отвечает"
fi
echo ""

# Проверка Grafana
echo "📊 Grafana API:"
if curl -s http://localhost:3000/api/health | grep -q "ok"; then
    echo "   ✅ Grafana готов"
else
    echo "   ❌ Grafana не отвечает"
fi
echo ""

# Проверка логов на ошибки
echo "⚠️ Последние ошибки в логах:"
docker-compose -f deploy/loki/docker-compose.yml logs --tail=50 | grep -i error | tail -5
```

#### Быстрые команды проверки

```bash
# Статус всех сервисов
task loki:status

# Детальный статус с информацией о портах
docker-compose -f deploy/loki/docker-compose.yml ps -a

# Проверка использования ресурсов
docker stats loki promtail grafana --no-stream
```

---

## 2. Common Operations

### 2.1. Просмотр логов

#### Просмотр логов всех сервисов

```bash
# В реальном времени
task loki:logs

# Последние 100 строк
docker-compose -f deploy/loki/docker-compose.yml logs --tail=100

# С временными метками
docker-compose -f deploy/loki/docker-compose.yml logs -t
```

#### Просмотр логов конкретного сервиса

```bash
# Loki логи
docker-compose -f deploy/loki/docker-compose.yml logs loki

# Promtail логи
docker-compose -f deploy/loki/docker-compose.yml logs promtail

# Grafana логи
docker-compose -f deploy/loki/docker-compose.yml logs grafana
```

#### Фильтрация логов

```bash
# Только ошибки
docker-compose -f deploy/loki/docker-compose.yml logs | grep -i error

# Логи за последние 5 минут
docker-compose -f deploy/loki/docker-compose.yml logs --since=5m

# Логи с конкретной меткой
docker-compose -f deploy/loki/docker-compose.yml logs loki | grep "component=ingester"
```

#### Поиск логов приложения в Loki

```bash
# Через API Loki
curl -G "http://localhost:3100/loki/api/v1/query_range" \
  --data-urlencode 'query={app="investor"}' \
  --data-urlencode 'start=2026-04-20T00:00:00Z' \
  --data-urlencode 'end=2026-04-20T23:59:59Z' \
  --data-urlencode 'limit=100' | jq .

# Через LogCLI (если установлен)
logcli query '{app="investor"}' --since=1h
```

### 2.2. Перезапуск сервисов

#### Перезапуск всех сервисов

```bash
# Мягкий перезапуск
task loki:down && task loki:up

# Быстрый перезапуск
docker-compose -f deploy/loki/docker-compose.yml restart

# Перезапуск с пересозданием контейнеров
docker-compose -f deploy/loki/docker-compose.yml up -d --force-recreate
```

#### Перезапуск конкретного сервиса

```bash
# Перезапуск Loki
docker-compose -f deploy/loki/docker-compose.yml restart loki

# Перезапуск Promtail
docker-compose -f deploy/loki/docker-compose.yml restart promtail

# Перезапуск Grafana
docker-compose -f deploy/loki/docker-compose.yml restart grafana
```

#### Перезапуск с обновлением конфигурации

```bash
# Обновить конфигурацию Loki
docker-compose -f deploy/loki/docker-compose.yml up -d --force-recreate loki

# Обновить конфигурацию Promtail
docker-compose -f deploy/loki/docker-compose.yml up -d --force-recreate promtail

# Обновить конфигурацию Grafana
docker-compose -f deploy/loki/docker-compose.yml up -d --force-recreate grafana
```

### 2.3. Очистка данных (Prune)

#### Очистка старых логов

```bash
# Очистка данных старше 7 дней (настраивается в loki-config.yml)
# Автоматически выполняется compactor'ом Loki

# Принудительная очистка через API
curl -X POST http://localhost:3100/loki/api/v1/index/delete \
  -H "Content-Type: application/json" \
  -d '{
    "matchers": "{app=\"investor\"}",
    "start": "2026-04-01T00:00:00Z",
    "end": "2026-04-10T00:00:00Z"
  }'
```

#### Полная очистка данных

```bash
# ⚠️ ВНИМАНИЕ: Удаляет ВСЕ данные!
task loki:prune

# Ручная очистка volumes
docker-compose -f deploy/loki/docker-compose.yml down -v

# Очистка конкретных volumes
docker volume rm loki_loki-data
docker volume rm loki_grafana-data
```

#### Очистка кэша Promtail

```bash
# Очистка позиций Promtail (будет читать файлы заново)
docker-compose -f deploy/loki/docker-compose.yml down -v promtail

# Или вручную удалить файл позиций
docker-compose -f deploy/loki/docker-compose.yml exec promtail \
  rm /var/lib/promtail/positions.yaml
```

#### Очистка старых dashboards в Grafana

```bash
# Через API Grafana
curl -X DELETE http://admin:admin@localhost:3000/api/dashboards/uid/<dashboard-uid>

# Через поиск и удаление
curl -H "Authorization: Bearer <token>" \
  http://localhost:3000/api/search?type=dash-db
```

---

## 3. Troubleshooting

### 3.1. Loki не запускается

#### Симптомы
- Контейнер loki exits immediately
- В логах ошибки конфигурации
- Port 3100 не слушается

#### Диагностика

```bash
# Проверить логи контейнера
docker-compose -f deploy/loki/docker-compose.yml logs loki

# Проверить статус контейнера
docker-compose -f deploy/loki/docker-compose.yml ps loki

# Проверить использование памяти
docker stats loki --no-stream

# Проверить занятость порта
lsof -i :3100
```

#### Решения

**Проблема 1: Ошибка конфигурации**

```bash
# Валидировать конфигурацию
docker run --rm -v $(pwd)/deploy/loki/loki-config.yml:/loki-config.yml \
  grafana/loki:2.9.0 -config.file=/loki-config.yml -verify-config

# Исправить ошибки в loki-config.yml
# Перезапустить Loki
docker-compose -f deploy/loki/docker-compose.yml up -d --force-recreate loki
```

**Проблема 2: Порт 3100 занят**

```bash
# Найти процесс, занимающий порт
lsof -i :3100

# Остановить конфликтующий процесс
kill -9 <PID>

# Или изменить порт в docker-compose.yml
# ports:
#   - "3101:3100"  # Использовать порт 3101
```

**Проблема 3: Недостаточно памяти**

```bash
# Увеличить лимит памяти в docker-compose.yml
# deploy:
#   resources:
#     limits:
#       memory: 1G  # Было 512M

# Перезапустить с новыми лимитами
docker-compose -f deploy/loki/docker-compose.yml up -d --force-recreate loki
```

**Проблема 4: Повреждённые данные**

```bash
# ⚠️ Остановка и очистка данных
docker-compose -f deploy/loki/docker-compose.yml down
docker volume rm loki_loki-data
docker-compose -f deploy/loki/docker-compose.yml up -d
```

### 3.2. Grafana не подключается к Loki

#### Симптомы
- В Grafana ошибка "Data source unavailable"
- Loki datasource не проходит health check
- Запросы в Explore возвращают ошибку

#### Диагностика

```bash
# Проверить доступность Loki
curl http://localhost:3100/ready

# Проверить сеть Docker
docker network inspect loki_default

# Проверить логи Grafana
docker-compose -f deploy/loki/docker-compose.yml logs grafana | grep -i "loki\|datasource"

# Проверить конфигурацию datasource
docker-compose -f deploy/loki/docker-compose.yml exec grafana \
  cat /etc/grafana/provisioning/datasources/loki.yml
```

#### Решения

**Проблема 1: Неправильный URL datasource**

```bash
# Проверить URL в datasource
# Должно быть: http://loki:3100 (внутри Docker сети)
# Или: http://localhost:3100 (для внешнего доступа)

# Исправить в deploy/loki/grafana/provisioning/datasources/loki.yml
# url: http://loki:3100

# Перезапустить Grafana
docker-compose -f deploy/loki/docker-compose.yml restart grafana
```

**Проблема 2: Сетевые проблемы**

```bash
# Проверить, что сервисы в одной сети
docker network inspect loki_default

# Если нет - добавить в docker-compose.yml
# networks:
#   - loki_default

# Пересоздать сеть
docker-compose -f deploy/loki/docker-compose.yml down
docker network rm loki_default
docker-compose -f deploy/loki/docker-compose.yml up -d
```

**Проблема 3: Loki ещё не готов**

```bash
# Подождать готовности Loki
sleep 10
curl http://localhost:3100/ready

# Перепроверить datasource в Grafana UI
# Configuration → Data sources → Loki → Save & test
```

### 3.3. Promtail не отправляет логи

#### Симптомы
- Логи приложения не появляются в Loki
- В логах Promtail ошибки подключения
- Статус targets показывает "DOWN"

#### Диагностика

```bash
# Проверить логи Promtail
docker-compose -f deploy/loki/docker-compose.yml logs promtail

# Проверить конфигурацию Promtail
docker-compose -f deploy/loki/docker-compose.yml exec promtail \
  cat /etc/promtail/config.yml

# Проверить доступность Loki из Promtail
docker-compose -f deploy/loki/docker-compose.yml exec promtail \
  curl http://loki:3100/ready

# Проверить, что файлы логов существуют
docker-compose -f deploy/loki/docker-compose.yml exec promtail \
  ls -la /var/log/investor/
```

#### Решения

**Проблема 1: Неправильный URL Loki**

```bash
# Проверить в promtail-config.yml
# clients:
#   - url: http://loki:3100/loki/api/v1/push

# Исправить и перезапустить
docker-compose -f deploy/loki/docker-compose.yml up -d --force-recreate promtail
```

**Проблема 2: Promtail не видит файлы логов**

```bash
# Проверить volume mounts в docker-compose.yml
# volumes:
#   - /var/log/investor:/var/log/investor:ro
#   - /var/run/docker.sock:/var/run/docker.sock

# Убедиться, что файлы существуют
ls -la /var/log/investor/

# Проверить права доступа
chmod 644 /var/log/investor/*.log
```

**Проблема 3: Ошибки парсинга JSON**

```bash
# Проверить формат логов
tail -5 /var/log/investor/app.log

# Исправить pipeline stages в promtail-config.yml
# pipeline_stages:
#   - json:
#       expressions:
#         level: level
#         component: component

# Перезапустить Promtail
docker-compose -f deploy/loki/docker-compose.yml restart promtail
```

### 3.4. Высокое потребление памяти

#### Симптомы
- Loki использует >1GB памяти
- Частые OOM kills
- Замедление запросов

#### Диагностика

```bash
# Проверить использование памяти
docker stats loki --no-stream

# Проверить метрики Loki
curl http://localhost:3100/metrics | grep go_memstats

# Проверить размер данных
du -sh $(docker volume inspect loki_loki-data | jq -r '.[0].Mountpoint')

# Проверить количество streams
curl -G http://localhost:3100/loki/api/v1/index/stats \
  --data-urlencode 'match={app="investor"}'
```

#### Решения

**Решение 1: Настроить retention**

```yaml
# loki-config.yml
limits_config:
  retention_period: 168h  # 7 дней

compactor:
  working_directory: /loki/compactor
  shared_store: filesystem
  compaction_interval: 10m
  retention_enabled: true
  retention_delete_delay: 2h
  retention_delete_worker_count: 150
```

**Решение 2: Ограничить ingestion rate**

```yaml
# loki-config.yml
limits_config:
  ingestion_rate_mb: 16
  ingestion_burst_size_mb: 24
  max_streams_per_user: 10000
  max_line_size: 256kb
```

**Решение 3: Увеличить ресурсы**

```yaml
# docker-compose.yml
deploy:
  resources:
    limits:
      cpus: '2'
      memory: 2G
    reservations:
      cpus: '1'
      memory: 1G
```

**Решение 4: Очистить старые данные**

```bash
# Принудительная compaction
curl -X POST http://localhost:3100/loki/api/v1/compactor/compact

# Удалить старые streams
curl -X DELETE "http://localhost:3100/loki/api/v1/index/delete" \
  -H "Content-Type: application/json" \
  -d '{"matchers": "{app=\"investor\"}", "start": "2026-04-01T00:00:00Z", "end": "2026-04-10T00:00:00Z"}'
```

---

## 4. Monitoring & Alerts

### 4.1. Health Check сервисов

#### Loki Health Endpoints

```bash
# Ready check
curl http://localhost:3100/ready
# Ответ: "ready" или "not ready"

# Status API
curl http://localhost:3100/loki/api/v1/status
# Ответ: {"status":"success","data":{...}}

# Metrics
curl http://localhost:3100/metrics
```

#### Grafana Health Endpoints

```bash
# Health check
curl http://localhost:3000/api/health
# Ответ: {"commit":"...","database":"ok","version":"10.2.0"}

# API status
curl http://localhost:3000/api/status

# Metrics
curl http://localhost:3000/metrics
```

#### Promtail Health Check

```bash
# Promtail не имеет HTTP API, проверяем через логи
docker-compose -f deploy/loki/docker-compose.yml logs promtail | tail -10

# Проверка через positions file
docker-compose -f deploy/loki/docker-compose.yml exec promtail \
  cat /var/lib/promtail/positions.yaml
```

#### Скрипт комплексной проверки

```bash
#!/bin/bash
# check-loki-health.sh

set -e

LOKI_URL="http://localhost:3100"
GRAFANA_URL="http://localhost:3000"

echo "🔍 Loki Stack Health Check"
echo "==========================="
echo ""

# Loki
echo "🟢 Loki:"
if curl -s $LOKI_URL/ready | grep -q "ready"; then
    echo "   ✅ Ready"
else
    echo "   ❌ Not ready"
    exit 1
fi

# Grafana
echo "📊 Grafana:"
if curl -s $GRAFANA_URL/api/health | grep -q "ok"; then
    echo "   ✅ Healthy"
else
    echo "   ❌ Unhealthy"
    exit 1
fi

# Promtail
echo "📤 Promtail:"
if docker-compose -f deploy/loki/docker-compose.yml ps promtail | grep -q "Up"; then
    echo "   ✅ Running"
else
    echo "   ❌ Not running"
    exit 1
fi

echo ""
echo "✅ All services healthy!"
```

### 4.2. Метрики для мониторинга

#### Ключевые метрики Loki

```promql
# Количество запросов в секунду
rate(loki_request_duration_seconds_count[5m])

# 99-й перцентиль длительности запросов
histogram_quantile(0.99, rate(loki_request_duration_seconds_bucket[5m]))

# Количество ошибок
rate(loki_request_duration_seconds_count{status_code=~"5.."}[5m])

# Использование памяти
go_memstats_heap_inuse_bytes{job="loki"}

# Количество активных streams
loki_ingester_memory_streams

# Размер хранилища
loki_ingester_chunks_stored_total
```

#### Ключевые метрики Promtail

```promql
# Количество отправленных строк
rate(promtail_sent_entries_total[5m])

# Количество ошибок отправки
rate(promtail_pusher_errors_total[5m])

# Задержка отправки
promtail_custom_labels_count

# Количество читаемых файлов
promtail_read_files
```

#### Ключевые метрики Grafana

```promql
# Количество запросов к datasource
rate(grafana_datasource_request_total{datasource="Loki"}[5m])

# Ошибки запросов
rate(grafana_datasource_request_error_total{datasource="Loki"}[5m])

# Время отклика
histogram_quantile(0.95, rate(grafana_datasource_request_duration_seconds_bucket[5m]))
```

### 4.3. Рекомендуемые алерты

#### Alert: Loki недоступен

```yaml
# grafana/provisioning/alerting/loki-alerts.yml
apiVersion: 1
groups:
  - orgId: 1
    name: Loki Stack
    folder: Loki
    interval: 1m
    rules:
      - uid: loki-down
        title: "Loki Down"
        condition: C
        data:
          - refId: A
            relativeTimeRange:
              from: 60
              to: 0
            datasourceUid: prometheus
            model:
              expr: up{job="loki"} == 0
              intervalMs: 1000
              maxDataPoints: 43200
          - refId: C
            relativeTimeRange:
              from: 0
              to: 0
            datasourceUid: __expr__
            model:
              expression: C
              type: math
        noDataState: Alerting
        execErrState: Alerting
        for: 1m
        annotations:
          summary: "Loki сервис недоступен"
          description: "Loki не отвечает более 1 минуты"
        labels:
          severity: critical
        isPaused: false
```

#### Alert: Promtail не отправляет логи

```yaml
      - uid: promtail-not-sending
        title: "Promtail Not Sending Logs"
        condition: C
        data:
          - refId: A
            relativeTimeRange:
              from: 300
              to: 0
            datasourceUid: prometheus
            model:
              expr: rate(promtail_sent_entries_total[5m]) == 0
              intervalMs: 1000
          - refId: C
            relativeTimeRange:
              from: 0
              to: 0
            datasourceUid: __expr__
            model:
              expression: C
              type: math
        noDataState: Alerting
        execErrState: Alerting
        for: 5m
        annotations:
          summary: "Promtail не отправляет логи"
          description: "Promtail не отправлял логи в течение 5 минут"
        labels:
          severity: warning
        isPaused: false
```

#### Alert: Высокий уровень ошибок в логах

```yaml
      - uid: high-error-rate
        title: "High Error Rate in Application Logs"
        condition: C
        data:
          - refId: A
            relativeTimeRange:
              from: 300
              to: 0
            datasourceUid: loki
            model:
              expr: sum(rate({app="investor"} |= "ERROR" [5m]))
              intervalMs: 1000
          - refId: C
            relativeTimeRange:
              from: 0
              to: 0
            datasourceUid: __expr__
            model:
              expression: "$A > 10"
              type: math
        noDataState: NoData
        execErrState: Alerting
        for: 5m
        annotations:
          summary: "Высокий уровень ошибок в приложении"
          description: "Более 10 ошибок в секунду в течение 5 минут"
        labels:
          severity: warning
        isPaused: false
```

#### Alert: Остановка_ingestion

```yaml
      - uid: ingestion-stopped
        title: "Log Ingestion Stopped"
        condition: C
        data:
          - refId: A
            relativeTimeRange:
              from: 600
              to: 0
            datasourceUid: loki
            model:
              expr: count_over_time({app="investor"}[10m])
              intervalMs: 1000
          - refId: C
            relativeTimeRange:
              from: 0
              to: 0
            datasourceUid: __expr__
            model:
              expression: "$A == 0"
              type: math
        noDataState: Alerting
        execErrState: Alerting
        for: 10m
        annotations:
          summary: "Поступление логов остановилось"
          description: "Не получено ни одного лога за 10 минут"
        labels:
          severity: critical
        isPaused: false
```

#### Alert: Высокое использование памяти Loki

```yaml
      - uid: loki-high-memory
        title: "Loki High Memory Usage"
        condition: C
        data:
          - refId: A
            relativeTimeRange:
              from: 300
              to: 0
            datasourceUid: prometheus
            model:
              expr: go_memstats_heap_inuse_bytes{job="loki"} / 1024 / 1024 / 1024
              intervalMs: 1000
          - refId: C
            relativeTimeRange:
              from: 0
              to: 0
            datasourceUid: __expr__
            model:
              expression: "$A > 1.5"
              type: math
        noDataState: NoData
        execErrState: Alerting
        for: 5m
        annotations:
          summary: "Высокое использование памяти Loki"
          description: "Loki использует более 1.5GB памяти"
        labels:
          severity: warning
        isPaused: false
```

---

## 5. Backup & Recovery

### 5.1. Как сделать backup данных Loki

#### Backup через volume snapshot

```bash
# Остановить Loki для консистентности
docker-compose -f deploy/loki/docker-compose.yml stop loki

# Создать backup volume
docker run --rm \
  -v loki_loki-data:/data:ro \
  -v $(pwd)/backups:/backup \
  alpine tar czf /backup/loki-data-$(date +%Y%m%d-%H%M%S).tar.gz -C /data .

# Запустить Loki
docker-compose -f deploy/loki/docker-compose.yml start loki

echo "✅ Backup создан в $(pwd)/backups/"
```

#### Backup через rsync

```bash
# Найти mount point volume
MOUNT_POINT=$(docker volume inspect loki_loki-data | jq -r '.[0].Mountpoint')

# Создать backup
rsync -avz --delete $MOUNT_POINT/ /backup/loki-data-$(date +%Y%m%d-%H%M%S)/

echo "✅ Backup создан через rsync"
```

#### Backup конфигураций

```bash
# Создать директорию для backup
mkdir -p backups/config/$(date +%Y%m%d-%H%M%S)

# Скопировать конфигурации
cp deploy/loki/loki-config.yml backups/config/$(date +%Y%m%d-%H%M%S)/
cp deploy/loki/promtail-config.yml backups/config/$(date +%Y%m%d-%H%M%S)/
cp -r deploy/loki/grafana backups/config/$(date +%Y%m%d-%H%M%S)/

# Скопировать dashboards из Grafana
curl -H "Authorization: Bearer admin:admin" \
  http://localhost:3000/api/search?type=dash-db | \
  jq '.[].uid' | while read uid; do
    curl -H "Authorization: Bearer admin:admin" \
      http://localhost:3000/api/dashboards/uid/$uid \
      > backups/config/$(date +%Y%m%d-%H%M%S)/dashboard-$uid.json
done

echo "✅ Конфигурации и dashboards забэкаплены"
```

#### Автоматический backup (cron)

```bash
# /etc/cron.d/loki-backup
0 2 * * * root /usr/local/bin/loki-backup.sh >> /var/log/loki-backup.log 2>&1
```

```bash
#!/bin/bash
# /usr/local/bin/loki-backup.sh

set -e

BACKUP_DIR="/backup/loki"
DATE=$(date +%Y%m%d-%H%M%S)
RETENTION_DAYS=7

# Создать backup
cd /opt/investor
docker-compose -f deploy/loki/docker-compose.yml stop loki
docker run --rm \
  -v loki_loki-data:/data:ro \
  -v $BACKUP_DIR:/backup \
  alpine tar czf /backup/loki-data-$DATE.tar.gz -C /data .
docker-compose -f deploy/loki/docker-compose.yml start loki

# Удалить старые backup'ы
find $BACKUP_DIR -name "loki-data-*.tar.gz" -mtime +$RETENTION_DAYS -delete

echo "[$DATE] Backup completed"
```

### 5.2. Как восстановить из backup

#### Восстановление volume

```bash
# Остановить Loki
docker-compose -f deploy/loki/docker-compose.yml stop loki

# Очистить текущие данные
docker volume rm loki_loki-data
docker volume create loki_loki-data

# Восстановить из backup
BACKUP_FILE="backups/loki-data-20260420-020000.tar.gz"
docker run --rm \
  -v loki_loki-data:/data \
  -v $(pwd)/$BACKUP_FILE:/backup.tar.gz:ro \
  alpine tar xzf /backup.tar.gz -C /data

# Запустить Loki
docker-compose -f deploy/loki/docker-compose.yml start loki

echo "✅ Данные восстановлены"
```

#### Восстановление конфигураций

```bash
# Восстановить конфигурации
BACKUP_DATE="20260420-020000"
cp backups/config/$BACKUP_DATE/loki-config.yml deploy/loki/
cp backups/config/$BACKUP_DATE/promtail-config.yml deploy/loki/
cp -r backups/config/$BACKUP_DATE/grafana/* deploy/loki/grafana/

# Перезапустить сервисы
docker-compose -f deploy/loki/docker-compose.yml up -d --force-recreate

echo "✅ Конфигурации восстановлены"
```

#### Восстановление dashboards

```bash
# Восстановить dashboards через API
for file in backups/config/$BACKUP_DATE/dashboard-*.json; do
    uid=$(basename $file | sed 's/dashboard-//' | sed 's/.json//')
    
    curl -X POST \
      -H "Content-Type: application/json" \
      -H "Authorization: Bearer admin:admin" \
      -d @$file \
      http://localhost:3000/api/dashboards/db
done

echo "✅ Dashboards восстановлены"
```

#### Полное восстановление

```bash
#!/bin/bash
# restore-loki.sh

set -e

BACKUP_DATE=$1

if [ -z "$BACKUP_DATE" ]; then
    echo "Использование: $0 <backup-date>"
    echo "Пример: $0 20260420-020000"
    exit 1
fi

echo "🔄 Восстановление Loki из backup $BACKUP_DATE..."

# Остановить сервисы
docker-compose -f deploy/loki/docker-compose.yml down

# Восстановить volume
docker volume rm loki_loki-data
docker run --rm \
  -v loki_loki-data:/data \
  -v $(pwd)/backups/loki-data-$BACKUP_DATE.tar.gz:/backup.tar.gz:ro \
  alpine tar xzf /backup.tar.gz -C /data

# Восстановить конфигурации
cp backups/config/$BACKUP_DATE/loki-config.yml deploy/loki/
cp backups/config/$BACKUP_DATE/promtail-config.yml deploy/loki/
cp -r backups/config/$BACKUP_DATE/grafana/* deploy/loki/grafana/

# Запустить сервисы
docker-compose -f deploy/loki/docker-compose.yml up -d

echo "✅ Восстановление завершено!"
```

---

## 6. Performance Tuning

### 6.1. Настройка retention

#### Оптимизация retention policies

```yaml
# deploy/loki/loki-config.yml

limits_config:
  # Период хранения логов
  retention_period: 168h  # 7 дней (по умолчанию)
  
  # Максимальный размер строки лога
  max_line_size: 256kb
  
  # Максимальное количество строк на запрос
  max_entries_limit_per_query: 5000

compactor:
  working_directory: /loki/compactor
  shared_store: filesystem
  
  # Включить retention
  retention_enabled: true
  
  # Задержка перед удалением (для безопасности)
  retention_delete_delay: 2h
  
  # Количество workers для удаления
  retention_delete_worker_count: 150
  
  # Интервал compaction
  compaction_interval: 10m
  
  # Интервал очистки
  retention_delete_all_interval: 15m
```

#### Применение retention

```bash
# Перезапустить Loki с новыми настройками
docker-compose -f deploy/loki/docker-compose.yml up -d --force-recreate loki

# Проверить применение
curl http://localhost:3100/loki/api/v1/status | jq .
```

### 6.2. Оптимизация запросов

#### Best practices для LogQL

```logql
# ❌ ПЛОХО: Запрос без фильтра по времени
{app="investor"}

# ✅ ХОРОШО: С фильтром по времени
{app="investor"} | __timestamp__ > 1682000000000000000

# ❌ ПЛОХО: Поиск по всему тексту
{app="investor"} |= "error"

# ✅ ХОРОШО: Использование label filter
{app="investor", level="error"}

# ❌ ПЛОХО: Агрегация за большой период
sum(rate({app="investor"}[1h]))

# ✅ ХОРОШО: Агрегация за короткий период
sum(rate({app="investor"}[5m]))
```

#### Индексация для ускорения

```yaml
# loki-config.yml

schema_config:
  configs:
    - from: 2026-01-01
      store: tsdb
      object_store: filesystem
      schema: v13
      index:
        prefix: index_
        period: 24h

ingester:
  # Размер chunk перед сбросом
  chunk_target_size: 1572864  # 1.5MB
  
  # Максимальное время жизни chunk
  chunk_idle_period: 30m
  
  # Максимальное количество строк в chunk
  chunk_block_size: 262144
```

#### Кэширование запросов

```yaml
# loki-config.yml

query_scheduler:
  max_outstanding_requests_per_tenant: 32768

frontend:
  # Включить кэш запросов
  compress_responses: true
  
  # Максимальное количество запросов в очереди
  max_outstanding_per_tenant: 4096

query_range:
  # Кэширование результатов запросов
  align_queries_with_step: true
  
  # Максимальное количество параллельных запросов
  max_retries: 5
  
  cache_results: true
  
  results_cache:
    cache:
      embedded_cache:
        enabled: true
        max_size_mb: 100
```

### 6.3. Resource Limits

#### Рекомендуемые лимиты для production

```yaml
# deploy/loki/docker-compose.yml

services:
  loki:
    deploy:
      resources:
        limits:
          cpus: '4'
          memory: 4G
        reservations:
          cpus: '2'
          memory: 2G
    
    # Tuning параметров
    command:
      - -config.file=/etc/loki/loki-config.yml
      - -target=all
      - -legacy-read-mode=false
      - -ingester.max-line-size=256kb
      - -ingester.max-chunk-size=1.5MB
      - -querier.max-samples-per-query=5000000

  promtail:
    deploy:
      resources:
        limits:
          cpus: '1'
          memory: 512M
        reservations:
          cpus: '0.5'
          memory: 256M

  grafana:
    deploy:
      resources:
        limits:
          cpus: '2'
          memory: 1G
        reservations:
          cpus: '1'
          memory: 512M
```

#### Monitoring resource usage

```bash
# Проверка использования CPU
docker stats loki promtail grafana --no-stream --format "table {{.Name}}\t{{.CPUPerc}}\t{{.MemUsage}}"

# Проверка через Prometheus
curl http://localhost:3100/metrics | grep -E "go_memstats|process_cpu"

# Алерт на превышение лимитов
# expr: container_memory_usage_bytes{name="loki"} / container_spec_memory_limit_bytes{name="loki"} > 0.9
```

#### Tuning для high-load

```yaml
# loki-config.yml для high-load

limits_config:
  # Увеличенные лимиты
  ingestion_rate_mb: 32
  ingestion_burst_size_mb: 64
  max_streams_per_user: 50000
  max_line_size: 512kb
  max_entries_limit_per_query: 10000
  
  # Cardinality лимиты
  cardinality_limit: 100000

ingester:
  # Увеличенные параметры
  lifecycler:
    ring:
      kvstore:
        store: memberlist
      replication_factor: 1
  
  chunk_idle_period: 1h
  chunk_block_size: 524288
  chunk_target_size: 3145728  # 3MB
  chunk_retain_period: 30s

querier:
  max_concurrent: 20
  timeout: 5m
  query_ingesters_within: 3h

query_scheduler:
  max_outstanding_requests_per_tenant: 65536
```

---

## 📞 Support & Escalation

### Контакты

- **Platform Team:** platform@example.com
- **On-call:** oncall@example.com
- **Slack:** #loki-support

### Эскалация

1. **Level 1:** Проверить runbook, попробовать troubleshooting шаги
2. **Level 2:** Обратиться в Platform Team
3. **Level 3:** On-call engineer
4. **Level 4:** Grafana Labs Support (для production issues)

### Полезные ссылки

- [Loki Documentation](https://grafana.com/docs/loki/)
- [LogQL Reference](https://grafana.com/docs/loki/latest/logql/)
- [Grafana Dashboards](https://grafana.com/grafana/dashboards/)
- [Promtail Configuration](https://grafana.com/docs/loki/latest/clients/promtail/configuration/)

---

## 📝 Changelog

| Version | Date | Changes |
|---------|------|---------|
| 1.0 | 2026-04-20 | Initial release |

---

**Document Owner:** Platform Engineering  
**Review Cycle:** Quarterly  
**Next Review:** 2026-07-20
