package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Port int
}

func LoadConfig() (*Config, error) {
	// Загружаем .env файл (игнорируем ошибку, если файла нет)
	if err := godotenv.Load(); err != nil {
		fmt.Println("DEBUG: .env file not found in root directory, using environment variables")
	}

	portStr := os.Getenv("CALENDAR_PORT")

	// Если порт не установлен, используем значение по умолчанию
	if portStr == "" {
		portStr = "8080"
		fmt.Println("DEBUG: Using default port 8080")
	}

	// Конвертируем порт в int
	port, err := strconv.Atoi(portStr)
	if err != nil {
		return nil, fmt.Errorf("invalid port: %v", err)
	}

	return &Config{
		Port: port,
	}, nil
}
