# Product Vision: Investor

## Problem Statement

Трейдеры и инвесторы на российском рынке (MOEX) тратят время на ручной мониторинг котировок акций. Нет удобного инструмента для получения real-time данных и настройки алертов.

## Solution

**Investor** — CLI приложение для автоматического получения котировок с биржи MOEX и отправки алертов при значимых изменениях цен.

## Target Audience

- Приватные трейдеры и инвесторы
- Активные участники российского фондового рынка (MOEX)
- Пользователи, которым нужен простой CLI для мониторинга акций

## Unique Value Proposition

> "Получайте данные с MOEX в реальном времени и настраивайте алерты — без сложных терминалов и дорогих подписок."

## Key Benefits

- Real-time котировки с MOEX (SBER, GAZP, MOEX и др.)
- Настраиваемые алерты (Telegram)
- Простой запуск и настройка через .env
- Low memory footprint

## Success Metrics

- Приложение запускается и получает данные за <5 секунд
- 100% uptime при работе 24/7
- Алерты доставляются <1 секунды после срабатывания

## Technical Stack

- Go 1.25+
- MOEX ISS API (https://iss.moex.com)
- Telegram Bot API (опционально)
- No external DB (in-memory)

## Current Status

- MVP работает: получает котировки SBER, GAZP, MOEX каждые 4 секунды
- Логирование через fmt.Println (debug mode)

---
**Version:** 0.1.0 | **Date:** 2026-04-19
