package handlers

import (
	"net/http"
	"strconv"
	"time"

	"eventBooker/internal/event"

	"github.com/wb-go/wbf/ginext"
)

type EventHandler struct {
	service *event.Service
}

// NewEventHandler - конструктор EventHandler
func NewEventHandler(s *event.Service) *EventHandler {
	return &EventHandler{service: s}
}

// POST /events
func (h *EventHandler) CreateEvent(c *ginext.Context) {
	var req struct {
		Name       string `json:"name"`
		Date       string `json:"date"` // формат от <input type="datetime-local">
		TotalSeats int    `json:"total_seats"`
		BookingTTL int    `json:"booking_ttl_minutes"`
	}

	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ginext.H{"error": err.Error()})
		return
	}

	// Парсим дату от браузера
	date, err := time.Parse("2006-01-02T15:04", req.Date)
	if err != nil {
		c.JSON(http.StatusBadRequest, ginext.H{"error": "invalid date format"})
		return
	}

	now := time.Now()

	// событие должно быть минимум через 1 минуту
	if date.Before(now.Add(1 * time.Minute)) {
		c.JSON(http.StatusBadRequest, ginext.H{
			"error": "event must be scheduled at least 1 minute in the future",
		})
		return
	}

	// Проверка количества мест
	if req.TotalSeats <= 0 {
		c.JSON(http.StatusBadRequest, ginext.H{
			"error": "total seats must be greater than 0",
		})
		return
	}

	// Проверка TTL брони
	if req.BookingTTL <= 0 {
		c.JSON(http.StatusBadRequest, ginext.H{
			"error": "booking ttl must be greater than 0",
		})
		return
	}

	e, err := h.service.CreateEvent(
		c.Request.Context(),
		req.Name,
		date,
		req.TotalSeats,
		time.Duration(req.BookingTTL)*time.Minute,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ginext.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, e)
}

// GET /events/:id
func (h *EventHandler) GetEvent(c *ginext.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, ginext.H{"error": "invalid id"})
		return
	}

	e, err := h.service.GetEvent(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, ginext.H{"error": "event not found"})
		return
	}

	c.JSON(http.StatusOK, e)
}

// GET /events
func (h *EventHandler) GetAllEvents(c *ginext.Context) {
	events, err := h.service.GetAll(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, ginext.H{"error": err.Error()})
		return
	}

	// Если events nil, вернуть пустой массив
	if events == nil {
		events = []*event.Event{}
	}

	c.JSON(http.StatusOK, events)
}
