package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
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
	cfg, err := config.LoadConfig("config.yaml")
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	if err := config.ValidateConfig(cfg); err != nil {
		log.Fatalf("failed to validate config: %v", err)
	}

	rabbitClient, err := rabbitmq.NewRabbitClient(*cfg)
	if err != nil {
		log.Fatalf("failed to init RabbitMQ client: %v", err)
	}

	producer := rabbitmq.NewProducer(rabbitClient)

	var store storage.NotificationStorage
	if cfg.Redis.URL != "" {
		redisClient := wbfredis.New(cfg.Redis.URL, "", 0)
		store = storage.NewRedisStorage(redisClient, 24*time.Hour)
		log.Println("Redis storage enabled")
	} else {
		store = storage.NewMemoryStorage()
		log.Println("In-memory storage enabled")
	}

	emailSender := sender.NewEmailSender(
		cfg.SMTP.Host,
		cfg.SMTP.Port,
		cfg.SMTP.Username,
		cfg.SMTP.Password,
	)

	telegramSender := &sender.TelegramSender{
		BotToken: cfg.Telegram.Token,
	}

	// Контекст для graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	consumer := rabbitmq.NewNotificationConsumer(
		rabbitClient.Rmq,
		store,
		emailSender,
		telegramSender,
	)

	go func() {
		if err := consumer.Start(ctx); err != nil {
			log.Printf("consumer stopped: %v", err)
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

	r.GET("/", func(c *gin.Context) {
		c.File("./frontend/index.html")
	})
	r.Static("/frontend", "./frontend")
	handlers.RegisterRoutes(r, h)

	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)

	srv := &http.Server{
		Addr:    addr,
		Handler: r,
	}

	// Запуск HTTP сервера в горутине
	go func() {
		log.Printf("server started at %s", addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// Gracegul shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("shutdown signal received")

	// Останавливаем consumer через контекст
	cancel()

	// Даем HTTP серверу время корректно завершить соединения
	ctxTimeout, cancelTimeout := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelTimeout()

	if err := srv.Shutdown(ctxTimeout); err != nil {
		log.Printf("server shutdown failed: %v", err)
	}

	// Закрываем RabbitMQ соединение
	if err := rabbitClient.Rmq.Close(); err != nil {
		log.Printf("failed to close RabbitMQ: %v", err)
	}

	log.Println("graceful shutdown completed")
}
