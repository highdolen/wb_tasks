package main

import (
	"context"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"syscall"
	"time"

	"L0_optimize/internal/cache"
	"L0_optimize/internal/config"
	"L0_optimize/internal/database"
	"L0_optimize/internal/handlers"
	"L0_optimize/internal/kafka"
	"L0_optimize/internal/service"
	"L0_optimize/internal/web"

	"github.com/gorilla/mux"
)

func main() {
	ctx := context.Background()

	// Загружаем конфиг
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Ошибка загрузки конфигурации: %v", err)
	}
	if err := cfg.Validate(); err != nil {
		log.Fatalf("Ошибка валидации конфигурации: %v", err)
	}

	// Подключение к базе
	dsn := "postgres://" + cfg.DB.User + ":" + cfg.DB.Password + "@" + cfg.DB.Host + ":" + cfg.DB.Port + "/" + cfg.DB.Name + "?sslmode=" + cfg.DB.SSLMode
	db, err := database.ConnectDB(dsn)
	if err != nil {
		log.Fatalf("Ошибка подключения к БД: %v", err)
	}
	defer db.Close()

	repo := database.NewOrderRepository(db)

	// Создаём кэш с TTL 30 минут
	orderCache := cache.New(30 * time.Minute)

	// Создаём адаптеры для сервисного слоя
	repoAdapter := service.NewRepositoryAdapter(repo)
	cacheAdapter := service.NewCacheAdapter(orderCache)

	// Создаём сервис заказов
	orderService := service.NewOrderService(repoAdapter, cacheAdapter)

	// Загружаем данные из БД в кэш через адаптер
	if err := orderService.LoadFromDB(ctx); err != nil {
		log.Fatalf("Ошибка загрузки кэша: %v", err)
	}
	log.Println("Кэш успешно загружен")

	// Создаём Kafka Consumer
	consumer := kafka.NewConsumer(
		[]string{cfg.Kafka.Broker},
		"orders",
		"group-1",
		repo,
		orderCache,
	)

	// Создаём контекст для graceful shutdown
	ctxWithCancel, cancel := context.WithCancel(ctx)
	defer cancel()

	// Подключаем handlers
	r := mux.NewRouter()
	orderHandler := handlers.NewOrderHandler(orderService)

	// API для работы с заказами
	r.HandleFunc("/order/{order_uid}", orderHandler.GetOrder).Methods("GET", "OPTIONS")

	// API для управления кешом
	r.HandleFunc("/cache/stats", orderHandler.GetCacheStats).Methods("GET", "OPTIONS")
	r.HandleFunc("/cache/invalidate/{order_uid}", orderHandler.InvalidateCache).Methods("POST", "DELETE", "OPTIONS")
	r.HandleFunc("/cache/invalidate", orderHandler.InvalidateCache).Methods("POST", "DELETE", "OPTIONS")

	err = r.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		path, _ := route.GetPathTemplate()
		methods, _ := route.GetMethods()
		log.Printf("Route registered: %s Methods: %v", path, methods)
		return nil
	})
	if err != nil {
		log.Printf("Ошибка при обходе маршрутов: %v", err)
	}

	// Подключаем веб-интерфейс
	web.RegisterWebHandlers(r)

	// Middleware - добавляет лишнюю нагрузку на CPU
	//r.Use(handlers.LoggingMiddleware)
	r.Use(handlers.CORSMiddleware)

	// Основной HTTP сервер
	srv := &http.Server{
		Addr:    cfg.Server.Port,
		Handler: r,
	}

	// Отдельный debug server для pprof
	pprofSrv := &http.Server{
		Addr: ":6060",
	}

	// Канал для перехвата сигналов
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Канал для уведомления о завершении shutdown
	shutdownComplete := make(chan bool, 1)

	// Горутина для обработки сигналов и graceful shutdown
	go func() {
		// Ожидаем сигнал
		sig := <-sigChan
		log.Printf("Получен сигнал %v, начинаем graceful shutdown...", sig)

		// Отменяем контекст для остановки Kafka consumer
		log.Println("Останавливаем Kafka consumer...")
		cancel()

		// Ждем немного, чтобы consumer успел обработать отмену контекста
		time.Sleep(100 * time.Millisecond)

		// Создаём контекст с таймаутом для shutdown HTTP серверов
		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer shutdownCancel()

		// Graceful shutdown основного HTTP сервера
		log.Println("Останавливаем HTTP сервер...")
		if err := srv.Shutdown(shutdownCtx); err != nil {
			log.Printf("Ошибка при graceful shutdown HTTP сервера: %v", err)
		} else {
			log.Println("HTTP сервер успешно остановлен")
		}

		// Graceful shutdown pprof сервера
		log.Println("Останавливаем pprof сервер...")
		if err := pprofSrv.Shutdown(shutdownCtx); err != nil {
			log.Printf("Ошибка при graceful shutdown pprof сервера: %v", err)
		} else {
			log.Println("pprof сервер успешно остановлен")
		}

		// Закрываем Kafka consumer
		log.Println("Закрываем Kafka consumer...")
		consumer.Close()
		log.Println("Kafka consumer успешно остановлен")

		// Закрываем кеш (останавливаем горутину очистки)
		log.Println("Останавливаем кеш...")
		orderCache.Close()
		log.Println("Кеш успешно остановлен")

		log.Println("Graceful shutdown завершён")
		shutdownComplete <- true
	}()

	// Запускаем Kafka Consumer в горутине
	go consumer.Start(ctxWithCancel)

	// Запускаем pprof сервер в горутине
	go func() {
		log.Println("pprof сервер запущен на :6060")
		if err := pprofSrv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("Ошибка pprof сервера: %v", err)
			sigChan <- syscall.SIGTERM
		}
	}()

	// Запускаем HTTP сервер в горутине
	go func() {
		log.Printf("HTTP сервер запущен на %s", cfg.Server.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("Ошибка HTTP сервера: %v", err)
			// Отправляем сигнал для shutdown, если сервер упал
			sigChan <- syscall.SIGTERM
		}
	}()

	log.Println("Сервер запущен. Нажмите Ctrl+C для остановки.")

	// Ожидаем завершения shutdown
	<-shutdownComplete
}
