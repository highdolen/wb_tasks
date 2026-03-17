package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"calendar/internal/calendar"
	"calendar/internal/config"
	"calendar/internal/handlers"
	"calendar/internal/logger"
	"calendar/internal/reminder"
)

// main — точка входа в приложение
func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	logg := logger.New(100)
	reminderService := reminder.New(100, logg)
	cal := calendar.New(reminderService)

	h := handlers.NewHandlers(cal)
	middleware := handlers.NewMiddleware(logg)

	router := h.SetupRoutes()
	router.Use(middleware.Logging)

	srv := &http.Server{
		Addr:    ":" + strconv.Itoa(cfg.Port),
		Handler: router,
	}

	// Фоновая горутина архиватора
	go func() {
		cal.ArchiveOldEvents()
		logg.Log("archive worker started")

		ticker := time.NewTicker(1 * time.Minute)
		defer ticker.Stop()

		for range ticker.C {
			cal.ArchiveOldEvents()
			logg.Log("archive worker executed")
		}
	}()

	idleConnsClosed := make(chan struct{})

	// Горутина graceful shutdown
	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

		<-sigChan
		logg.Log("received shutdown signal")

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := srv.Shutdown(ctx); err != nil {
			logg.Log(fmt.Sprintf("HTTP shutdown error: %v", err))
		}

		close(idleConnsClosed)
	}()

	logg.Log(fmt.Sprintf("server starting on port %d", cfg.Port))

	// Отдельная горутина запускает HTTP-сервер
	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logg.Log(fmt.Sprintf("ListenAndServe error: %v", err))
			close(idleConnsClosed)
		}
	}()

	<-idleConnsClosed
	logg.Log("server stopped")
}
