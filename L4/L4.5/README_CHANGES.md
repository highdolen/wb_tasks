# README_CHANGES

## Что уже сделано

### 1. Инфраструктура и запуск
- Поднят локальный стенд через `docker-compose`:
  - приложение
  - PostgreSQL
  - Kafka
  - Kafka UI
- Исправлено применение миграции при старте PostgreSQL
- Добавлен проброс порта `6060` для `pprof`

### 2. Профилирование
- Подключен `net/http/pprof`
- Добавлен отдельный debug endpoint:
  - `http://localhost:6060/debug/pprof/`
- Проверены:
  - CPU profile
  - heap profile
  - trace

### 3. Benchmark'и
Добавлены benchmark'и для двух уровней.

#### Service layer
- `BenchmarkOrderService_GetOrderByUID_CacheHit`
- `BenchmarkOrderService_GetOrderByUID_CacheMiss`
- `BenchmarkOrderService_GetOrderByUIDWithRefresh`

#### Handler layer
- `BenchmarkOrderHandler_GetOrder_CacheHit`
- `BenchmarkOrderHandler_GetOrder_CacheMiss`
- `BenchmarkOrderHandler_GetOrder_WithRefresh`

### 4. Исправления корректности
- Нормализована ошибка `ErrOrderNotFound`
- Исправлен docker/migration path
- Проверен рабочий ingestion path через Kafka → PostgreSQL → API

### 5. Оптимизации
Оставлены в финальной версии:
- `cache.Get()` переведен с `Lock` на `RLock`
- Убран лишний hot-path logging при нагрузочном профилировании
- `GetOrderByUIDWithRefresh()` больше не делает двойное чтение из БД
- `OrderRepository.GetOrderByUID()` сокращен с нескольких запросов до схемы:
  - `JOIN` для `orders + delivery + payment`
  - отдельный запрос для `items`

### 6. Что пробовалось и было откатано
- Пробовалось кэширование уже сериализованного JSON для cache-hit path
- По `pprof` это уменьшало вклад `encoding/json`
- Но benchmark'и показали ухудшение по времени и памяти
- Поэтому решение было откатано

---

## Промежуточные результаты

### Benchmark service
Текущие значения:
- CacheHit: `229.7 ns/op`, `464 B/op`, `2 allocs/op`
- CacheMiss: `494.9 ns/op`, `1200 B/op`, `5 allocs/op`
- WithRefresh: `100.0 ns/op`, `16 B/op`, `1 allocs/op`

### Benchmark handlers
Текущие значения:
- CacheHit: `5506 ns/op`, `7842 B/op`, `25 allocs/op`
- CacheMiss: `4930 ns/op`, `7841 B/op`, `25 allocs/op`
- WithRefresh: `5235 ns/op`, `8274 B/op`, `28 allocs/op`

### CPU profile
На финальном варианте:
- cache-hit path в основном упирается в `net/http` и syscalls записи ответа
- refresh-path в основном упирается в PostgreSQL path (`pgx`) и I/O
- внутренние лишние узкие места были убраны

---

## Что осталось сделать
- Финально оформить историю коммитов
- Собрать основной `README.md`
- При необходимости добавить `benchstat` сравнение в итоговую документацию
- Финально описать выводы по `trace`
