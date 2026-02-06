package main

import (
	"fmt"
	"log"

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

	// Создание новый экземпляр HTTP-сервера с режимом debug
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

	// Регистрация маршрутов (POST, GET, DELETE)
	h.Register(r)

	log.Printf("Starting server at %s...", addr)

	// Запускаем сервер на указанном адресе
	if err := r.Run(addr); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
