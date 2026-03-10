# WarehouseControl

**WarehouseControl** — сервис для учёта товаров на складе с историей изменений и ролевой моделью доступа.

Любой склад — это не просто список товаров. Это движущаяся система, где кто-то добавляет, кто-то списывает, кто-то редактирует остатки, и важно понимать, кто, что, когда и зачем сделал.

Сервис реализует:

- CRUD операции с товарами
- Логирование всех изменений с указанием пользователя и времени
- Ролевую модель доступа (`admin`, `manager`, `viewer`)
- Экспорт истории в CSV
- Веб-интерфейс для пользователей с разными ролями
- Авторизацию через JWT

> История изменений реализована через **триггеры в PostgreSQL** — это антипаттерн, который специально используется для образовательных целей.

---

## Основные возможности

### CRUD товаров

- **Создать товар** — `POST /items`  
- **Получить список товаров** — `GET /items`  
- **Обновить товар** — `PUT /items/{id}`  
- **Удалить товар** — `DELETE /items/{id}`  

### История изменений

- Сохраняется автоматически через триггер в PostgreSQL
- Доступна через API `GET /items/{id}/history`
- Содержит: старое и новое значение, пользователя, время изменения, тип действия (`INSERT`, `UPDATE`, `DELETE`)
- Возможность просмотра отличий между версиями (diff)

### Роли и авторизация

- `admin` — полный доступ, CRUD, просмотр истории
- `manager` — просмотр и редактирование товаров
- `viewer` — только просмотр товаров
- JWT используется для передачи роли и проверки доступа на каждом запросе

### Экспорт истории

- **GET /audit/export** — экспорт истории изменений в CSV

### Web интерфейс

- Вход пользователя с выбором роли  
- Просмотр списка товаров  
- Создание, редактирование, удаление (по правам)  
- Просмотр истории изменений  

---

## Структура проекта

```
warehouseControl/
│   .env
│   config.yaml
│   docker-compose.yml
│   Dockerfile
│   go.mod
│   go.sum
│
├───cmd
│   └───server
│           main.go
│
├───frontend
│   │   index.html
│   │
│   ├───pages
│   │       login.html
│   │       items.html
│   │       audit.html
│   │
│   └───static
│       ├───components
│       │       table.js
│       ├───css
│       │       style.css
│       └───js
│               api.js
│               app.js
│               auth.js
│               items.js
│               audit.js
│               history.js
│
├───internal
│   ├───audit
│   │       model.go
│   │       repository.go
│   │       service.go
│   ├───auth
│   │       jwt.go
│   │       middleware.go
│   │       roles.go
│   │       service.go
│   ├───config
│   │       config.go
│   │       loadConfig.go
│   ├───export
│   │       repository.go
│   │       service.go
│   ├───handlers
│   │       auth.go
│   │       items.go
│   │       audit.go
│   │       export.go
│   │       routes.go
│   ├───items
│   │       model.go
│   │       repository.go
│   │       service.go
│   ├───storage
│   │       postgres.go
│   └───users
│           model.go
│           repository.go
│           service.go
└───migrations
        0001_init.up.sql
        0001_init.down.sql
```

---

## База данных

### Таблица `users`

| Поле      | Тип  | Описание             |
|-----------|-----|--------------------|
| id        | SERIAL PRIMARY KEY | идентификатор пользователя |
| username  | TEXT UNIQUE NOT NULL | имя пользователя |
| role      | TEXT NOT NULL | роль пользователя (`admin`, `manager`, `viewer`) |

### Таблица `items`

| Поле       | Тип | Описание           |
|------------|-----|------------------|
| id         | SERIAL PRIMARY KEY | идентификатор товара |
| name       | TEXT NOT NULL | название товара |
| quantity   | INT NOT NULL | количество на складе |
| updated_at | TIMESTAMP DEFAULT now() | время последнего обновления |

### Таблица `item_history`

| Поле       | Тип      | Описание                     |
|------------|----------|-----------------------------|
| id         | SERIAL PRIMARY KEY | идентификатор записи истории |
| item_id    | INT      | id товара                     |
| action     | TEXT     | действие (`INSERT`/`UPDATE`/`DELETE`) |
| old_value  | JSONB    | старые значения              |
| new_value  | JSONB    | новые значения              |
| changed_by | TEXT     | кто сделал изменение         |
| changed_at | TIMESTAMP DEFAULT now() | время изменения |

> История реализована через триггер `items_audit_trigger` и функцию `audit_item_changes()`.

---

## Примеры API

### Авторизация (JWT)

**POST /login**

```bash
curl -X POST http://localhost:8080/login \
-H "Content-Type: application/json" \
-d '{
  "username": "alex",
  "role": "admin"
}'
```

Ответ содержит токен JWT:

```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVC..."
}
```

---

### Создание товара

**POST /items**

```bash
curl -X POST http://localhost:8080/items \
-H "Authorization: Bearer $TOKEN" \
-H "Content-Type: application/json" \
-d '{
  "name": "Laptop",
  "quantity": 10
}'
```

---

### Получение списка товаров

**GET /items**

```bash
curl -X GET http://localhost:8080/items \
-H "Authorization: Bearer $TOKEN"
```

---

### Обновление товара

**PUT /items/{id}**

```bash
curl -X PUT http://localhost:8080/items/1 \
-H "Authorization: Bearer $TOKEN" \
-H "Content-Type: application/json" \
-d '{
  "name": "Laptop",
  "quantity": 20
}'
```

---

### Удаление товара

**DELETE /items/{id}**

```bash
curl -X DELETE http://localhost:8080/items/2 \
-H "Authorization: Bearer $TOKEN"
```

---

### Получение истории товара

**GET /items/{id}/history**

```bash
curl -X GET http://localhost:8080/items/1/history \
-H "Authorization: Bearer $TOKEN"
```

---

### Экспорт истории в CSV

**GET /audit/export**

```bash
curl http://localhost:8080/audit/export \
-H "Authorization: Bearer $TOKEN" \
-o history.csv
```

---

## Запуск проекта через Docker

```bash
docker compose up --build
```

Сервис будет доступен по адресу:

```
http://localhost:8080
```

Сервисы:

1. **PostgreSQL** — база данных  
2. **Migrate** — применяет миграции  
3. **App** — запускает сервер WarehouseControl  

---

## Конфигурация

Файл `config.yaml` и `.env` для настроек:

- Подключение к PostgreSQL  
- Порты сервера  
- JWT секрет  

---

## Технологии

- Go  
- PostgreSQL  
- Docker / docker-compose  
- Gin (`ginext`)  
- JavaScript, HTML, CSS

