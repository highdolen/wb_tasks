package handlers

import "github.com/wb-go/wbf/ginext"

// SetupRouter настраивает все HTTP-маршруты API
func SetupRouter(
	router *ginext.Engine,
	itemHandler *ItemHandler,
	analyticsHandler *AnalyticsHandler,
	exportHandler *ExportHandler,
) {
	// Общая группа API
	api := router.Group("/api")

	// Работа с финансовыми операциями
	items := api.Group("/items")
	{
		// Создать новую запись
		items.POST("", itemHandler.Create)

		// Получить список записей
		items.GET("", itemHandler.GetAll)

		// Обновить запись по ID
		items.PUT("/:id", itemHandler.Update)

		// Удалить запись по ID
		items.DELETE("/:id", itemHandler.Delete)
	}

	// Метрики и агрегированные данные
	analytics := api.Group("/analytics")
	{
		// Общая аналитика (суммы, медиана и т.д.)
		analytics.GET("", analyticsHandler.GetAnalytics)

		// Группированная аналитика (например по категориям или датам)
		analytics.GET("/grouped", analyticsHandler.GetGroupedAnalytics)

		// Экспорт аналитики в CSV
		analytics.GET("/export/csv", exportHandler.ExportCSV)
	}
}
