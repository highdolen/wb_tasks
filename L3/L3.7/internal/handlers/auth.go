package handlers

import (
	"net/http"
	"warehouseControl/internal/auth"

	"github.com/wb-go/wbf/ginext"
)

// AuthHandler - обработчик запросов авторизации
type AuthHandler struct {
	authService *auth.Service
}

// NewAuthHandler - создание обработчика авторизации
func NewAuthHandler(authService *auth.Service) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

// LoginRequest - структура запроса логина
type LoginRequest struct {
	Username string `json:"username"`
	Role     string `json:"role"`
}

// Login - обработчик входа пользователя и выдачи JWT токена
func (h *AuthHandler) Login(c *ginext.Context) {

	var req LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ginext.H{"error": "invalid request"})
		return
	}

	token, err := h.authService.GenerateToken(req.Username, req.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ginext.H{"error": "could not generate token"})
		return
	}

	c.JSON(http.StatusOK, ginext.H{
		"token": token,
	})
}
