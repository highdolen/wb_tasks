package handlers

import (
	"net/http"
	"strconv"
	"warehouseControl/internal/audit"

	"github.com/wb-go/wbf/ginext"
)

// AuditHandler - обработчик запросов для истории изменений
type AuditHandler struct {
	service *audit.Service
}

// NewAuditHandler - создание обработчика истории
func NewAuditHandler(service *audit.Service) *AuditHandler {
	return &AuditHandler{
		service: service,
	}
}

// GetItemHistory - получение истории изменений конкретного товара
func (h *AuditHandler) GetItemHistory(c *ginext.Context) {

	idStr := c.Param("id")
	id, _ := strconv.Atoi(idStr)

	history, err := h.service.GetHistory(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ginext.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, history)
}

// FilterHistory - фильтрация истории изменений
func (h *AuditHandler) FilterHistory(c *ginext.Context) {

	user := c.Query("user")
	action := c.Query("action")
	from := c.Query("from")
	to := c.Query("to")

	history, err := h.service.FilterHistory(
		c.Request.Context(),
		user,
		action,
		from,
		to,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, ginext.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, history)
}

// GetDiff - получение различий между версиями товара
func (h *AuditHandler) GetDiff(c *ginext.Context) {

	idStr := c.Param("id")
	id, _ := strconv.Atoi(idStr)

	history, err := h.service.GetHistory(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ginext.H{"error": err.Error()})
		return
	}

	if len(history) == 0 {
		c.JSON(http.StatusNotFound, ginext.H{"error": "no history"})
		return
	}

	diff, err := h.service.GetDiff(history[0].OldValue, history[0].NewValue)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ginext.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, diff)
}
