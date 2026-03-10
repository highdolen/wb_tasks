package main

import (
	"context"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"warehouseControl/internal/audit"
	"warehouseControl/internal/auth"
	"warehouseControl/internal/config"
	"warehouseControl/internal/export"
	"warehouseControl/internal/handlers"
	"warehouseControl/internal/items"
	"warehouseControl/internal/storage"

	"github.com/wb-go/wbf/ginext"
)

func main() {

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	cfg, err := config.LoadConfig("config.yaml")
	if err != nil {
		log.Fatalf("config error: %v", err)
	}

	if err := config.ValidateConfig(cfg); err != nil {
		log.Fatalf("config validation error: %v", err)
	}

	db, err := storage.NewPostgres(cfg)
	if err != nil {
		log.Fatalf("db connection error: %v", err)
	}
	defer func() {
		if err := db.Master.Close(); err != nil {
			log.Printf("db close error: %v", err)
		}
	}()

	// repositories
	itemRepo := items.NewRepository(db)
	auditRepo := audit.NewRepository(db)
	exportRepo := export.NewRepository(db)

	// services
	itemService := items.NewService(itemRepo)
	auditService := audit.NewService(auditRepo)
	exportService := export.NewService(exportRepo)
	authService := auth.NewService(cfg.Auth.JWTSecret)

	// handlers
	authHandler := handlers.NewAuthHandler(authService)
	itemsHandler := handlers.NewItemsHandler(itemService)
	auditHandler := handlers.NewAuditHandler(auditService)
	exportHandler := handlers.NewExportHandler(exportService)

	h := &handlers.Handlers{
		Auth:   authHandler,
		Items:  itemsHandler,
		Audit:  auditHandler,
		Export: exportHandler,
	}

	engine := ginext.New("debug")
	engine.Use(ginext.Logger(), ginext.Recovery())

	// API routes
	handlers.RegisterRoutes(engine, h, cfg.Auth.JWTSecret)

	// FRONTEND
	engine.Static("/static", "./frontend/static")
	engine.Static("/pages", "./frontend/pages")
	engine.StaticFile("/", "./frontend/index.html")

	srv := &http.Server{
		Addr:    ":8080",
		Handler: engine.Engine,
	}

	go func() {
		log.Println("server started on :8080")

		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server error: %v", err)
		}
	}()

	<-ctx.Done()

	log.Println("shutdown signal received")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Printf("server shutdown error: %v", err)
	}

	log.Println("server stopped")
}
