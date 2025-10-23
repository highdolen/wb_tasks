package downloader

import (
	"fmt"
	"net/url"
	"path/filepath"
	"time"
)

// Config содержит настройки загрузчика
type Config struct {
	BaseURL       string        // Базовый URL для скачивания
	Depth         int           // Глубина рекурсии
	OutputDir     string        // Директория для сохранения
	MaxConcurrent int           // Максимум одновременных загрузок
	Timeout       time.Duration // Таймаут для HTTP-запросов
}

// Validate проверяет корректность конфигурации
func (c *Config) Validate() error {
	if c.BaseURL == "" {
		return fmt.Errorf("BaseURL не может быть пустым")
	}

	// Парсим URL чтобы убедиться в его корректности
	_, err := url.Parse(c.BaseURL)
	if err != nil {
		return fmt.Errorf("некорректный URL: %v", err)
	}

	if c.Depth < 1 {
		return fmt.Errorf("глубина должна быть >= 1")
	}

	if c.MaxConcurrent < 1 {
		return fmt.Errorf("количество одновременных загрузок должно быть >= 1")
	}

	if c.Timeout == 0 {
		c.Timeout = 30 * time.Second
	}

	return nil
}

// GetOutputPath возвращает полный путь для сохранения файла
func (c *Config) GetOutputPath(urlPath string) string {
	return filepath.Join(c.OutputDir, urlPath)
}
