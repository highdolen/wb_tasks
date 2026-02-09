package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"imageProcessor/internal/app"
	"imageProcessor/internal/config"
	"imageProcessor/internal/handlers"

	"github.com/wb-go/wbf/ginext"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Загрузка конфигурации
	cfg, err := config.LoadConfig("config.yaml")
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	if err := config.ValidateConfig(cfg); err != nil {
		log.Fatalf("failed to validate config: %v", err)
	}

	// Инициализация приложения (storage, broker, service)
	application, err := app.New(ctx, cfg)
	if err != nil {
		log.Fatalf("failed to init app: %v", err)
	}

	log.Println("app started successfully")

	// Gin engine
	engine := ginext.New("debug")
	engine.Use(ginext.Logger(), ginext.Recovery())

	// Регистрация API-хендлеров
	h := handlers.New(application.Service)
	handlers.RegisterRoutes(engine, h)

	// Статика (CSS/JS)
	engine.Static("/static", "./frontend")

	// Главная страница
	engine.LoadHTMLGlob("frontend/*")
	engine.GET("/", func(c *ginext.Context) {
		c.HTML(200, "index.html", nil)
	})

	addr := cfg.Server.Host + ":" + strconv.Itoa(cfg.Server.Port)

	// Запуск HTTP сервера
	go func() {
		log.Println("HTTP server started at", addr)
		if err := engine.Run(addr); err != nil {
			log.Fatalf("http server error: %v", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("shutdown signal received")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	// Останавливаем фоновые воркеры, Kafka и т.д.
	application.Close()

	// Ждём завершения
	<-shutdownCtx.Done()
	log.Println("server stopped gracefully")
}
