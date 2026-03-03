package handlers

import (
	"net/http"
	"time"

	"salesTracker/internal/service"

	"github.com/wb-go/wbf/ginext"
)

// AnalyticsHandler структура
type AnalyticsHandler struct {
	service *service.AnalyticsService
}

// NewAnalyticsHandler - конструктор AnalyticsHandler
func NewAnalyticsHandler(s *service.AnalyticsService) *AnalyticsHandler {
	return &AnalyticsHandler{service: s}
}

// GetAnalytics - возврат аналитики
func (h *AnalyticsHandler) GetAnalytics(c *ginext.Context) {
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

	data, err := h.service.GetAnalytics(c.Request.Context(), from, to)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ginext.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, data)
}

// GetGroupedAnalytics - возврат сгрупированной аналитики
func (h *AnalyticsHandler) GetGroupedAnalytics(c *ginext.Context) {
	fromStr := c.Query("from")
	toStr := c.Query("to")
	groupBy := c.Query("group_by")
	sort := c.Query("sort")

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

	data, err := h.service.GetGroupedAnalytics(
		c.Request.Context(),
		from,
		to,
		groupBy,
		sort,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ginext.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, data)
}
