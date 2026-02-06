FROM golang:1.24-alpine

WORKDIR /app

# Сначала зависимости (для кеша)
COPY go.mod go.sum ./
RUN go mod download

# Копируем весь проект
COPY . .

# Собираем бинарник
RUN go build -o server ./cmd/server

EXPOSE 8080

CMD ["./server"]
