package items

import (
	"context"
	"database/sql"
	"log"

	"github.com/wb-go/wbf/dbpg"
)

// Repository - репозиторий для работы с таблицей товаров
type Repository struct {
	db *dbpg.DB
}

// NewRepository - создание репозитория товаров
func NewRepository(db *dbpg.DB) *Repository {
	return &Repository{db: db}
}

// CreateItem - создание нового товара
func (r *Repository) CreateItem(ctx context.Context, username string, item *Item) error {

	tx, err := r.db.Master.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		if err := tx.Rollback(); err != nil && err != sql.ErrTxDone {
			log.Printf("tx rollback error: %v", err)
		}
	}()

	// устанавливаем пользователя для триггера
	_, err = tx.ExecContext(ctx, `SELECT set_config('app.user', $1, true)`, username)
	if err != nil {
		return err
	}

	query := `
	INSERT INTO items (name, quantity)
	VALUES ($1, $2)
	RETURNING id, updated_at
	`

	err = tx.QueryRowContext(ctx, query, item.Name, item.Quantity).
		Scan(&item.ID, &item.UpdatedAt)
	if err != nil {
		return err
	}

	return tx.Commit()
}

// GetItems - получение товаров
func (r *Repository) GetItems(ctx context.Context) ([]Item, error) {

	query := `
	SELECT id, name, quantity, updated_at
	FROM items
	ORDER BY id
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			log.Printf("rows close error: %v", err)
		}
	}()

	var items []Item

	for rows.Next() {

		var item Item

		if err := rows.Scan(
			&item.ID,
			&item.Name,
			&item.Quantity,
			&item.UpdatedAt,
		); err != nil {
			return nil, err
		}

		items = append(items, item)
	}

	return items, nil
}

// UpdateItem - обновление товара
func (r *Repository) UpdateItem(ctx context.Context, id int, item *Item, username string) error {

	tx, err := r.db.Master.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		if err := tx.Rollback(); err != nil && err != sql.ErrTxDone {
			log.Printf("tx rollback error: %v", err)
		}
	}()

	_, err = tx.ExecContext(ctx, `SELECT set_config('app.user', $1, true)`, username)
	if err != nil {
		return err
	}

	query := `
	UPDATE items
	SET name=$1, quantity=$2, updated_at=now()
	WHERE id=$3
	`

	_, err = tx.ExecContext(ctx, query, item.Name, item.Quantity, id)
	if err != nil {
		return err
	}

	return tx.Commit()
}

// DeleteItem - удаление товара
func (r *Repository) DeleteItem(ctx context.Context, id int, username string) error {

	tx, err := r.db.Master.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		if err := tx.Rollback(); err != nil && err != sql.ErrTxDone {
			log.Printf("tx rollback error: %v", err)
		}
	}()

	_, err = tx.ExecContext(ctx, `SELECT set_config('app.user', $1, true)`, username)
	if err != nil {
		return err
	}

	query := `DELETE FROM items WHERE id=$1`

	_, err = tx.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	return tx.Commit()
}
