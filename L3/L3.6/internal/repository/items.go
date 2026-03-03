package repository

import (
	"context"
	"fmt"

	"salesTracker/internal/models"

	"github.com/wb-go/wbf/dbpg"
)

// ItemRepository - репозиторий для работы с таблицей items
type ItemRepository struct {
	db *dbpg.DB
}

// NewItemRepository - создает новый экземпляр репозитория
func NewItemRepository(db *dbpg.DB) *ItemRepository {
	return &ItemRepository{db: db}
}

// Create - создает новую запись операции (income / expense) в базе данных
func (r *ItemRepository) Create(ctx context.Context, item *models.Item) (int64, error) {
	query := `
	INSERT INTO items (type, amount, category, created_at)
	VALUES ($1, $2, $3, $4)
	RETURNING id
	`

	var id int64
	err := r.db.QueryRowContext(
		ctx,
		query,
		item.Type,
		item.Amount,
		item.Category,
		item.CreatedAt,
	).Scan(&id)

	return id, err
}

// GetAll - возвращает список операций с возможностью фильтрации
// по типу, категории и диапазону дат
func (r *ItemRepository) GetAll(ctx context.Context, filter models.ItemFilter) ([]models.Item, error) {

	query := `
	SELECT id, type, amount, category, created_at
	FROM items
	WHERE 1=1
	`

	args := []interface{}{}
	i := 1

	if filter.Type != "" {
		query += fmt.Sprintf(" AND type=$%d", i)
		args = append(args, filter.Type)
		i++
	}

	if filter.Category != "" {
		query += fmt.Sprintf(" AND category ILIKE $%d", i)
		args = append(args, "%"+filter.Category+"%")
		i++
	}

	if filter.From != nil {
		query += fmt.Sprintf(" AND created_at >= $%d", i)
		args = append(args, *filter.From)
		i++
	}

	if filter.To != nil {
		query += fmt.Sprintf(" AND created_at <= $%d", i)
		args = append(args, *filter.To)
		i++
	}

	query += " ORDER BY created_at DESC"

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []models.Item

	for rows.Next() {
		var item models.Item

		err := rows.Scan(
			&item.ID,
			&item.Type,
			&item.Amount,
			&item.Category,
			&item.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		items = append(items, item)
	}

	return items, nil
}

// Update - обновляет существующую операцию по ID
func (r *ItemRepository) Update(ctx context.Context, id int64, item *models.Item) error {
	query := `
	UPDATE items
	SET type=$1, amount=$2, category=$3, created_at=$4
	WHERE id=$5
	`

	_, err := r.db.ExecContext(
		ctx,
		query,
		item.Type,
		item.Amount,
		item.Category,
		item.CreatedAt,
		id,
	)

	return err
}

// Delete - удаляет операцию из базы данных по ID
func (r *ItemRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM items WHERE id=$1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}
