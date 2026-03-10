package users

import (
	"context"
	"database/sql"
)

// Repository - репозиторий для работы с пользователями
type Repository struct {
	db *sql.DB
}

// NewRepository - создание репозитория пользователей
func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		db: db,
	}
}

// GetByUsername - получение пользователя по имени
func (r *Repository) GetByUsername(ctx context.Context, username string) (*User, error) {

	query := `
	SELECT id, username, role
	FROM users
	WHERE username = $1
	`

	var user User

	err := r.db.QueryRowContext(ctx, query, username).Scan(
		&user.ID,
		&user.Username,
		&user.Role,
	)

	if err != nil {
		return nil, err
	}

	return &user, nil
}
