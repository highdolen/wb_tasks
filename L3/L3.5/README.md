# EventBooker

**EventBooker** — сервис бронирования мест на мероприятия с автоматической отменой неоплаченных броней через заданный интервал времени (TTL).  

Проект моделирует реальную систему регистрации на события, где важно освобождать места, если пользователь не подтвердил бронь вовремя.

Сервис включает:

- REST API
- Фоновый обработчик истечения бронирований
- Веб-интерфейс для пользователей и администратора
- Работу с PostgreSQL
- Docker окружение
- Уведомления через Telegram

---

## Основные возможности

### Работа с мероприятиями

- Создание мероприятий
- Просмотр списка событий
- Информация о свободных местах
- Настройка TTL бронирования для каждого события

### Бронирование

- Пользователь может забронировать место
- Бронь имеет статус `pending`
- Пользователь может подтвердить бронь (оплатить)
- Статус меняется на `confirmed`

### Автоматическая отмена бронирований

Если бронь не подтверждена в течение TTL:

- Бронь автоматически отменяется
- Статус меняется на `canceled`
- Место освобождается
- Пользователь получает уведомление (Telegram)

### Поддержка нескольких пользователей

Система поддерживает:

- Регистрацию пользователей
- Роли (`user` / `admin`)
- Уникальные бронирования на событие

### Web интерфейс

Есть две страницы:

**Администратор:**

- Создание событий
- Просмотр событий

**Пользователь:**

- Просмотр событий
- Бронирование
- Подтверждение брони
- Отслеживание статуса

---

## Архитектура проекта

Проект построен по принципам **Clean Architecture**:


handlers → services → repositories → database



Слои:

- `handlers` — HTTP обработчики
- `services` — бизнес логика
- `repositories` — работа с базой
- `scheduler` — фоновые процессы
- `notification` — интеграции (Telegram)

---

## Структура проекта

```
eventBooker/
│
├── cmd/
│   └── server/
│       └── main.go              # точка входа приложения, запуск сервера
│
├── internal/                    # внутренняя логика приложения
│   ├── booking/                 # логика бронирований (model, repository, service)
│   ├── event/                   # логика мероприятий
│   ├── user/                    # работа с пользователями
│   ├── handlers/                # HTTP обработчики (API)
│   ├── scheduler/               # фоновый обработчик отмены просроченных броней
│   ├── notification/            # уведомления (Telegram)
│   ├── storage/                 # подключение к PostgreSQL
│   └── config/                  # загрузка и валидация конфигурации
│
├── frontend/                    # простой веб-интерфейс
│   ├── login.html
│   ├── user.html
│   ├── admin.html
│   └── static/
│       ├── css/
│       │   └── style.css
│       └── js/
│           ├── login.js
│           ├── user.js
│           └── admin.js
│
├── migrations/                  # SQL миграции базы данных
│   ├── 0001_init.up.sql
│   └── 0001_init.down.sql
│
├── config.yaml                  # основной конфигурационный файл
├── .env                         # переменные окружения
├── Dockerfile                   # сборка приложения
├── docker-compose.yml           # запуск приложения и PostgreSQL
├── go.mod
└── go.sum
```

---

---

## Как работает TTL бронирования

Каждое событие имеет параметр `booking_ttl` — время жизни брони.  

Когда пользователь бронирует место:

1. Создаётся запись в таблице `bookings`
2. Статус устанавливается `pending`
3. Устанавливается `expires_at`

Если пользователь не подтверждает бронь:

Фоновый scheduler:

- Проверяет просроченные брони
- Отменяет их
- Возвращает место в событие

Scheduler запускается каждые 10 секунд.

---

## База данных

Используется PostgreSQL.  

### Таблица `users`

| Поле        | Тип     |
|------------|--------|
| id         | BIGSERIAL PRIMARY KEY |
| email      | TEXT UNIQUE NOT NULL |
| name       | TEXT NOT NULL |
| telegram_id| TEXT |
| role       | VARCHAR(20) ('user' / 'admin') |
| created_at | TIMESTAMP NOT NULL DEFAULT NOW() |

