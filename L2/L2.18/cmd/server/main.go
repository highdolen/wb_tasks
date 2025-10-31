package main

import (
	"log"
	"net/http"
	"strconv"

	"secondBlock/L2.18/internal/calendar"
	"secondBlock/L2.18/internal/config"
	"secondBlock/L2.18/internal/handlers"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	cal := calendar.New()
	h := handlers.NewHandlers(cal)
	router := h.SetupRoutes()

	// Add logging middleware
	routerWithMiddleware := handlers.LoggingMiddleware(router) // ← теперь работает!

	log.Printf("Server starting on port %d", cfg.Port)
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(cfg.Port), routerWithMiddleware))
}
