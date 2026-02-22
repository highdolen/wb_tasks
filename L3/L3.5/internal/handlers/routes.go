package handlers

import "github.com/wb-go/wbf/ginext"

// RegisterRoutes - регистрация маршрутов
func RegisterRoutes(engine *ginext.Engine, eventHandler *EventHandler, bookingHandler *BookingHandler, userHandler *UserHandler) {
	// Группа для событий
	api := engine.Group("/events")
	api.POST("", eventHandler.CreateEvent)     // создать событие
	api.GET("", eventHandler.GetAllEvents)     // список всех событий
	api.GET("/:id", eventHandler.GetEvent)     // детали события
	api.POST("/:id/book", bookingHandler.Book) // забронировать место

	// Подтверждение брони через отдельный маршрут
	engine.POST("/bookings/:id/confirm", bookingHandler.Confirm)

	// Пользователи
	users := engine.Group("/users")
	users.POST("", userHandler.Register) // регистрация или получение пользователя

	// Маршрут для получения всех броней текущего пользователя
	users.GET("/:id/bookings", userHandler.GetBookings)
}
