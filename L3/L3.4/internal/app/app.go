package app

import (
	"context"
	"log"

	"imageProcessor/internal/broker"
	"imageProcessor/internal/config"
	"imageProcessor/internal/service"
	"imageProcessor/internal/storage"
)

// App — главный контейнер приложения
type App struct {
	Service *service.ImageService
	Broker  *broker.Broker
}

// New - создание нового приложения с зависимостями
func New(ctx context.Context, cfg *config.AppConfig) (*App, error) {
	// storage
	st, err := storage.New(cfg.Storage.OriginalPath, cfg.Storage.ProcessedPath)
	if err != nil {
		return nil, err
	}

	// kafka broker
	br := broker.New(
		cfg.Kafka.Brokers,
		cfg.Kafka.Topic,
		cfg.Kafka.GroupID,
	)

	// service
	svc := service.New(st, br)

	// запуск воркера
	br.StartWorker(ctx, func(id string) {
		if err := svc.ProcessImage(id); err != nil {
			log.Println("failed to process image:", err)
		}
	})

	return &App{
		Service: svc,
		Broker:  br,
	}, nil
}

// Close - остановка фоновых воркеров и брокера сообщений
func (a *App) Close() {
	a.Broker.Close()
}
