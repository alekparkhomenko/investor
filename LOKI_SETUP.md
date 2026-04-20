# Loki Stack - Быстрый старт

## 📦 Что входит

- **Loki 2.9.0** - хранилище логов (порт 3100)
- **Promtail 2.9.0** - агент для сбора логов
- **Grafana 10.2.0** - UI для просмотра (порт 3000, логин/пароль: admin/admin)

---

## 🚀 Запуск

### Шаг 1: Запустите Loki stack

```bash
cd /Users/alekparkhomenko/Projects/investor

# Запуск всех сервисов
docker-compose -f docker-compose.loki.yml up -d

# Проверка статуса
docker-compose -f docker-compose.loki.yml ps
```

Ожидаемый вывод:
```
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
docker-compose -f docker-compose.yml -f docker-compose.loki.yml up -d

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

### Loki (loki/loki-config.yml)

Основные настройки:
- `retention_period: 168h` - хранение 7 дней
- `http_listen_port: 3100` - порт API

### Promtail (promtail/promtail-config.yml)

Собирает логи:
- Docker контейнеры
- Файлы `/var/log/investor/*.log`
- Системные логи

### Grafana (grafana/provisioning/)

Автоматическая настройка:
- Loki datasource
- Default datasource

---

## 🛑 Остановка

```bash
# Остановить все сервисы
docker-compose -f docker-compose.loki.yml down

# Остановить с удалением данных (осторожно!)
docker-compose -f docker-compose.loki.yml down -v
```

---

## 🔧 Troubleshooting

### Loki не запускается

```bash
# Проверить логи
docker-compose -f docker-compose.loki.yml logs loki

# Пересоздать контейнер
docker-compose -f docker-compose.loki.yml up -d --force-recreate loki
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
   docker-compose -f docker-compose.loki.yml logs promtail
   ```

### Grafana не видит Loki

1. Проверьте, что Loki запущен:
   ```bash
   docker-compose -f docker-compose.loki.yml ps loki
   ```

2. Проверьте datasource в Grafana:
   - **Configuration** → **Data sources** → **Loki**
   - Нажмите **Save & test**
   - Должно быть: `Data source is working`

---

## 📚 Дополнительные ресурсы

- [Loki Documentation](https://grafana.com/docs/loki/)
- [LogQL Guide](https://grafana.com/docs/loki/latest/logql/)
- [Grafana Dashboards](https://grafana.com/grafana/dashboards/?dataSource=loki)

---

## 🎯 Next Steps

1. ✅ Запустить Loki stack
2. ✅ Запустить приложение с `LOKI_ENABLED=true`
3. ✅ Открыть Grafana (http://localhost:3000)
4. ✅ Найти логи в Explore
5. 📊 Создать дашборд с метриками
6. 🔔 Настроить алерты на ошибки
