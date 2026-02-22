package user

import (
	"context"
	"database/sql"

	"github.com/wb-go/wbf/dbpg"
)

type UserRepository struct {
	db *dbpg.DB
}

// NewRepository - конструктор UserRepository
func NewRepository(db *dbpg.DB) *UserRepository {
	return &UserRepository{db: db}
}

// CREATE
func (r *UserRepository) Create(
	ctx context.Context,
	name,
	email,
	telegramID,
	role string,
) (*User, error) {

	query := `
		INSERT INTO users (name, email, telegram_id, role)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at
	`

	u := &User{
		Name:       name,
		Email:      email,
		TelegramID: telegramID,
		Role:       role,
	}

	err := r.db.Master.QueryRowContext(
		ctx,
		query,
		name,
		email,
		telegramID,
		role,
	).Scan(&u.ID, &u.CreatedAt)

	if err != nil {
		return nil, err
	}

	return u, nil
}

// GET BY ID
func (r *UserRepository) GetByID(ctx context.Context, id int64) (*User, error) {
	query := `
		SELECT id, name, email, telegram_id, role, created_at
		FROM users
		WHERE id = $1
	`

	var u User

	err := r.db.Master.QueryRowContext(
		ctx,
		query,
		id,
	).Scan(
		&u.ID,
		&u.Name,
		&u.Email,
		&u.TelegramID,
		&u.Role,
		&u.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &u, nil
}

// GET BY EMAIL
func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*User, error) {
	query := `
		SELECT id, name, email, telegram_id, role, created_at
		FROM users
		WHERE email = $1
	`

	var u User

	err := r.db.Master.QueryRowContext(
		ctx,
		query,
		email,
	).Scan(
		&u.ID,
		&u.Name,
		&u.Email,
		&u.TelegramID,
		&u.Role,
		&u.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return &u, nil
}

// UPDATE TELEGRAM ID
func (r *UserRepository) UpdateTelegramID(
	ctx context.Context,
	userID int64,
	telegramID string,
) error {

	query := `
		UPDATE users
		SET telegram_id = $1
		WHERE id = $2
	`

	_, err := r.db.Master.ExecContext(
		ctx,
		query,
		telegramID,
		userID,
	)

	return err
}

// UPDATE ROLE
func (r *UserRepository) UpdateRole(
	ctx context.Context,
	userID int64,
	role string,
) error {

	query := `
		UPDATE users
		SET role = $1
		WHERE id = $2
	`

	_, err := r.db.Master.ExecContext(
		ctx,
		query,
		role,
		userID,
	)

	return err
}
