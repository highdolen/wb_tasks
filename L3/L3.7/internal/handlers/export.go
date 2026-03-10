package handlers

import (
	"net/http"

	"warehouseControl/internal/export"

	"github.com/wb-go/wbf/ginext"
)

// ExportHandler - обработчик экспорта данных
type ExportHandler struct {
	service *export.Service
}

// NewExportHandler - создание обработчика экспорта
func NewExportHandler(service *export.Service) *ExportHandler {
	return &ExportHandler{service: service}
}

// ExportHistoryCSV - экспорт истории изменений в CSV файл
func (h *ExportHandler) ExportHistoryCSV(c *ginext.Context) {

	c.Header("Content-Type", "text/csv")
	c.Header("Content-Disposition", "attachment; filename=history.csv")

	err := h.service.ExportHistoryCSV(c.Request.Context(), c.Writer)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ginext.H{
			"error": err.Error(),
		})
		return
	}
}
