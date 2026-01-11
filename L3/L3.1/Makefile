run:
	go run ./cmd/server/main.go

build:
	go build -o notifier ./cmd/server/main.go

docker-build:
	docker build -t notifier-l3.1 .

up:
	docker-compose up --build

down:
	docker-compose down

restart:
	docker-compose down && docker-compose up --build
