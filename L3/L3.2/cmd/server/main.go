package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"shortener/internal/cache"
	"shortener/internal/config"
	"shortener/internal/db"
	"shortener/internal/handlers"
	"shortener/internal/repository"

	"github.com/gin-contrib/cors"
	"github.com/wb-go/wbf/ginext"
)

func main() {
	// Load config
	cfg, err := config.LoadConfig("config.yaml")
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}
	if err := config.ValidateConfig(cfg); err != nil {
		log.Fatalf("failed to validate config: %v", err)
	}

	// Redis
	store, err := cache.New(cfg)
	if err != nil {
		log.Fatal(err)
	}
	if store != nil {
		defer store.Close()
		log.Println("Redis enabled")
	} else {
		log.Println("Redis disabled")
	}

	// Postgres
	storage, err := db.New(cfg)
	if err != nil {
		log.Fatal(err)
	}
	defer storage.Master.Close()
	log.Println("Postgres connected")

	// Repository
	repo := repository.New(cfg, storage, store)
	// Handlers
	h := handlers.NewRepository(repo)

	// HTTP server
	engine := ginext.New("debug")
	engine.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5500"},
		AllowMethods:     []string{"GET", "POST", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type"},
		AllowCredentials: true,
	}))

	// Подключение роутов
	handlers.RegisterRoutes(engine, h)

	addr := cfg.Server.Host + ":" + strconv.Itoa(cfg.Server.Port)

	//graceful shutdown
	httpServer := &http.Server{
		Addr:    addr,
		Handler: engine,
	}

	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()

	go func() {
		log.Println("Server started on", addr)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server error: %v", err)
		}
	}()

	<-ctx.Done()
	log.Println("Shutdown signal received")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := httpServer.Shutdown(shutdownCtx); err != nil {
		log.Printf("graceful shutdown failed: %v", err)
	} else {
		log.Println("Server stopped gracefully")
	}
}
