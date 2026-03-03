package models

import "time"

// Item — основная сущность финансовой операции
type Item struct {
	ID        int64     `json:"id"`         // Уникальный ID записи
	Type      string    `json:"type"`       // Тип операции
	Amount    float64   `json:"amount"`     // Сумма
	Category  string    `json:"category"`   // Категория
	CreatedAt time.Time `json:"created_at"` // Дата создания
}

// ItemFilter используется для фильтрации записей при запросах к API
type ItemFilter struct {
	Type     string     // Фильтр по типу операции
	Category string     // Фильтр по категории
	From     *time.Time // Начальная дата диапазона
	To       *time.Time // Конечная дата диапазона
}
