Отлично — вот это уже уровень «я думаю как инженер, а не просто пишу код». Давай соберём **понятную и реалистичную схему под Docker Compose**, которую можно реально запустить локально и показать на собесе.

---

# 🧩 Общая архитектура (Compose-ready)

У тебя будет **5 контейнеров**:

1. `ingestor` — читает WebSocket, пишет в Kafka
2. `alert-engine` — читает из Kafka, проверяет условия
3. `notifier` — отправляет в Telegram
4. `kafka` — брокер сообщений
5. `redis` — быстрый стор алертов

(+ опционально `postgres` и `zookeeper`, если не используешь KRaft)

---

# 🔄 Схема взаимодействия

```
          WebSocket API (Finnhub)
                    │
                    ▼
              [ ingestor ]
                    │
                    ▼
         (topic: raw_quotes)
              Kafka
                    │
                    ▼
           [ alert-engine ]
        │                │
        ▼                ▼
     Redis        (topic: alert_triggered)
                          │
                          ▼
                    [ notifier ]
                          │
                          ▼
                     Telegram API
```

---

# 🐳 Docker Compose (базовый пример)

Вот минимальный, но «боевой» вариант:

```yaml
version: '3.9'

services:
  kafka:
    image: bitnami/kafka:latest
    container_name: kafka
    ports:
      - "9092:9092"
    environment:
      - KAFKA_CFG_NODE_ID=0
      - KAFKA_CFG_PROCESS_ROLES=broker,controller
      - KAFKA_CFG_CONTROLLER_QUORUM_VOTERS=0@kafka:9093
      - KAFKA_CFG_LISTENERS=PLAINTEXT://:9092,CONTROLLER://:9093
      - KAFKA_CFG_ADVERTISED_LISTENERS=PLAINTEXT://kafka:9092
      - KAFKA_CFG_CONTROLLER_LISTENER_NAMES=CONTROLLER
      - ALLOW_PLAINTEXT_LISTENER=yes

  redis:
    image: redis:7
    container_name: redis
    ports:
      - "6379:6379"

  ingestor:
    build: ./ingestor
    depends_on:
      - kafka
    environment:
      - KAFKA_BROKER=kafka:9092

  alert-engine:
    build: ./alert-engine
    depends_on:
      - kafka
      - redis
    environment:
      - KAFKA_BROKER=kafka:9092
      - REDIS_ADDR=redis:6379

  notifier:
    build: ./notifier
    depends_on:
      - kafka
    environment:
      - KAFKA_BROKER=kafka:9092
      - TELEGRAM_TOKEN=your_token_here
```

---

# 📦 Структура проекта

```
project/
│
├── docker-compose.yml
│
├── ingestor/
│   ├── main.go
│   └── Dockerfile
│
├── alert-engine/
│   ├── main.go
│   └── Dockerfile
│
├── notifier/
│   ├── main.go
│   └── Dockerfile
│
└── shared/
    ├── models.go
    └── kafka.go
```

---

# 🔥 Важные детали (которые выделят тебя)

## 1. Healthchecks (очень круто смотрится)

Добавь, например:

```yaml
healthcheck:
  test: ["CMD", "nc", "-z", "localhost", "9092"]
  interval: 10s
  timeout: 5s
  retries: 5
```

---

## 2. Wait-for Kafka (иначе всё падает)

Compose не гарантирует готовность Kafka.

Решение:

* retry loop в Go
* или `wait-for-it.sh`

---

## 3. Topic init (часто забывают)

Можно добавить init-контейнер:

```yaml
kafka-init:
  image: bitnami/kafka:latest
  depends_on:
    - kafka
  entrypoint: >
    bash -c "
    kafka-topics.sh --create --topic raw_quotes --bootstrap-server kafka:9092 --partitions 3 --replication-factor 1;
    kafka-topics.sh --create --topic alert_triggered --bootstrap-server kafka:9092 --partitions 3 --replication-factor 1;
    "
```

---

## 4. Логирование (must-have)

Просто, но важно:

```go
log.Printf("Received quote: %s = %f", symbol, price)
```

---

## 5. Конфигурация через env

Ты уже это сделал — это плюс на интервью.

---

# 🧠 Как это объяснить на собесе

Коротко и мощно:

> "Я использую Apache Kafka как центральную шину событий.
> Ingestor пишет сырые котировки, Alert Engine обрабатывает их и генерирует события алертов.
> Redis используется как in-memory store для быстрого доступа к подпискам.
> Все сервисы изолированы и запускаются через Docker Compose, что позволяет легко масштабировать Alert Engine через consumer groups."

---

# 🚀 Хочешь следующий шаг?

Можем дальше углубиться:

* написать **dockerfile для каждого сервиса**
* реализовать **ingestor (WebSocket + Kafka producer)**
* или **alert-engine с Redis и worker pool**

Скажи, куда копать дальше 👇
