package models

// SummaryData содержит агрегированные метрики
type SummaryData struct {
	Sum          float64 `json:"sum"`           // Общая сумма
	Avg          float64 `json:"avg"`           // Среднее значение
	Count        int64   `json:"count"`         // Количество операций
	Median       float64 `json:"median"`        // Медиана
	Percentile90 float64 `json:"percentile_90"` // 90-й перцентиль
}

// Analytics — общий объект аналитики
type Analytics struct {
	Income  SummaryData `json:"income"`  // Аналитика по доходам
	Expense SummaryData `json:"expense"` // Аналитика по расходам
	Balance float64     `json:"balance"` // Баланс (Income - Expense)
}
