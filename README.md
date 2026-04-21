# Investor — MOEX Stock Quote Ingestor

Приложение для мониторинга котировок акций MOEX с системой уведомлений в Telegram.

## 🚀 Быстрый старт

### Требования

- [Docker](https://docker.com)
- [Task](https://taskfile.dev) (опционально)

### Запуск всех сервисов

```bash
# Через Task (рекомендуется)
task docker:up

# Или напрямую через Docker
docker compose -f deploy/app/docker-compose.yml up -d --build
docker compose -f deploy/loki/docker-compose.yml up -d
```

### Остановка

```bash
task docker:down
```

## 🐳 Сервисы

| Сервис | URL | Описание |
|--------|-----|----------|
| **App** | http://localhost:8080 | Основное приложение |
| **Swagger UI** | http://localhost:8080/swagger/index.html | API документация |
| **PostgreSQL** | localhost:5432 | База данных |
| **Grafana** | http://localhost:3000 | Логирование / Дашборды |
| **Loki** | http://localhost:3100 | Логи |

### Учетные данные

| Сервис | Логин | Пароль |
|--------|-------|--------|
| Grafana | `admin` | `admin` |
| PostgreSQL | `investor` | `investor` |

## 📡 API Endpoints

| метод | endpoint | описание |
|-------|---------|----------|
| GET | `/api/v1/tickers` | Список доступных тикеров |
| GET | `/api/v1/portfolio` | Портфель пользователя |
| POST | `/api/v1/portfolio` | Добавить тикеры в портфель |
| DELETE | `/api/v1/portfolio/{ticker}` | Удалить тикер |

## 🔧 Конфигурация

Переменные окружения в `.env`:

```bash
# База данных
DATABASE_URL=postgres://investor:investor@postgres:5432/investor?sslmode=disable

# MOEX тикеры
SYMBOLS=SBER,GAZP,MOEX

# Интервал опроса
POLL_INTERVAL=4s

# Telegram
TELEGRAM_TOKEN=<ваш_токен>
```

## 🛠 Разработка

### Локальный запуск

```bash
# Установить зависимости
go work sync
go mod tidy -compat=1.24

# Запустить
go run ./investor/cmd/main.go
```

### Тесты

```bash
go test ./...
```

### Линтинг

```bash
task lint
```

### Миграции БД

```bash
# Применить миграции
task db:migrate:apply

# Docker контейнер
docker exec -i app-postgres-1 psql -U investor -d investor < investor/migrations/001_portfolio.up.sql
```

## 📂 Структура проекта

```
investor/
├── cmd/              # Точка входа
├── internal/         # Внутренние пакеты
│   ├── app/        # Бизнес-логика
│   ├── config/     # Конфигурация
│   ├── http/      # HTTP сервер
│   ├── ingestor/  # MOEX ingester
│   ├── model/     # Модели
│   ├── storage/   # PostgreSQL
│   └── metrics/  # Метрики
└── migrations/   # SQL миграции

plantform/             # Shared пакеты
└── pkg/
    ├── closer/    # Graceful shutdown
    └── logger/   # Structured logging

deploy/
├── app/           # docker-compose для app
├── postgres/     # PostgreSQL
└── loki/        # Loki + Grafana
```

## 🏗 Tech Stack

- **Go 1.25** — Язык
- **PostgreSQL** — База данных
- **Gorilla/Mux** — HTTP роутер
- **pgx** — PostgreSQL драйвер
- **zap** — Structured logging
- **Docker** — Контейнеризация
- **Loki + Grafana** — Логирование