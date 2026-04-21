# Loki Stack - Быстрый старт

## 📦 Что входит

- **Loki 2.9.0** - хранилище логов (порт 3100)
- **Promtail 2.9.0** - агент для сбора логов
- **Grafana 10.2.0** - UI для просмотра (порт 3000, логин/пароль: admin/admin)

---

## 🎯 Taskfile Commands

Проект использует Taskfile для управления Loki stack. Все команды доступны через `task`.

### Основные команды

| Команда | Aliases | Описание |
|---------|---------|----------|
| `task loki:up` | `loki:start` | Запуск Loki stack |
| `task loki:down` | `loki:stop` | Остановка Loki stack |
| `task loki:logs` | `loki:log` | Просмотр логов всех сервисов |
| `task loki:status` | `loki:ps` | Статус сервисов |
| `task loki:validate` | `loki:config`, `loki:check` | Валидация конфигурации |
| `task loki:prune` | `loki:clean` | Очистка (удаление volumes) |

### Примеры использования

**Запуск Loki stack:**
```bash
task loki:up
```

**Просмотр логов в реальном времени:**
```bash
task loki:logs
```

**Проверка статуса сервисов:**
```bash
task loki:status
```

**Валидация конфигурации:**
```bash
task loki:validate
```

**Остановка с очисткой данных:**
```bash
task loki:prune
```

---

## 🚀 Быстрый старт

### Шаг 1: Запустите Loki stack

```bash
cd /Users/alekparkhomenko/Projects/investor

# Запуск всех сервисов через Taskfile
task loki:up

# Проверка статуса
task loki:status
```

Ожидаемый вывод:
```
📊 Статус сервисов Loki stack...
NAME       IMAGE                        STATUS
loki       grafana/loki:2.9.0           Up
promtail   grafana/promtail:2.9.0       Up
grafana    grafana/grafana:10.2.0       Up
```

### Шаг 2: Откройте Grafana

Перейдите в браузере: **http://localhost:3000**

- **Логин:** `admin`
- **Пароль:** `admin`

Grafana уже настроена с datasource Loki!

> 💡 **Совет:** Используйте `task loki:status` для проверки статуса сервисов и `task loki:logs` для просмотра логов в реальном времени.

---

## 📝 Настройка приложения

### Вариант A: Запуск приложения с Loki

В `.env` уже настроено:

```bash
LOKI_ENABLED=true
LOKI_HOST=http://localhost:3100
LOKI_ENV=development
```

Запустите приложение:

```bash
cd investor
go run cmd/main.go
```

Логи будут:
- ✅ Писаться в консоль
- ✅ Отправляться в Loki

### Вариант B: Запуск приложения в Docker

```bash
# Запуск investor вместе с Loki stack
docker-compose -f docker-compose.yml -f deploy/loki/docker-compose.yml up -d

# Просмотр логов
docker-compose logs -f investor
```

---

## 🔍 Поиск логов в Grafana

1. Откройте **Explore** (иконка компаса слева)
2. Выберите datasource **Loki** (сверху)
3. Введите LogQL запрос:

```logql
{app="investor"}
```

Или с фильтрами:

```logql
{app="investor"} |= "ERROR"
{app="investor"} | json | component="moex-ingestor"
{app="investor"} | json | duration_ms > 1000
```

### Примеры запросов

**Все логи приложения:**
```logql
{app="investor"}
```

**Только ошибки:**
```logql
{app="investor"} |= "ERROR"
```

**Медленные запросы (>1с):**
```logql
{app="investor"} | json | duration_ms > 1000
```

**Логи по компоненту:**
```logql
{app="investor"} | json | component="app"
```

**Агрегация по уровням:**
```logql
sum by (level) (count_over_time({app="investor"}[1h]))
```

---

## 📊 Дашборды

### Создание дашборда

1. **Dashboards** → **New Dashboard** → **Add visualization**
2. Выберите **Loki**
3. Введите запрос, например:
   ```logql
   rate({app="investor"} |= "ERROR" [5m])
   ```
4. Выберите визуализацию (Graph, Stat, Table)
5. **Apply** → **Save dashboard**

### Готовые запросы для дашбордов

**Количество логов по уровням:**
```logql
sum by (level) (count_over_time({app="investor"}[5m]))
```

**Количество ошибок:**
```logql
sum(rate({app="investor"} |= "ERROR" [5m]))
```

**Средняя длительность запросов:**
```logql
avg by (component) (rate({app="investor"} | json | duration_ms [5m]))
```

**Топ медленных запросов:**
```logql
topk(5, {app="investor"} | json | duration_ms > 500)
```

