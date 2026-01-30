package handlers

import (
	"github.com/wb-go/wbf/ginext"
)

// RegisterRoutes — регистрирует все маршруты для приложения
func RegisterRoutes(engine *ginext.Engine, h *Handler) {
	engine.POST("/shorten", h.CreateShortLink)
	engine.GET("/s/:code", h.Redirect)
	engine.GET("/analytics/:code", h.Analytics)
}
