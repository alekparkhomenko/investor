**План реализации MVP: Data Ingestor (MOEX ISS) на Go**

## Этап 0. Подготовка проекта

1. Инициализировать Go-модуль:

   ```bash
   go mod init mvp-ingestor
   ```
2. Создать базовую структуру:

   ```
   /cmd/app/main.go
   /internal/ingestor/moex.go
   /internal/model/quote.go
   ```

---

## Этап 1. Описание модели данных

1. Создать структуру `Quote`:

   ```go
   type Quote struct {
       Symbol string
       Price  float64
       Time   int64
   }
   ```

2. Создать структуру для ответа ISS API Московская биржа:

   ```go
   type ISSResponse struct {
       MarketData struct {
           Columns []string        `json:"columns"`
           Data    [][]interface{} `json:"data"`
       } `json:"marketdata"`
   }
   ```

---

## Этап 2. Реализация HTTP-клиента

1. Создать HTTP-клиент:

   * timeout: 5 секунд
   * переиспользуемый (`http.Client`)

2. Реализовать функцию:

   ```go
   func fetchISS(ctx context.Context, client *http.Client, symbol string) (*ISSResponse, error)
   ```

3. Внутри:

   * сформировать URL:

     ```
     https://iss.moex.com/iss/engines/stock/markets/shares/securities/{SYMBOL}.json
     ```
   * выполнить GET-запрос
   * декодировать JSON
   * закрыть `resp.Body`

---

## Этап 3. Парсинг данных ISS

1. Реализовать функцию:

   ```go
   func parseQuote(symbol string, resp *ISSResponse) (Quote, error)
   ```

2. Логика:

   * проверить, что `data` не пустой
   * найти индекс колонки `LAST` в `columns`
   * извлечь значение цены по индексу
   * привести к `float64`

3. Вернуть `Quote`:

   * Symbol → symbol
   * Price → LAST
   * Time → `time.Now().Unix()`

---

## Этап 4. Реализация polling-ингестора

1. Реализовать функцию:

   ```go
   func StartMOEXIngestor(ctx context.Context, symbol string, out chan<- Quote)
   ```

2. Внутри:

   * создать `time.Ticker` (интервал 1–2 секунды)
   * в цикле:

     * вызвать `fetchISS`
     * вызвать `parseQuote`
     * отправить результат в канал

3. Обработать:

   * ошибки HTTP
   * ошибки парсинга
   * `ctx.Done()` (graceful shutdown)

---

## Этап 5. Конкурентная обработка

1. Запуск ingestion в goroutine:

   ```go
   go StartMOEXIngestor(ctx, "SBER", quotesChan)
   ```

2. Использовать буферизированный канал:

   ```go
   quotesChan := make(chan Quote, 100)
   ```

3. Реализовать consumer (временно для логирования):

   ```go
   for q := range quotesChan {
       log.Printf("%s = %.2f", q.Symbol, q.Price)
   }
   ```

---

## Этап 6. Graceful Shutdown

1. Использовать `context.WithCancel`
2. Обработать сигналы ОС (SIGINT, SIGTERM)
3. При остановке:

   * остановить ticker
   * завершить goroutine
   * закрыть канал (при необходимости)

---

## Этап 7. Логирование

1. Логировать:

   * ошибки HTTP
   * ошибки парсинга
   * успешные котировки

2. Использовать стандартный `log` пакет

---

## Этап 8. Конфигурация

1. Добавить переменные окружения:

   * `SYMBOL=SBER`
   * `POLL_INTERVAL=2s`

2. Считать их в `main.go`

---

## Этап 9. Тестирование

1. Написать unit-тест:

   * на `parseQuote`
   * использовать mock JSON от ISS

2. Проверить:

   * корректное извлечение цены
   * обработку пустых данных

---

## Этап 10. Проверка результата

1. Запустить сервис:

   ```bash
   go run ./cmd/app
   ```

2. Ожидаемый результат:

   * каждые 1–2 секунды логируется цена акции (например SBER)
   * сервис не падает при ошибках сети
   * корректно завершается при остановке

---

## Результат этапа

Готовый Data Ingestor:

* получает данные с ISS API
* обрабатывает нестандартный JSON (columns/data)
* работает асинхронно
* устойчив к ошибкам
* готов к интеграции с Alert Engine

---

## Следующий этап

Реализация Alert Engine (обработка условий и генерация алертов)
