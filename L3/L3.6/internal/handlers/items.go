package handlers

import (
	"net/http"
	"strconv"
	"time"

	"salesTracker/internal/models"
	"salesTracker/internal/service"

	"github.com/wb-go/wbf/ginext"
)

// ItemHandle - обработка HTTP-запросов, связанных с транзакциями
type ItemHandler struct {
	service *service.ItemService
}

// NewItemHandler - конструктор ItemHandler
func NewItemHandler(s *service.ItemService) *ItemHandler {
	return &ItemHandler{service: s}
}

// Create - создение новой записи
func (h *ItemHandler) Create(c *ginext.Context) {
	var item models.Item

	if err := c.ShouldBindJSON(&item); err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	// Пользователь должен передать дату создания записи
	if item.CreatedAt.IsZero() {
		c.JSON(http.StatusBadRequest, ginext.H{"error": "created_at is required"})
		return
	}

	id, err := h.service.Create(c, &item)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, ginext.H{"id": id})
}

// GetAll - возвращение списка записей с возможностью фильтрации
func (h *ItemHandler) GetAll(c *ginext.Context) {
	fromStr := c.Query("from")
	toStr := c.Query("to")

	var fromTime, toTime *time.Time
	if fromStr != "" {
		t, err := time.Parse(time.RFC3339, fromStr)
		if err == nil {
			fromTime = &t
		}
	}
	if toStr != "" {
		t, err := time.Parse(time.RFC3339, toStr)
		if err == nil {
			toTime = &t
		}
	}

	filter := models.ItemFilter{
		Type:     c.Query("type"),
		Category: c.Query("category"),
		From:     fromTime,
		To:       toTime,
	}

	items, err := h.service.GetAll(c, filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, items)
}

// Delete - удаление записи по ID
func (h *ItemHandler) Delete(c *ginext.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)

	err := h.service.Delete(c, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, ginext.H{"status": "deleted"})
}

// Update - обновление существующей записи
func (h *ItemHandler) Update(c *ginext.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)

	var item models.Item
	if err := c.ShouldBindJSON(&item); err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	err := h.service.Update(c, id, &item)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, ginext.H{"status": "updated"})
}
