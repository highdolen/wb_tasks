package storage

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
)

type Storage struct {
	originalPath  string
	processedPath string
}

// New - инициализация Storage
func New(originalPath, processedPath string) (*Storage, error) {
	// создаём папки если их нет
	if err := os.MkdirAll(originalPath, 0755); err != nil {
		return nil, fmt.Errorf("create original dir: %w", err)
	}

	if err := os.MkdirAll(processedPath, 0755); err != nil {
		return nil, fmt.Errorf("create processed dir: %w", err)
	}

	return &Storage{
		originalPath:  originalPath,
		processedPath: processedPath,
	}, nil
}

// SaveOriginal - сохранение оригинального изображения в хранилище
func (s *Storage) SaveOriginal(id string, r io.Reader) (string, error) {
	path := filepath.Join(s.originalPath, id)

	file, err := os.Create(path)
	if err != nil {
		return "", fmt.Errorf("create original file: %w", err)
	}
	defer func() {
		if err := file.Close(); err != nil {
			log.Println("failed to close file:", err)
		}
	}()

	if _, err := io.Copy(file, r); err != nil {
		return "", fmt.Errorf("write original file: %w", err)
	}

	return path, nil
}

// SaveProcessed - сохранения обработанного изображения в хранилище
func (s *Storage) SaveProcessed(id string, r io.Reader) (string, error) {
	path := filepath.Join(s.processedPath, id)

	file, err := os.Create(path)
	if err != nil {
		return "", fmt.Errorf("create processed file: %w", err)
	}
	defer func() {
		if err := file.Close(); err != nil {
			log.Println("failed to close file:", err)
		}
	}()

	if _, err := io.Copy(file, r); err != nil {
		return "", fmt.Errorf("write processed file: %w", err)
	}

	return path, nil
}

// GetOriginal - получение оригинального изображения
func (s *Storage) GetOriginal(id string) (*os.File, error) {
	path := filepath.Join(s.originalPath, id)
	return os.Open(path)
}

// GetProcessed - получение обработанного изображения
func (s *Storage) GetProcessed(id string) (*os.File, error) {
	path := filepath.Join(s.processedPath, id)
	return os.Open(path)
}

// Delete - удаление изображения полностью
func (s *Storage) Delete(id string) error {
	orig := filepath.Join(s.originalPath, id)
	proc := filepath.Join(s.processedPath, id)

	_ = os.Remove(orig)
	_ = os.Remove(proc)

	return nil
}
