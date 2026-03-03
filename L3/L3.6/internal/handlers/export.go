package handlers

import (
	"net/http"
	"salesTracker/internal/service"
	"time"

	"github.com/wb-go/wbf/ginext"
)

type ExportHandler struct {
	service *service.ExportService
}

// NewExportHandler - конструктор ExportHandler
func NewExportHandler(s *service.ExportService) *ExportHandler {
	return &ExportHandler{service: s}
}

// ExportCSV - экспорт csv
func (h *ExportHandler) ExportCSV(c *ginext.Context) {
	fromStr := c.Query("from")
	toStr := c.Query("to")

	from, err := time.Parse(time.RFC3339, fromStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, ginext.H{"error": "invalid from date"})
		return
	}

	to, err := time.Parse(time.RFC3339, toStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, ginext.H{"error": "invalid to date"})
		return
	}

	c.Header("Content-Disposition", "attachment; filename=analytics.csv")
	c.Header("Content-Type", "text/csv")

	if err := h.service.ExportCSV(c.Request.Context(), from, to, c.Writer); err != nil {
		c.JSON(http.StatusInternalServerError, ginext.H{"error": err.Error()})
		return
	}
}
