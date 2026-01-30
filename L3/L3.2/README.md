# Shortener — сервис сокращения URL с аналитикой

Mini-сервис для создания коротких ссылок с возможностью отслеживания переходов: кто перешёл, когда и с какого устройства.

---

## Возможности

- Создание коротких ссылок (`POST /shorten`)
- Переход по короткой ссылке (`GET /s/{short_code}`)
- Просмотр аналитики переходов (`GET /analytics/{short_code}`)
- Поддержка кастомных коротких кодов
- Кэширование популярных ссылок через Redis

---

## Структура проекта

```
.
├── cmd/server/main.go        # Точка входа
├── internal/
│   ├── cache/               # Redis
│   ├── config/              # Конфиги
│   ├── db/                  # Postgres
│   ├── handlers/            # HTTP-обработчики
│   └── repository/          # Работа с БД и логика ссылок
├── migrations/              # SQL миграции
├── frontend/                # HTML+JS фронтенд
├── Dockerfile
├── docker-compose.yml
├── config.yaml
└── .env
```

---

## Настройка

1. Скопируй `.env.example` в `.env` и укажи свои данные:

```env
# PostgreSQL
POSTGRES_HOST=postgres
POSTGRES_PORT=5432
POSTGRES_USER=your_user
POSTGRES_PASSWORD=your_password
POSTGRES_DB=shortener

# Redis
REDIS_ADDR=redis:6379
REDIS_PASSWORD=
REDIS_TTL=3600

# Server
SERVER_HOST=0.0.0.0
SERVER_PORT=8080
```

2. Настрой `config.yaml` при необходимости:

```yaml
server:
  host: "0.0.0.0"
  port: 8080

redis:
  enabled: true
  addr: "redis:6379"
  password: ""
  db: 0
  ttl: "3600s"

postgres:
  max_open_conns: 10
  max_idle_conns: 5
  conn_max_lifetime: "5m"
```

---

## Запуск через Docker

```bash
docker-compose up --build
```

- Приложение будет доступно на `http://localhost:8080`
- Postgres и Redis запускаются в контейнерах

---

## REST API

### 1. Создать короткую ссылку

```bash
curl -X POST http://localhost:8080/shorten \
-H "Content-Type: application/json" \
-d '{"url":"https://example.com"}'
```

Пример с кастомным кодом:

```bash
curl -X POST http://localhost:8080/shorten \
-H "Content-Type: application/json" \
-d '{"url":"https://example.com","custom_code":"mycode"}'
```

Ответ:

```json
{
  "code": "abc123"
}
```

---

### 2. Переход по короткой ссылке

```bash
curl -v http://localhost:8080/s/abc123
```

- Перенаправляет на оригинальный URL

---

### 3. Получить аналитику

```bash
curl -v http://localhost:8080/analytics/abc123
```

Ответ:

```json
[
  {
    "date": "2026-01-30T00:00:00Z",
    "user_agent": "Mozilla/5.0 ...",
    "count": 5
  }
]
```

> ⚠️ Важно: указывать именно `short_code`, а не полный URL.

---

## Фронтенд

- Расположен в папке `frontend/`
- Простая HTML + JS страница для:
  - Создания ссылок
  - Просмотра аналитики

> При запуске фронтенда локально убедись, что в `main.go` настроен CORS:

```go
AllowOrigins: []string{"http://localhost:5500"} // адрес фронтенда
```

---

## Дополнительно

- Все ссылки и переходы сохраняются в Postgres
- Популярные ссылки кэшируются в Redis
- Graceful shutdown реализован — безопасное завершение работы сервера

---

## Миграции SQL

Пример миграции для таблиц `short_links` и `visits`:

```sql
-- short_links
CREATE TABLE IF NOT EXISTS short_links (
    id BIGSERIAL PRIMARY KEY,
    short_code VARCHAR(32) NOT NULL UNIQUE,
    original_url TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    expires_at TIMESTAMPTZ,
    is_active BOOLEAN NOT NULL DEFAULT TRUE
);

-- visits
CREATE TABLE IF NOT EXISTS visits (
    id BIGSERIAL PRIMARY KEY,
    short_link_id BIGINT NOT NULL REFERENCES short_links(id) ON DELETE CASCADE,
    user_agent TEXT,
    ip_address TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
```

---

## Примечания

- Для локальной разработки фронтенд можно запускать через простой сервер Python:

```bash
cd frontend
python -m http.server 5500
```

- Для успешной работы API убедись, что фронтенд и сервер находятся в CORS-совместимых сетях.

