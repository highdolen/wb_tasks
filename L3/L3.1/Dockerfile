# Stage 1: Build
FROM golang:1.24 AS builder

WORKDIR /app

# Копируем go.mod и go.sum
COPY go.mod go.sum ./

# Загружаем зависимости
RUN go mod download

# Копируем весь проект
COPY . .

# Собираем бинарник
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -o notifier ./cmd/server/main.go

# Stage 2: Runtime
FROM alpine:3.20

WORKDIR /app

# Копируем бинарник и конфиг
COPY --from=builder /app/notifier .
COPY config.yaml .

# Порт приложения
EXPOSE 8080

# Запуск приложения
CMD ["./notifier"]
