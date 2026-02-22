# STAGE 1: BUILD 
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Нужны для сборки некоторых пакетов
RUN apk add --no-cache git

# Кэшируем зависимости
COPY go.mod go.sum ./
RUN go mod download

# Копируем весь проект
COPY . .

# Собираем бинарник
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o eventbooker ./cmd/server


# STAGE 2: RUNTIME 
FROM alpine:latest

WORKDIR /app

# Для корректной работы time / ssl
RUN apk add --no-cache ca-certificates

# Копируем бинарник
COPY --from=builder /app/eventbooker .

# Копируем config
COPY --from=builder /app/config.yaml .

# Копируем frontend
COPY --from=builder /app/frontend ./frontend

# Открываем порт
EXPOSE 8080

# Запуск
CMD ["./eventbooker"]
