package handlers

import (
	"warehouseControl/internal/auth"

	"github.com/wb-go/wbf/ginext"
)

// Handlers - структура всех обработчиков приложения
type Handlers struct {
	Auth   *AuthHandler
	Items  *ItemsHandler
	Audit  *AuditHandler
	Export *ExportHandler
}

// RegisterRoutes - регистрация всех HTTP маршрутов
func RegisterRoutes(r *ginext.Engine, h *Handlers, jwtSecret string) {

	// маршрут авторизации
	r.POST("/login", h.Auth.Login)

	api := r.Group("/")
	api.Use(auth.AuthMiddleware(jwtSecret))

	{
		// ITEMS

		// создание товара (только admin)
		api.POST("/items",
			auth.RequireRoles("admin"),
			h.Items.CreateItem,
		)

		// получение списка товаров
		api.GET("/items",
			auth.RequireRoles("admin", "manager", "viewer"),
			h.Items.GetItems,
		)

		// обновление товара
		api.PUT("/items/:id",
			auth.RequireRoles("admin", "manager"),
			h.Items.UpdateItem,
		)

		// удаление товара
		api.DELETE("/items/:id",
			auth.RequireRoles("admin"),
			h.Items.DeleteItem,
		)

		// HISTORY

		// история изменений конкретного товара
		api.GET("/items/:id/history",
			auth.RequireRoles("admin", "manager"),
			h.Audit.GetItemHistory,
		)

		// фильтрация истории изменений
		api.GET("/audit/filter",
			auth.RequireRoles("admin", "manager"),
			h.Audit.FilterHistory,
		)

		// экспорт истории в CSV
		api.GET("/audit/export",
			auth.RequireRoles("admin", "manager"),
			h.Export.ExportHistoryCSV,
		)
	}
}
