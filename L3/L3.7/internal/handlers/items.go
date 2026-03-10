package handlers

import (
	"net/http"
	"strconv"

	"warehouseControl/internal/items"

	"github.com/wb-go/wbf/ginext"
)

// ItemsHandler - обработчик запросов для работы с товарами
type ItemsHandler struct {
	service *items.Service
}

// NewItemsHandler - создание обработчика товаров
func NewItemsHandler(service *items.Service) *ItemsHandler {
	return &ItemsHandler{service: service}
}

// CreateItem - создание нового товара
func (h *ItemsHandler) CreateItem(c *ginext.Context) {

	username := c.GetString("username")

	var item items.Item
	if err := c.ShouldBindJSON(&item); err != nil {
		c.JSON(http.StatusBadRequest, ginext.H{"error": "invalid body"})
		return
	}

	if err := h.service.CreateItem(c.Request.Context(), &item, username); err != nil {
		c.JSON(http.StatusInternalServerError, ginext.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, item)
}

// GetItems - получение списка всех товаров
func (h *ItemsHandler) GetItems(c *ginext.Context) {

	itemsList, err := h.service.GetItems(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, ginext.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, itemsList)
}

// UpdateItem - обновление данных товара
func (h *ItemsHandler) UpdateItem(c *ginext.Context) {

	username := c.GetString("username")

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ginext.H{"error": "invalid id"})
		return
	}

	var item items.Item
	if err := c.ShouldBindJSON(&item); err != nil {
		c.JSON(http.StatusBadRequest, ginext.H{"error": "invalid body"})
		return
	}

	if err := h.service.UpdateItem(c.Request.Context(), id, &item, username); err != nil {
		c.JSON(http.StatusInternalServerError, ginext.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, ginext.H{"status": "updated"})
}

// DeleteItem - удаление товара
func (h *ItemsHandler) DeleteItem(c *ginext.Context) {

	username := c.GetString("username")

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ginext.H{"error": "invalid id"})
		return
	}

	if err := h.service.DeleteItem(c.Request.Context(), id, username); err != nil {
		c.JSON(http.StatusInternalServerError, ginext.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, ginext.H{"status": "deleted"})
}
