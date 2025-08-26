# Project-go-url-shortener

[![Go](https://img.shields.io/badge/Go-1.20-blue)](https://golang.org/)  
[![License](https://img.shields.io/badge/license-MIT-green)](LICENSE)

**URL Shortener** — это простой и быстрый сервис для сокращения URL, реализованный на Go с использованием Gin, Zap, Swagger и Prometheus.  

Проект предназначен для демонстрации навыков backend-разработки.

---

## Функционал

- Создание коротких URL по длинным ссылкам
- Перенаправление с короткого URL на исходный
- Rate limiting для защиты API
- Healthcheck эндпоинт
- Swagger документация для API
- Prometheus метрики для мониторинга
- Логирование через Zap

---

## Структура проекта
```bash
Project-go-url-shortener/
├─ cmd/shortener/ # Точка входа приложения
│ └─ main.go
├─ handler/ # Бизнес-хендлеры
│ ├─ url.go # CreateShortUrl и HandleShortUrlRedirect
│ └─ health.go # Healthcheck
├─ internal/
│ └─ http/middleware/ # RateLimit, Logging и другие middleware
├─ store/ # Хранилище URL-сокращений
├─ docs/ # Swagger документация
├─ go.mod
├─ go.sum
└─ README.md
```

---

## Установка и запуск

### 1. Клонирование репозитория
```bash
git clone https://github.com/DmitryApareev/Project-go-url-shortener.git
cd Project-go-url-shortener
`````

### 2. Установка зависимостей
```bash
go mod tidy
```

### 3. Генерация Swagger документации
```bash
go install github.com/swaggo/swag/cmd/swag@latest
swag init -g ./cmd/shortener/main.go --output ./docs
```

### 4. Запуск сервера
```bash
go run ./cmd/shortener
```
По умолчанию сервер запускается на порту 8080

### Пример запроса на создание короткого URL
```http
POST /api/shorten
Content-Type: application/json

{
  "long_url": "https://example.com/some/long/path",
  "user_id": "123"
}
```

### Пример ответа:
```http
{
  "message": "short url created successfully",
  "short_url": "http://localhost:9808/abc123"
}
```

### Запустить все unit-тесты:
```bach
go test ./... -v
```

---

## Зависимости

### [Gin](https://github.com/gin-gonic/gin) — HTTP-фреймворк

### [Zap](https://github.com/uber-go/zap) — логирование

### [Prometheus client](https://github.com/prometheus/client_golang) — метрики

### [Swaggo](https://github.com/swaggo/swag) — генерация Swagger документации