package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"secondBlock/L2.18/internal/calendar"
	"secondBlock/L2.18/internal/config"
	"secondBlock/L2.18/internal/handlers"
)

func main() {
	// Загружаем конфиг
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	// Инициализируем календарь и хэндлеры
	cal := calendar.New()
	h := handlers.NewHandlers(cal)
	router := h.SetupRoutes()

	router.Use(handlers.LoggingMiddleware)

	srv := &http.Server{
		Addr:    ":" + strconv.Itoa(cfg.Port),
		Handler: router,
	}

	// graceful shutdown
	idleConnsClosed := make(chan struct{})
	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		<-sigChan

		log.Println("Received signal, shutting down...")

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := srv.Shutdown(ctx); err != nil {
			log.Printf("HTTP server Shutdown: %v", err)
		}
		close(idleConnsClosed)
	}()

	log.Printf("Server starting on port %d", cfg.Port)
	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(
			err, http.ErrServerClosed,
		) {
			log.Fatalf("ListenAndServe: %v", err)
		}
	}()

	<-idleConnsClosed
	log.Println("Server stopped")
}
