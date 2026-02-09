package service

import (
	"context"
	"fmt"
	"io"
	"log"
	"mime"
	"path/filepath"

	"imageProcessor/internal/broker"
	"imageProcessor/internal/storage"
)

// ImageService — основной слой бизнес-логики приложения
type ImageService struct {
	storage   *storage.Storage
	broker    *broker.Broker
	processor *Processor
}

// New - создание нового экземпляра ImageService и инициализация Processor
func New(storage *storage.Storage, broker *broker.Broker) *ImageService {
	return &ImageService{
		storage:   storage,
		broker:    broker,
		processor: NewProcessor(),
	}
}

// UploadImage - сохранение оригинальное изображение в хранилище
func (s *ImageService) UploadImage(ctx context.Context, id string, r io.Reader) error {
	// Сохраняем оригинальный файл
	if _, err := s.storage.SaveOriginal(id, r); err != nil {
		return fmt.Errorf("save original: %w", err)
	}

	// Публикуем задачу в Kafka
	if err := s.broker.PublishTask(ctx, id); err != nil {
		return fmt.Errorf("publish task: %w", err)
	}

	return nil
}

// ProcessImage - выполнение обработки изображения
func (s *ImageService) ProcessImage(id string) error {
	// Открываем оригинальное изображение из хранилища
	origFile, err := s.storage.GetOriginal(id)
	if err != nil {
		return fmt.Errorf("open original: %w", err)
	}
	defer func() {
		if err := origFile.Close(); err != nil {
			log.Println("failed to close original file:", err)
		}
	}()

	// Определяем формат файла по расширению
	format := detectFormat(id)

	// Передаём изображение в Processor для resize + watermark
	buf, err := s.processor.Process(origFile, format)
	if err != nil {
		return fmt.Errorf("process: %w", err)
	}

	// Сохраняем обработанное изображение в отдельную директорию
	if _, err := s.storage.SaveProcessed(id, buf); err != nil {
		return fmt.Errorf("save processed: %w", err)
	}

	log.Println("image processed:", id)
	return nil
}

// GetImage - возвращение обработанное изображение из хранилища
func (s *ImageService) GetImage(id string) (io.ReadCloser, error) {
	return s.storage.GetProcessed(id)
}

// DeleteImage - удаление оригинала и обработанного изображения
func (s *ImageService) DeleteImage(id string) error {
	return s.storage.Delete(id)
}

// detectFormat -  формат изображения по расширению файла
func detectFormat(filename string) string {
	ext := filepath.Ext(filename)

	switch mime.TypeByExtension(ext) {
	case "image/png":
		return "png"
	case "image/gif":
		return "gif"
	default:
		return "jpeg"
	}
}
