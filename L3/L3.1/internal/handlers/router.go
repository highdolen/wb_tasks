package handlers

import "github.com/wb-go/wbf/ginext"

func RegisterRoutes(r *ginext.Engine, h *NotificationHandlers) {
	api := r.Group("/notify")

	api.POST("", h.CreateNotification)
	api.GET("/:id", h.GetNotification)
	api.DELETE("/:id", h.DeleteNotification)
}
