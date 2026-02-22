package handlers

import (
	"net/http"
	"strconv"

	"eventBooker/internal/booking"
	"eventBooker/internal/event"
	"eventBooker/internal/user"

	"github.com/wb-go/wbf/ginext"
)

// BookingHandler - структура хэндлера бронирования
type BookingHandler struct {
	service      *booking.Service
	userService  *user.Service
	eventService *event.Service
	botToken     string
}

// NewBookingHandler - конструктор BookingHandler
func NewBookingHandler(
	b *booking.Service,
	u *user.Service,
	e *event.Service,
	botToken string,
) *BookingHandler {
	return &BookingHandler{
		service:      b,
		userService:  u,
		eventService: e,
		botToken:     botToken,
	}
}

// bookRequest - структура запроса
type bookRequest struct {
	Email      string `json:"email"`
	Name       string `json:"name"`
	TelegramID string `json:"telegram_id"`
	Role       string `json:"role"`
}

// Book - ручка бронирования события
func (h *BookingHandler) Book(c *ginext.Context) {
	eventIDStr := c.Param("id")
	eventID, err := strconv.ParseInt(eventIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, ginext.H{"error": "invalid event id"})
		return
	}

	var req bookRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ginext.H{"error": "invalid body"})
		return
	}

	u, err := h.userService.GetOrCreate(c.Request.Context(), req.Email, req.Name, req.TelegramID, req.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ginext.H{"error": err.Error()})
		return
	}

	ev, err := h.eventService.GetEvent(c.Request.Context(), eventID)
	if err != nil {
		c.JSON(http.StatusNotFound, ginext.H{"error": "event not found"})
		return
	}

	ttlSeconds := int64(ev.BookingTTL.Seconds())

	newBookingID, err := h.service.BookSeatWithID(c.Request.Context(), eventID, u.ID, ttlSeconds)
	if err != nil {
		c.JSON(http.StatusConflict, ginext.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, ginext.H{"booking_id": newBookingID})
}

// Confirm -ручка оплаты бронирования
func (h *BookingHandler) Confirm(c *ginext.Context) {
	bookingIDStr := c.Param("id")
	bookingID, err := strconv.ParseInt(bookingIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, ginext.H{"error": "invalid booking id"})
		return
	}

	if err := h.service.ConfirmBooking(c.Request.Context(), bookingID); err != nil {
		c.JSON(http.StatusConflict, ginext.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, ginext.H{"message": "booking confirmed"})
}
