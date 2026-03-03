package models

// GroupedAnalytics - группировка по ключу
type GroupedAnalytics struct {
	GroupKey   string  `json:"group_key"`   // Ключ группировки
	IncomeSum  float64 `json:"income_sum"`  // Сумма доходов
	ExpenseSum float64 `json:"expense_sum"` // Сумма расходов
	Profit     float64 `json:"profit"`      // Прибыль (Income - Expense)
	Avg        float64 `json:"avg"`         // Среднее значение
	Count      int64   `json:"count"`       // Количество операций
}
