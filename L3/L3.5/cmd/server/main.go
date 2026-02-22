package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"eventBooker/internal/booking"
	"eventBooker/internal/config"
	"eventBooker/internal/event"
	"eventBooker/internal/handlers"
	"eventBooker/internal/notification"
	"eventBooker/internal/scheduler"
	"eventBooker/internal/user"

	"github.com/wb-go/wbf/dbpg"
	"github.com/wb-go/wbf/ginext"
)

func main() {
	// Контекст для graceful shutdown
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// Загружаем конфиг
	cfg, err := config.LoadConfig("config.yaml")
	if err != nil {
		log.Fatalf("config load error: %v", err)
	}
	if err := config.ValidateConfig(cfg); err != nil {
		log.Fatalf("config validation error: %v", err)
	}

	// Подключение к базе
	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.Postgres.Host,
		cfg.Postgres.Port,
		cfg.Postgres.User,
		cfg.Postgres.Password,
		cfg.Postgres.DBName,
		cfg.Postgres.SSLMode,
	)

	db, err := dbpg.New(dsn, nil, nil)
	if err != nil {
		log.Fatalf("db connection error: %v", err)
	}
	defer func() {
		if err := db.Master.Close(); err != nil {
			log.Printf("db close error: %v", err)
		}
	}()

	// Репозитории
	eventRepo := event.NewRepository(db)
	bookingRepo := booking.NewRepository(db)
	userRepo := user.NewRepository(db)

	// Сервисы
	eventService := event.NewService(eventRepo)
	userService := user.NewService(userRepo)
	telegramClient := notification.NewTelegram(cfg.Telegram.BotToken)
	bookingService := booking.NewService(bookingRepo, telegramClient)

	// Планировщик для отмены просроченных броней
	expScheduler := scheduler.NewExpirationScheduler(bookingService, cfg.Scheduler.Interval)
	expScheduler.Start(ctx)

	// Gin engine
	engine := ginext.New("debug")
	engine.Use(func(c *ginext.Context) {
		c.Writer.Header().Set("Cache-Control", "no-store")
		c.Next()
	})
	engine.Use(ginext.Logger(), ginext.Recovery())

	// Handlers
	eventHandler := handlers.NewEventHandler(eventService)
	bookingHandler := handlers.NewBookingHandler(bookingService, userService, eventService, cfg.Telegram.BotToken)
	userHandler := handlers.NewUserHandler(userService, bookingService, eventService)

	handlers.RegisterRoutes(engine, eventHandler, bookingHandler, userHandler)

	// Подключение фронта
	engine.Static("/static", "./frontend/static")
	engine.LoadHTMLGlob("frontend/*.html")

	// Главная страница
	engine.GET("/", func(c *ginext.Context) {
		c.HTML(http.StatusOK, "login.html", nil)
	})

	//Страница пользователя
	engine.GET("/user", func(c *ginext.Context) {
		c.HTML(http.StatusOK, "user.html", nil)
	})

	//Страница админа
	engine.GET("/admin", func(c *ginext.Context) {
		c.HTML(http.StatusOK, "admin.html", nil)
	})

	// Запуск сервера
	addr := cfg.Server.Host + ":" + strconv.Itoa(cfg.Server.Port)
	srv := &http.Server{
		Addr:    addr,
		Handler: engine.Engine,
	}

	go func() {
		log.Printf("server started at %s", addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen error: %v", err)
		}
	}()

	//Graceful shutdown
	<-ctx.Done()
	log.Println("shutdown signal received")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_ = srv.Shutdown(shutdownCtx)
	log.Println("server exited gracefully")
}
