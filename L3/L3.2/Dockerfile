# Берём Go образ
FROM golang:1.24-alpine AS builder

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем модули и делаем кеш для сборки
COPY go.mod go.sum ./
RUN go mod download

# Копируем весь проект
COPY . .

# Собираем бинарник
RUN go build -o shortener cmd/server/main.go

# Второй слой для лёгкого контейнера
FROM alpine:3.18

WORKDIR /app

# Копируем бинарник из builder
COPY --from=builder /app/shortener .

# Экспортируем порт, который слушает наш сервис
EXPOSE 8080

# Команда запуска
CMD ["./shortener"]
