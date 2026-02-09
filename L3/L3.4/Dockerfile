# ---- Stage 1: Build ----
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Копируем go.mod и go.sum для кеширования слоев
COPY go.mod go.sum ./
RUN go mod download

# Копируем весь проект
COPY . ./

# Собираем бинарник
RUN go build -o image_processor ./cmd/server

# ---- Stage 2: Final ----
FROM alpine:latest

WORKDIR /app

# Копируем бинарник из builder
COPY --from=builder /app/image_processor ./

# Копируем фронт
COPY --from=builder /app/frontend ./frontend

# Копируем конфиг
COPY --from=builder /app/config.yaml ./config.yaml

# Порт приложения
EXPOSE 8080

# Запуск приложения
CMD ["./image_processor"]
