package repository

import (
	"context"
	"fmt"
	"salesTracker/internal/models"
	"time"
)

// GetGroupedAnalytics - возвращает сгруппированную аналитику
func (r *ItemRepository) GetGroupedAnalytics(
	ctx context.Context,
	from time.Time,
	to time.Time,
	groupBy string,
	sort string,
) ([]models.GroupedAnalytics, error) {

	var groupExpr string

	switch groupBy {
	case "day":
		groupExpr = "DATE(created_at)"
	case "week":
		groupExpr = "DATE_TRUNC('week', created_at)"
	case "category":
		groupExpr = "category"
	default:
		groupExpr = "DATE(created_at)"
	}

	if sort != "asc" {
		sort = "desc"
	}

	query := fmt.Sprintf(`
		SELECT
			%s as group_key,
			COALESCE(SUM(CASE WHEN type = 'income' THEN amount ELSE 0 END),0) as income_sum,
			COALESCE(SUM(CASE WHEN type = 'expense' THEN amount ELSE 0 END),0) as expense_sum,
			COALESCE(AVG(amount),0),
			COUNT(*)
		FROM items
		WHERE created_at BETWEEN $1 AND $2
		GROUP BY group_key
		ORDER BY group_key %s
	`, groupExpr, sort)

	rows, err := r.db.QueryContext(ctx, query, from, to)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []models.GroupedAnalytics

	for rows.Next() {
		var g models.GroupedAnalytics

		err := rows.Scan(
			&g.GroupKey,
			&g.IncomeSum,
			&g.ExpenseSum,
			&g.Avg,
			&g.Count,
		)
		if err != nil {
			return nil, err
		}

		g.Profit = g.IncomeSum - g.ExpenseSum
		result = append(result, g)
	}

	return result, nil
}
