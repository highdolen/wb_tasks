package handlers

import (
	"net/http"
	"strconv"

	"eventBooker/internal/booking"
	"eventBooker/internal/event"
	"eventBooker/internal/user"

	"github.com/wb-go/wbf/ginext"
)

type UserHandler struct {
	userService    *user.Service
	bookingService *booking.Service
	eventService   *event.Service
}

// NewUserHandler - конструктор UserHandler
func NewUserHandler(u *user.Service, b *booking.Service, e *event.Service) *UserHandler {
	return &UserHandler{
		userService:    u,
		bookingService: b,
		eventService:   e,
	}
}

// POST /users — регистрация или получение пользователя
type userRequest struct {
	Name       string `json:"name"`
	Email      string `json:"email"`
	TelegramID string `json:"telegram_id"`
	Role       string `json:"role"`
}

// Register - ручка регистрации пользователя
func (h *UserHandler) Register(c *ginext.Context) {
	var req userRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ginext.H{"error": "invalid request"})
		return
	}

	u, err := h.userService.GetOrCreate(
		c.Request.Context(),
		req.Email,
		req.Name,
		req.TelegramID,
		req.Role,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ginext.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, ginext.H{
		"id":          u.ID,
		"name":        u.Name,
		"email":       u.Email,
		"telegram_id": u.TelegramID,
		"role":        u.Role,
	})
}

// GET /users/:id/bookings — получение всех броней пользователя
func (h *UserHandler) GetBookings(c *ginext.Context) {
	userIDParam := c.Param("id")

	userID, err := strconv.ParseInt(userIDParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, ginext.H{"error": "invalid user id"})
		return
	}

	bookingsList, err := h.bookingService.GetUserBookings(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ginext.H{"error": err.Error()})
		return
	}

	var result []ginext.H
	for _, b := range bookingsList {
		ev, _ := h.eventService.GetEvent(c.Request.Context(), b.EventID)

		result = append(result, ginext.H{
			"id":         b.ID,
			"event_id":   b.EventID,
			"event_name": ev.Name,
			"status":     b.Status,
			"created_at": b.CreatedAt,
			"expires_at": b.ExpiresAt,
		})
	}

	c.JSON(http.StatusOK, result)
}
