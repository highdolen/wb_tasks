package main

import (
	"context"
	"log"
	"net/http"
	"os/signal"
	"salesTracker/internal/config"
	"salesTracker/internal/handlers"
	"salesTracker/internal/repository"
	"salesTracker/internal/service"
	"salesTracker/internal/storage"
	"syscall"
	"time"

	"github.com/wb-go/wbf/ginext"
)

func main() {
	// Контекст для graceful shutdown
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// Загрузка конфига
	cfg, err := config.LoadConfig("config.yaml")
	if err != nil {
		log.Fatalf("config error: %v", err)
	}

	// Валидация конфига
	if err := config.ValidateConfig(cfg); err != nil {
		log.Fatalf("config validation error: %v", err)
	}

	// Подключение к БД
	db, err := storage.NewPostgres(cfg)
	if err != nil {
		log.Fatalf("db error: %v", err)
	}
	defer func() {
		if err := db.Master.Close(); err != nil {
			log.Printf("db close error: %v", err)
		}
	}()

	// Репозиторий
	itemRepo := repository.NewItemRepository(db)

	// Сервисы
	itemService := service.NewItemService(itemRepo)
	analyticsService := service.NewAnalyticsService(itemRepo)
	exportService := service.NewExportService(itemRepo)

	// Gin Engine
	engine := ginext.New("debug")
	engine.Use(ginext.Logger(), ginext.Recovery())

	// Handlers
	itemHandler := handlers.NewItemHandler(itemService)
	analyticsHandler := handlers.NewAnalyticsHandler(analyticsService)
	exportHandler := handlers.NewExportHandler(exportService)

	// API
	handlers.SetupRouter(
		engine,
		itemHandler,
		analyticsHandler,
		exportHandler,
	)

	// Фронт
	engine.Static("/static", "./frontend/static")
	engine.LoadHTMLGlob("frontend/*.html")

	engine.GET("/", func(c *ginext.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})

	// HTTP сервер
	srv := &http.Server{
		Addr:    ":8080",
		Handler: engine.Engine,
	}

	// Запуск сервера
	go func() {
		log.Println("server started on :8080")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server error: %v", err)
		}
	}()

	// Ожидание сигнала завершения
	<-ctx.Done()
	log.Println("shutdown signal received...")

	// Контекст для корректного завершения
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Printf("server shutdown error: %v", err)
	}

	log.Println("server exited gracefully")
}
