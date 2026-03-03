package repository

import (
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"time"
)

// ExportCSV - экспортирует аналитику и сгруппированные данные в CSV
func (r *ItemRepository) ExportCSV(
	ctx context.Context,
	from time.Time,
	to time.Time,
	w io.Writer,
) error {

	// Получаем агрегированную аналитику
	analytics, err := r.GetAnalytics(ctx, from, to)
	if err != nil {
		return err
	}

	// Получаем сгруппированную аналитику по категориям
	grouped, err := r.GetGroupedAnalytics(ctx, from, to, "category", "desc")
	if err != nil {
		return err
	}

	writer := csv.NewWriter(w)

	// Заголовок отчёта
	if err := writer.Write([]string{"REPORT"}); err != nil {
		return err
	}

	// Income статистика
	if err := writer.Write([]string{"income_sum", fmt.Sprint(analytics.Income.Sum)}); err != nil {
		return err
	}
	if err := writer.Write([]string{"income_avg", fmt.Sprint(analytics.Income.Avg)}); err != nil {
		return err
	}
	if err := writer.Write([]string{"income_count", fmt.Sprint(analytics.Income.Count)}); err != nil {
		return err
	}
	if err := writer.Write([]string{"income_median", fmt.Sprint(analytics.Income.Median)}); err != nil {
		return err
	}
	if err := writer.Write([]string{"income_p90", fmt.Sprint(analytics.Income.Percentile90)}); err != nil {
		return err
	}

	// Пустая строка-разделитель
	if err := writer.Write([]string{}); err != nil {
		return err
	}

	// Expense статистика
	if err := writer.Write([]string{"expense_sum", fmt.Sprint(analytics.Expense.Sum)}); err != nil {
		return err
	}
	if err := writer.Write([]string{"expense_avg", fmt.Sprint(analytics.Expense.Avg)}); err != nil {
		return err
	}
	if err := writer.Write([]string{"expense_count", fmt.Sprint(analytics.Expense.Count)}); err != nil {
		return err
	}
	if err := writer.Write([]string{"expense_median", fmt.Sprint(analytics.Expense.Median)}); err != nil {
		return err
	}
	if err := writer.Write([]string{"expense_p90", fmt.Sprint(analytics.Expense.Percentile90)}); err != nil {
		return err
	}

	// Пустая строка-разделитель
	if err := writer.Write([]string{}); err != nil {
		return err
	}

	// Баланс
	if err := writer.Write([]string{"balance", fmt.Sprint(analytics.Balance)}); err != nil {
		return err
	}

	if err := writer.Write([]string{}); err != nil {
		return err
	}

	// Заголовок блока группировки
	if err := writer.Write([]string{"GROUPED BY CATEGORY"}); err != nil {
		return err
	}

	if err := writer.Write([]string{
		"group_key",
		"income_sum",
		"expense_sum",
		"profit",
		"avg",
		"count",
	}); err != nil {
		return err
	}

	// Записываем строки группированной аналитики
	for _, g := range grouped {
		if err := writer.Write([]string{
			g.GroupKey,
			fmt.Sprint(g.IncomeSum),
			fmt.Sprint(g.ExpenseSum),
			fmt.Sprint(g.Profit),
			fmt.Sprint(g.Avg),
			fmt.Sprint(g.Count),
		}); err != nil {
			return err
		}
	}

	// Принудительно записываем буфер в writer
	writer.Flush()

	// Проверяем ошибку записи CSV
	if err := writer.Error(); err != nil {
		return err
	}

	return nil
}
