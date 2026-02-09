package handlers

import "github.com/wb-go/wbf/ginext"

//RegisterRoutes - регистрация маршрутов
func RegisterRoutes(r *ginext.Engine, h *Handler) {
	r.POST("/upload", h.Upload)
	r.GET("/image/:id", h.GetImage)
	r.DELETE("/image/:id", h.DeleteImage)
}
