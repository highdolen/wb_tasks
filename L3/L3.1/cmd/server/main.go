package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"delayed_notifier/internal/config"
	"delayed_notifier/internal/handlers"
	"delayed_notifier/internal/rabbitmq"
	"delayed_notifier/internal/sender"
	"delayed_notifier/internal/storage"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/wb-go/wbf/ginext"
	wbfredis "github.com/wb-go/wbf/redis"
)

func main() {

	// Загрузка конфигурации
	cfg, err := config.LoadConfig("config.yaml")
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	if err := config.ValidateConfig(cfg); err != nil {
		log.Fatalf("failed to validate config: %v", err)
	}

	// Инициализация клиента RabbitMQ
	rabbitClient, err := rabbitmq.NewRabbitClient(*cfg)
	if err != nil {
		log.Fatalf("failed to init RabbitMQ client: %v", err)
	}
	defer rabbitClient.Rmq.Close()

	//Инициализация продюсера RabbitMQ
	producer := rabbitmq.NewProducer(rabbitClient)

	// Хранилище уведомлений
	var store storage.NotificationStorage

	if cfg.Redis.URL != "" {
		redisClient := wbfredis.New(cfg.Redis.URL, "", 0)
		store = storage.NewRedisStorage(redisClient, 24*time.Hour)
		log.Println("Redis storage enabled")
	} else {
		store = storage.NewMemoryStorage()
		log.Println("In-memory storage enabled")
	}

	// Senders
	emailSender := sender.NewEmailSender(
		cfg.SMTP.Host,
		cfg.SMTP.Port,
		cfg.SMTP.Username,
		cfg.SMTP.Password,
	)

	telegramSender := &sender.TelegramSender{
		BotToken: cfg.Telegram.Token,
	}

	// Consumer
	ctx := context.Background()

	consumer := rabbitmq.NewNotificationConsumer(
		rabbitClient.Rmq,
		store,
		emailSender,
		telegramSender,
	)

	go func() {
		if err := consumer.Start(ctx); err != nil {
			log.Fatalf("consumer failed: %v", err)
		}
	}()

	// HTTP server
	h := handlers.NewNotificationHandlers(store, producer)

	r := ginext.New("debug")

	r.Use(
		ginext.Logger(),
		ginext.Recovery(),
		cors.New(cors.Config{
			AllowOrigins: []string{"*"},
			AllowMethods: []string{"GET", "POST", "DELETE", "OPTIONS"},
			AllowHeaders: []string{"Origin", "Content-Type"},
			MaxAge:       12 * time.Hour,
		}),
	)

	// Главная страница
	r.GET("/", func(c *gin.Context) {
		c.File("./frontend/index.html")
	})

	// CSS + JS
	r.Static("/frontend", "./frontend")

	// API
	handlers.RegisterRoutes(r, h)

	// Запуск сервера
	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	log.Printf("server started at %s", addr)

	if err := r.Run(addr); err != nil {
		log.Fatal(err)
	}
}
