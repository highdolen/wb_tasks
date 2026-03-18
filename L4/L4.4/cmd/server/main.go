package main

import (
	"log"

	"gcMetrics/internal/app"
)

func main() {
	// Точка сбора приложения
	application, err := app.New()
	if err != nil {
		log.Fatalf("failed to initialize app: %v", err)
	}

	// Точка запуска
	if err := application.Run(); err != nil {
		log.Fatalf("application stopped with error: %v", err)
	}
}
