# GC Memory Analyzer

Утилита на Go для анализа памяти и сборщика мусора (GC) с экспортом
метрик в формате Prometheus.

## Возможности

-   Сбор метрик через `runtime.ReadMemStats`
-   Настройка GC через `debug.SetGCPercent`
-   HTTP endpoint `/metrics` (Prometheus format)
-   Профилирование через `pprof`
-   Healthcheck endpoint `/healthz`
-   Тестовый endpoint `/debug/alloc` для генерации нагрузки

## Архитектура проекта

    gcMetrics/
    │   config.yaml
    │   go.mod
    │   go.sum
    │   README.md
    │
    ├── cmd/
    │   └── server/
    │       └── main.go
    │
    ├── internal/
    │   ├── app/
    │   │   └── app.go
    │   │
    │   ├── config/
    │   │   ├── config.go
    │   │   └── loadConfig.go
    │   │
    │   ├── debug/
    │   │   └── alloc.go
    │   │
    │   ├── httpserver/
    │   │   └── server.go
    │   │
    │   ├── metrics/
    │   │   ├── collector.go
    │   │   └── prometheus.go
    │   │
    │   └── profiler/
    │       └── pprof.go
    │
    └── pkg/
        └── meminfo/
            └── meminfo.go

## Запуск

Из корня проекта:

    go run ./cmd/server

## Конфигурация

Файл `config.yaml` в корне проекта:

``` yaml
host: "0.0.0.0"
port: "9090"
metrics_path: "/metrics"
health_path: "/healthz"
gc_percent: 50
```

## Доступные endpoints

-   `/metrics` --- метрики Prometheus
-   `/healthz` --- проверка состояния
-   `/debug/pprof/` --- профилирование
-   `/debug/pprof/heap`
-   `/debug/pprof/profile`
-   `/debug/alloc?mb=50` --- генерация аллокаций

## Примеры запросов

### Проверка сервера

    curl http://localhost:9090/healthz

### Получение метрик

    curl http://localhost:9090/metrics

### Генерация нагрузки

    curl "http://localhost:9090/debug/alloc?mb=50"

### Снятие heap profile

    go tool pprof http://localhost:9090/debug/pprof/heap

### CPU profile

    go tool pprof "http://localhost:9090/debug/pprof/profile?seconds=5"

## Проверенные метрики

-   `go_gc_mallocs_total`
-   `go_gc_cycles_total`
-   `go_memory_alloc_bytes`
-   `go_gc_last_time_unix`
-   `go_gc_pause_total_ns`
-   `go_runtime_goroutines`

## Описание

Программа собирает статистику памяти и GC из runtime Go и предоставляет
её через HTTP endpoint в формате Prometheus. Поддерживается
профилирование и настройка GC.

## Примечание

Endpoint `/debug/alloc` используется только для тестирования и
демонстрации работы GC.
