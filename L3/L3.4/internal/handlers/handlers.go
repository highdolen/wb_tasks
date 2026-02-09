package handlers

import (
	"io"
	"log"
	"net/http"

	"imageProcessor/internal/service"

	"github.com/google/uuid"
	"github.com/wb-go/wbf/ginext"
)

// Handler — HTTP-слой приложения
type Handler struct {
	service *service.ImageService
}

// New - создание нового HTTP-обработчика с зависимостью на ImageService
func New(service *service.ImageService) *Handler {
	return &Handler{service: service}
}

// Upload - принятие файла от пользователя, его сохранение и обработка
func (h *Handler) Upload(c *ginext.Context) {
	// Получаем файл из формы
	file, _, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, ginext.H{"error": "file is required"})
		return
	}
	defer func() {
		if err := file.Close(); err != nil {
			log.Println("failed to close file:", err)
		}
	}()

	// Генерируем уникальный ID для изображения
	id := uuid.New().String() + ".jpg"

	// Сохраняем оригинал и ставим задачу в очередь
	if err := h.service.UploadImage(c.Request.Context(), id, file); err != nil {
		c.JSON(http.StatusInternalServerError, ginext.H{"error": err.Error()})
		return
	}

	if err := h.service.ProcessImage(id); err != nil {
		log.Println("failed to process image:", err)
		c.JSON(http.StatusInternalServerError, ginext.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusAccepted, ginext.H{"id": id, "status": "processed"})
}

// GetImage - возврат обработанного изображение пользователю
func (h *Handler) GetImage(c *ginext.Context) {
	id := c.Param("id")

	// Получаем обработанное изображение
	img, err := h.service.GetImage(id)
	if err != nil {
		c.JSON(http.StatusNotFound, ginext.H{"status": "processing"})
		return
	}
	defer func() {
		if err := img.Close(); err != nil {
			log.Println("failed to close image:", err)
		}
	}()

	// Отдаём изображение как HTTP-ответ
	c.Header("Content-Type", "image/jpeg")
	c.Status(http.StatusOK)
	_, _ = io.Copy(c.Writer, img)
}

// DeleteImage - удалени изображение (и оригинал, и обработанна версия)
func (h *Handler) DeleteImage(c *ginext.Context) {
	id := c.Param("id")

	if err := h.service.DeleteImage(id); err != nil {
		c.JSON(http.StatusInternalServerError, ginext.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, ginext.H{"status": "deleted"})
}
