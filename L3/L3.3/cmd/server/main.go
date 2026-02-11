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

	"comment/internal/config"
	"comment/internal/handlers"
	"comment/internal/service"
	"comment/internal/storage"

	"github.com/wb-go/wbf/ginext"
)

func main() {
	// Загрузка конфигурации
	cfg, err := config.LoadConfig("config.yaml")
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	// Проверка корректности конфигурации
	if err := config.ValidateConfig(cfg); err != nil {
		log.Fatalf("failed to validate config: %v", err)
	}

	// Формирование адреса сервера
	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)

	// Создание нового экземпляра HTTP-сервера с режимом debug
	r := ginext.New("debug")

	//разрешаем запросы с любого источника
	r.Use(func(c *ginext.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// Обработка preflight-запросов
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204) // Возвращаем 204 и прерываем обработку
			return
		}

		// Передача управления следующему роуту
		c.Next()
	})

	// Встроенные middleware Gin
	r.Use(ginext.Logger(), ginext.Recovery())

	// Создание in-memory хранилища для комментариев
	store := storage.NewMemoryStorage()

	// Создание сервиса для работы с комментариями
	svc := service.NewCommentService(store)

	// Создание хендлера, который связывает HTTP-запросы с сервисом
	h := handlers.NewCommentHandler(svc)
	// Регистрация запрос
	h.Register(r)

	// оздаём свой http.Server
	srv := &http.Server{
		Addr:    addr,
		Handler: r,
	}

	// Запускаем сервер в отдельной горутине
	go func() {
		log.Printf("Starting server at %s...", addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// Канал для сигналов ОС
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Блокируемся до получения сигнала
	<-quit
	log.Println("Shutdown signal received...")

	// Даём серверу 5 секунд на завершение запросов
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Запускаем graceful shutdown
	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("server forced to shutdown: %v\n", err)
	}

	log.Println("server stopped gracefully")
}
