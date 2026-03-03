package repository

import (
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"time"
)

// ExportCSV - экспортирует аналитику и сгруппированные данные в CSV файл
func (r *ItemRepository) ExportCSV(ctx context.Context, from time.Time, to time.Time, w io.Writer) error {

	analytics, err := r.GetAnalytics(ctx, from, to)
	if err != nil {
		return err
	}

	grouped, err := r.GetGroupedAnalytics(ctx, from, to, "category", "desc")
	if err != nil {
		return err
	}

	writer := csv.NewWriter(w)

	writer.Write([]string{"REPORT"})
	writer.Write([]string{"income_sum", fmt.Sprint(analytics.Income.Sum)})
	writer.Write([]string{"income_avg", fmt.Sprint(analytics.Income.Avg)})
	writer.Write([]string{"income_count", fmt.Sprint(analytics.Income.Count)})
	writer.Write([]string{"income_median", fmt.Sprint(analytics.Income.Median)})
	writer.Write([]string{"income_p90", fmt.Sprint(analytics.Income.Percentile90)})

	writer.Write([]string{})

	writer.Write([]string{"expense_sum", fmt.Sprint(analytics.Expense.Sum)})
	writer.Write([]string{"expense_avg", fmt.Sprint(analytics.Expense.Avg)})
	writer.Write([]string{"expense_count", fmt.Sprint(analytics.Expense.Count)})
	writer.Write([]string{"expense_median", fmt.Sprint(analytics.Expense.Median)})
	writer.Write([]string{"expense_p90", fmt.Sprint(analytics.Expense.Percentile90)})

	writer.Write([]string{})

	writer.Write([]string{"balance", fmt.Sprint(analytics.Balance)})
	writer.Write([]string{})

	writer.Write([]string{"GROUPED BY CATEGORY"})
	writer.Write([]string{"group_key", "income_sum", "expense_sum", "profit", "avg", "count"})

	for _, g := range grouped {
		writer.Write([]string{
			g.GroupKey,
			fmt.Sprint(g.IncomeSum),
			fmt.Sprint(g.ExpenseSum),
			fmt.Sprint(g.Profit),
			fmt.Sprint(g.Avg),
			fmt.Sprint(g.Count),
		})
	}

	writer.Flush()
	return writer.Error()
}