### Таблица `events`

| Поле           | Тип     |
|----------------|--------|
| id             | BIGSERIAL PRIMARY KEY |
| name           | TEXT NOT NULL |
| date           | TIMESTAMP NOT NULL |
| total_seats    | INT NOT NULL |
| available_seats| INT NOT NULL |
| booking_ttl    | INT NOT NULL |
| created_at     | TIMESTAMP NOT NULL DEFAULT NOW() |

### Таблица `bookings`

| Поле       | Тип     |
|------------|--------|
| id         | BIGSERIAL PRIMARY KEY |
| event_id   | BIGINT REFERENCES events(id) |
| user_id    | BIGINT REFERENCES users(id) |
| status     | TEXT ('pending', 'confirmed', 'canceled') |
| created_at | TIMESTAMP NOT NULL DEFAULT NOW() |
| expires_at | TIMESTAMP NOT NULL |

---

## Индексы

Для производительности используются индексы:

- Поиск просроченных броней
- Поиск броней по событию
- Поиск броней пользователя
- Защита от дублирующих броней

---

## API

Базовый URL сервера:

```
http://localhost:8080
```

Перед тестированием убедитесь, что сервис запущен:

```bash
docker compose up --build
```

---

### Создать мероприятие

**POST /events**

```bash
curl -X POST http://localhost:8080/events \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Go Conference",
    "date": "2026-03-01T18:00",
    "total_seats": 50,
    "booking_ttl_minutes": 10
  }'
```

---

### Получить список мероприятий

**GET /events**

```bash
curl http://localhost:8080/events
```

---

### Получить одно мероприятие

**GET /events/{id}**

Пример:

```bash
curl http://localhost:8080/events/1
```

---

### Забронировать место на мероприятии

**POST /events/{id}/book**

```bash
curl -X POST http://localhost:8080/events/1/book \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@mail.com",
    "name": "User",
    "telegram_id": "123456",
    "role": "user"
  }'
```

---

### Подтвердить бронь (оплата)

**POST /bookings/{id}/confirm**

```bash
curl -X POST http://localhost:8080/bookings/1/confirm
```

---

### Получить брони пользователя

**GET /users/{id}/bookings**

```bash
curl http://localhost:8080/users/1/bookings
```

## Web интерфейс

После запуска сервиса доступны три страницы:

| Страница | URL | Описание |
|----------|-----|----------|
| Главная | `/` | страница входа |
| Пользователь | `/user` | просмотр событий, бронирование и подтверждение |
| Администратор | `/admin` | создание и управление событиями |

---

## Запуск проекта через Docker

```bash
docker compose up --build
```

После запуска сервис будет доступен по адресу:

```
http://localhost:8080
```

---

## Что происходит при запуске

Docker поднимает следующие сервисы:

1. **PostgreSQL** — база данных  
2. **Migrations** — применяются SQL миграции  
3. **Application** — запускается сервер EventBooker

---

## Конфигурация

Основной конфигурационный файл:

```
config.yaml
```

Пример конфигурации:

```yaml
server:
  host: 0.0.0.0
  port: 8080

scheduler:
  interval: 10s
```

---

## Переменные окружения

Файл:

```
.env
```

Используется для настройки:

- подключения к PostgreSQL
- Telegram уведомлений

---

## Уведомления

Сервис поддерживает уведомления через Telegram:

- подтверждение брони
- отмена просроченной брони

Если Telegram токен не указан — уведомления просто отключаются.

---

## Безопасность данных

В проекте используются:

- транзакции базы данных
- ограничения (constraints)
- индексы
- уникальные ключи

Это предотвращает:

- гонки данных
- овербукинг мест
- дублирование бронирований

---

## Технологии

- Go
- PostgreSQL
- Docker
- Gin (ginext)
- JavaScript
- HTML / CSS