---

## ⚙️ Конфигурация

Структура файлов:
```
deploy/loki/
├── docker-compose.yml        # Docker Compose конфигурация
├── loki-config.yml           # Конфигурация Loki
├── promtail-config.yml       # Конфигурация Promtail
└── grafana/
    └── provisioning/
        ├── datasources/      # Datasources (Loki)
        └── dashboards/       # Дашборды
```

### Loki (deploy/loki/loki-config.yml)

Основные настройки:
- `retention_period: 168h` - хранение 7 дней
- `http_listen_port: 3100` - порт API

### Promtail (deploy/loki/promtail-config.yml)

Собирает логи:
- Docker контейнеры
- Файлы `/var/log/investor/*.log`
- Системные логи

### Grafana (deploy/loki/grafana/provisioning/)

Автоматическая настройка:
- Loki datasource
- Default datasource

---

## 🛑 Остановка

```bash
# Остановить все сервисы через Taskfile
task loki:down

# Остановить с удалением данных (осторожно!)
task loki:prune
```

Или используя docker-compose напрямую:

```bash
# Остановить все сервисы
docker-compose -f deploy/loki/docker-compose.yml down

# Остановить с удалением данных (осторожно!)
docker-compose -f deploy/loki/docker-compose.yml down -v
```

---

## 🔧 Troubleshooting

### Loki не запускается

```bash
# Проверить логи через Taskfile
task loki:logs

# Или напрямую
docker-compose -f deploy/loki/docker-compose.yml logs loki

# Пересоздать контейнер
docker-compose -f deploy/loki/docker-compose.yml up -d --force-recreate loki
```

### Логи не попадают в Loki

1. Проверьте, что приложение запущено с `LOKI_ENABLED=true`

2. Проверьте connectivity:
    ```bash
    curl http://localhost:3100/ready
    ```
    Должно вернуть: `ready`

3. Проверьте логи Promtail:
    ```bash
    # Через Taskfile
    task loki:logs
    
    # Или напрямую
    docker-compose -f deploy/loki/docker-compose.yml logs promtail
    ```

4. Проверьте конфигурацию Promtail:
    ```bash
    task loki:validate
    ```

### Grafana не видит Loki

1. Проверьте, что Loki запущен:
    ```bash
    # Через Taskfile
    task loki:status
    
    # Или напрямую
    docker-compose -f deploy/loki/docker-compose.yml ps loki
    ```

2. Проверьте datasource в Grafana:
    - **Configuration** → **Data sources** → **Loki**
    - Нажмите **Save & test**
    - Должно быть: `Data source is working`

### Распространённые проблемы

| Проблема | Решение |
|----------|---------|
| Порт 3100 занят | Проверьте: `lsof -i :3100`, остановите конфликтующий сервис |
| Порт 3000 занят | Проверьте: `lsof -i :3000`, измените порт в docker-compose.yml |
| Loki не принимает логи | Проверьте `LOKI_HOST` в .env, должен быть `http://localhost:3100` |
| Promtail не отправляет логи | Проверьте логи: `task loki:logs`, убедитесь что Loki доступен |
| Grafana не показывает логи | Проверьте datasource, выберите правильный time range в Explore |

### Команды для диагностики

```bash
# Проверка статуса всех сервисов
task loki:status

# Проверка доступности Loki API
curl http://localhost:3100/ready
curl http://localhost:3100/loki/api/v1/status

# Просмотр логов конкретного сервиса
docker-compose -f deploy/loki/docker-compose.yml logs loki
docker-compose -f deploy/loki/docker-compose.yml logs promtail
docker-compose -f deploy/loki/docker-compose.yml logs grafana

# Валидация конфигурации
task loki:validate

# Перезапуск конкретного сервиса
docker-compose -f deploy/loki/docker-compose.yml restart loki
```

---

## 📚 Дополнительные ресурсы

- [Loki Documentation](https://grafana.com/docs/loki/)
- [LogQL Guide](https://grafana.com/docs/loki/latest/logql/)
- [Grafana Dashboards](https://grafana.com/grafana/dashboards/?dataSource=loki)
- [Taskfile Documentation](https://taskfile.dev/)

---

## 🎯 Next Steps

1. ✅ Запустить Loki stack: `task loki:up`
2. ✅ Запустить приложение с `LOKI_ENABLED=true`
3. ✅ Открыть Grafana (http://localhost:3000)
4. ✅ Найти логи в Explore
5. 📊 Создать дашборд с метриками
6. 🔔 Настроить алерты на ошибки
