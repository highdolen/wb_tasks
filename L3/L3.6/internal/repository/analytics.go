package repository

import (
	"context"
	"salesTracker/internal/models"
	"time"
)

// GetAnalytics - возвращает общую аналитику по доходам и расходам
// за указанный период времени
func (r *ItemRepository) GetAnalytics(ctx context.Context, from time.Time, to time.Time) (*models.Analytics, error) {
	query := `
	SELECT
		COALESCE(SUM(amount) FILTER (WHERE type='income'),0),
		COALESCE(AVG(amount) FILTER (WHERE type='income'),0),
		COUNT(*) FILTER (WHERE type='income'),
		COALESCE(percentile_cont(0.5) WITHIN GROUP (ORDER BY amount) FILTER (WHERE type='income'),0),
		COALESCE(percentile_cont(0.9) WITHIN GROUP (ORDER BY amount) FILTER (WHERE type='income'),0),

		COALESCE(SUM(amount) FILTER (WHERE type='expense'),0),
		COALESCE(AVG(amount) FILTER (WHERE type='expense'),0),
		COUNT(*) FILTER (WHERE type='expense'),
		COALESCE(percentile_cont(0.5) WITHIN GROUP (ORDER BY amount) FILTER (WHERE type='expense'),0),
		COALESCE(percentile_cont(0.9) WITHIN GROUP (ORDER BY amount) FILTER (WHERE type='expense'),0)
	FROM items
	WHERE created_at BETWEEN $1 AND $2
	`

	row := r.db.QueryRowContext(ctx, query, from, to)

	var a models.Analytics
	err := row.Scan(
		&a.Income.Sum,
		&a.Income.Avg,
		&a.Income.Count,
		&a.Income.Median,
		&a.Income.Percentile90,
		&a.Expense.Sum,
		&a.Expense.Avg,
		&a.Expense.Count,
		&a.Expense.Median,
		&a.Expense.Percentile90,
	)
	if err != nil {
		return nil, err
	}

	a.Balance = a.Income.Sum - a.Expense.Sum
	return &a, nil
}
