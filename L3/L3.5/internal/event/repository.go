package event

import (
	"context"
	"time"

	"github.com/wb-go/wbf/dbpg"
)

type PostgresRepository struct {
	db *dbpg.DB
}

// NewRepository - конструктор PostgresRepository
func NewRepository(db *dbpg.DB) *PostgresRepository {
	return &PostgresRepository{db: db}
}

// Create - создание нового события
func (r *PostgresRepository) Create(ctx context.Context, e *Event) error {
	ttlSec := int64(e.BookingTTL.Seconds())
	query := `
		INSERT INTO events (name, date, total_seats, available_seats, booking_ttl)
		VALUES ($1, $2, $3, $3, $4)
		RETURNING id, created_at
	`
	return r.db.QueryRowContext(ctx, query,
		e.Name, e.Date, e.TotalSeats, ttlSec,
	).Scan(&e.ID, &e.CreatedAt)
}

// GetByID - получение событие по айди
func (r *PostgresRepository) GetByID(ctx context.Context, id int64) (*Event, error) {
	var ttlSec int64
	e := &Event{}
	err := r.db.QueryRowContext(ctx, `
		SELECT id, name, date, total_seats, available_seats, booking_ttl, created_at
		FROM events
		WHERE id = $1
	`, id).Scan(
		&e.ID, &e.Name, &e.Date, &e.TotalSeats, &e.AvailableSeats, &ttlSec, &e.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	e.BookingTTL = time.Duration(ttlSec) * time.Second // <-- обратно в Duration
	return e, nil
}

// GetAll - получение всех событий
func (r *PostgresRepository) GetAll(ctx context.Context) ([]*Event, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, name, date, total_seats, available_seats, booking_ttl, created_at
		FROM events
		ORDER BY date ASC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []*Event

	for rows.Next() {
		var ttlSec int64
		e := &Event{}
		err := rows.Scan(
			&e.ID,
			&e.Name,
			&e.Date,
			&e.TotalSeats,
			&e.AvailableSeats,
			&ttlSec,
			&e.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		e.BookingTTL = time.Duration(ttlSec) * time.Second
		events = append(events, e)
	}

	return events, nil
}
