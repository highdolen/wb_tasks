package booking

import (
	"context"
	"errors"

	"eventBooker/internal/user"

	"github.com/wb-go/wbf/dbpg"
)

type BookingRepository struct {
	db *dbpg.DB
}

func NewRepository(db *dbpg.DB) *BookingRepository {
	return &BookingRepository{db: db}
}

// Create - создание брони и возврат ID
func (r *BookingRepository) Create(ctx context.Context, eventID, userID int64, ttlSeconds int64) (int64, error) {
	tx, err := r.db.Master.BeginTx(ctx, nil)
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	res, err := tx.ExecContext(ctx, `
		UPDATE events
		SET available_seats = available_seats - 1
		WHERE id = $1 AND available_seats > 0
	`, eventID)
	if err != nil {
		return 0, err
	}

	rows, _ := res.RowsAffected()
	if rows == 0 {
		return 0, errors.New("no available seats")
	}

	var bookingID int64
	err = tx.QueryRowContext(ctx, `
		INSERT INTO bookings (event_id, user_id, status, created_at, expires_at)
		VALUES ($1, $2, 'pending', NOW(), NOW() + $3 * interval '1 second')
		RETURNING id
	`, eventID, userID, ttlSeconds).Scan(&bookingID)
	if err != nil {
		return 0, err
	}

	return bookingID, tx.Commit()
}

// Confirm - подтверждение оплаты брони
func (r *BookingRepository) Confirm(ctx context.Context, bookingID int64) error {
	res, err := r.db.Master.ExecContext(ctx, `
		UPDATE bookings
		SET status = 'confirmed'
		WHERE id = $1
		  AND status = 'pending'
		  AND expires_at > NOW()
	`, bookingID)

	if err != nil {
		return err
	}

	rows, _ := res.RowsAffected()
	if rows == 0 {
		return errors.New("cannot confirm booking")
	}

	return nil
}

func (r *BookingRepository) CancelExpired(ctx context.Context) ([]user.User, error) {
	tx, err := r.db.Master.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	rows, err := tx.QueryContext(ctx, `
		UPDATE bookings b
		SET status = 'canceled'
		FROM users u
		WHERE b.user_id = u.id
		  AND b.status = 'pending'
		  AND b.expires_at <= NOW()
		RETURNING b.event_id, u.id, u.name, u.email, u.telegram_id
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []user.User
	var eventIDs []int64

	for rows.Next() {
		var u user.User
		var eventID int64
		if err := rows.Scan(&eventID, &u.ID, &u.Name, &u.Email, &u.TelegramID); err != nil {
			return nil, err
		}
		users = append(users, u)
		eventIDs = append(eventIDs, eventID)
	}

	for _, eventID := range eventIDs {
		_, err := tx.ExecContext(ctx, `
			UPDATE events
			SET available_seats = available_seats + 1
			WHERE id = $1
		`, eventID)
		if err != nil {
			return nil, err
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return users, nil
}

// GetUserByBookingID - получение пользователя по айди брони
func (r *BookingRepository) GetUserByBookingID(ctx context.Context, bookingID int64) (*user.User, error) {
	query := `
		SELECT u.id, u.name, u.email, u.telegram_id
		FROM bookings b
		JOIN users u ON u.id = b.user_id
		WHERE b.id = $1
	`

	var u user.User
	err := r.db.Master.QueryRowContext(ctx, query, bookingID).Scan(
		&u.ID, &u.Name, &u.Email, &u.TelegramID,
	)
	if err != nil {
		return nil, err
	}

	return &u, nil
}

// GetByUserID - получение брони по юсер айди
func (r *BookingRepository) GetByUserID(ctx context.Context, userID int64) ([]Booking, error) {
	rows, err := r.db.Master.QueryContext(ctx, `
		SELECT id, event_id, status, created_at, expires_at
		FROM bookings
		WHERE user_id = $1
		ORDER BY created_at DESC
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var bookings []Booking
	for rows.Next() {
		var b Booking
		var status string
		if err := rows.Scan(&b.ID, &b.EventID, &status, &b.CreatedAt, &b.ExpiresAt); err != nil {
			return nil, err
		}
		b.Status = Status(status)
		bookings = append(bookings, b)
	}

	return bookings, nil
}
