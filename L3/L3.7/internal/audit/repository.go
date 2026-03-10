package audit

import (
	"context"
	"log"

	"github.com/wb-go/wbf/dbpg"
)

// Repository - репозиторий для работы с таблицей истории изменений
type Repository struct {
	db *dbpg.DB
}

// NewRepository - создание нового репозитория истории
func NewRepository(db *dbpg.DB) *Repository {
	return &Repository{
		db: db,
	}
}

// GetHistory - получение истории изменений конкретного товара
func (r *Repository) GetHistory(ctx context.Context, itemID int) ([]History, error) {

	query := `
	SELECT 
		id,
		item_id,
		action,
		COALESCE(old_value::text, ''),
		COALESCE(new_value::text, ''),
		COALESCE(changed_by, ''),
		changed_at
	FROM item_history
	WHERE item_id = $1
	ORDER BY changed_at DESC
	`

	rows, err := r.db.Master.QueryContext(ctx, query, itemID)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			log.Printf("rows close error: %v", err)
		}
	}()

	var history []History

	for rows.Next() {
		var h History

		err := rows.Scan(
			&h.ID,
			&h.ItemID,
			&h.Action,
			&h.OldValue,
			&h.NewValue,
			&h.ChangedBy,
			&h.ChangedAt,
		)
		if err != nil {
			return nil, err
		}

		history = append(history, h)
	}

	return history, nil
}

// FilterHistory - фильтрация истории изменений по пользователю, действию и дате
func (r *Repository) FilterHistory(
	ctx context.Context,
	user, action, from, to string,
) ([]History, error) {

	query := `
	SELECT 
		id,
		item_id,
		action,
		COALESCE(old_value::text,''),
		COALESCE(new_value::text,''),
		COALESCE(changed_by,''),
		changed_at
	FROM item_history
	WHERE ($1 = '' OR changed_by = $1)
	AND ($2 = '' OR action = $2)
	AND ($3 = '' OR changed_at::date >= $3::date)
	AND ($4 = '' OR changed_at::date <= $4::date)
	ORDER BY changed_at DESC
	`

	rows, err := r.db.Master.QueryContext(
		ctx,
		query,
		user,
		action,
		from,
		to,
	)
	if err != nil {
		return nil, err
	}

	defer func() {
		if err := rows.Close(); err != nil {
			log.Printf("rows close error: %v", err)
		}
	}()

	var history []History

	for rows.Next() {

		var h History

		err := rows.Scan(
			&h.ID,
			&h.ItemID,
			&h.Action,
			&h.OldValue,
			&h.NewValue,
			&h.ChangedBy,
			&h.ChangedAt,
		)

		if err != nil {
			return nil, err
		}

		history = append(history, h)
	}

	return history, nil
}
