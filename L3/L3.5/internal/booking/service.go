package booking

import (
	"context"
	"fmt"

	"eventBooker/internal/notification"
	"eventBooker/internal/user"
)

// Интерфейс репозитория
type Repository interface {
	Create(ctx context.Context, eventID, userID int64, ttlSeconds int64) (int64, error)
	Confirm(ctx context.Context, bookingID int64) error
	CancelExpired(ctx context.Context) ([]user.User, error)
	GetUserByBookingID(ctx context.Context, bookingID int64) (*user.User, error)
	GetByUserID(ctx context.Context, userID int64) ([]Booking, error)
}

type Service struct {
	repo Repository
	tg   notification.Sender
}

// NewService - конструктор сервиса
func NewService(r Repository, tg notification.Sender) *Service {
	return &Service{
		repo: r,
		tg:   tg,
	}
}

// BookSeatWithID - бронирование места
func (s *Service) BookSeatWithID(ctx context.Context, eventID, userID int64, ttlSeconds int64) (int64, error) {
	return s.repo.Create(ctx, eventID, userID, ttlSeconds)
}

// ConfirmBooking - подтверждение оплаты брони
func (s *Service) ConfirmBooking(ctx context.Context, bookingID int64) error {
	if err := s.repo.Confirm(ctx, bookingID); err != nil {
		return err
	}

	u, err := s.repo.GetUserByBookingID(ctx, bookingID)
	if err != nil {
		return err
	}

	if u != nil && u.TelegramID != "" {
		if err := s.tg.Send(ctx, u.TelegramID, "✅ Ваша бронь подтверждена!"); err != nil {
			fmt.Println("telegram send error:", err)
		}
	}

	return nil
}

// CancelExpiredBookings - отмена созданной брони
func (s *Service) CancelExpiredBookings(ctx context.Context) error {
	users, err := s.repo.CancelExpired(ctx)
	if err != nil {
		return err
	}

	for _, u := range users {
		if u.TelegramID == "" {
			continue
		}
		if err := s.tg.Send(ctx, u.TelegramID, "❌ Ваша бронь истекла."); err != nil {
			fmt.Println("telegram send error:", err)
		}
	}

	return nil
}

// GetUserBookings - получение всех броней конкретного пользователя
func (s *Service) GetUserBookings(ctx context.Context, userID int64) ([]Booking, error) {
	return s.repo.GetByUserID(ctx, userID)
}
