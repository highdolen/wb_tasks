package export

import (
	"context"
	"log"
	"warehouseControl/internal/audit"

	"github.com/wb-go/wbf/dbpg"
)

// Repository - репозиторий для получения истории изменений из базы
type Repository struct {
	db *dbpg.DB
}

// NewRepository - создание репозитория экспорта
func NewRepository(db *dbpg.DB) *Repository {
	return &Repository{db: db}
}

// GetAllHistory - получение всей истории изменений товаров
func (r *Repository) GetAllHistory(ctx context.Context) ([]audit.History, error) {

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
	ORDER BY changed_at DESC
	`

	rows, err := r.db.Master.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			log.Printf("rows close error: %v", err)
		}
	}()

	var history []audit.History

	for rows.Next() {
		var h audit.History

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
