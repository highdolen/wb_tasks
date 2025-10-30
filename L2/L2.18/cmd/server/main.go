package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"secondBlock/L2.18/internal/config"
)

func main() {
	ctx := context.Background()

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Ошибка в загрузке конфига: %v", err)
	}

	if err := cfg.Validate(); err != nil {
		log.Fatalf("Ошибка валидации конфига: %v", err)
	}

	// Для примера — просто заглушка
	mux := http.NewServeMux()
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "OK")
	})

	serverAddr := ":" + cfg.Server.Port
	log.Printf("Сервер запущен на %s", serverAddr)

	if err := http.ListenAndServe(serverAddr, mux); err != nil {
		log.Fatalf("Ошибка при запуске сервера: %v", err)
	}

	<-ctx.Done()
}
