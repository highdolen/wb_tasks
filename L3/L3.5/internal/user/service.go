package user

import (
	"context"
)

// Repository - интерфейс
type Repository interface {
	Create(ctx context.Context, name, email, telegramID, role string) (*User, error)
	GetByID(ctx context.Context, id int64) (*User, error)
	GetByEmail(ctx context.Context, email string) (*User, error)
	UpdateTelegramID(ctx context.Context, userID int64, telegramID string) error
	UpdateRole(ctx context.Context, userID int64, role string) error
}

type Service struct {
	repo Repository
}

// NewService - конструктор Service
func NewService(r Repository) *Service {
	return &Service{repo: r}
}

// GetOrCreate - вход или создание нового пользователя
func (s *Service) GetOrCreate(
	ctx context.Context,
	email,
	name,
	telegramID,
	role string,
) (*User, error) {

	// Проверяем существует ли пользователь
	u, err := s.repo.GetByEmail(ctx, email)
	if err == nil && u != nil {

		// обновляем telegram_id если изменился
		if telegramID != "" && u.TelegramID != telegramID {
			if err := s.repo.UpdateTelegramID(ctx, u.ID, telegramID); err != nil {
				return nil, err
			}
			u.TelegramID = telegramID
		}

		// обновляем роль если передана
		if role != "" && u.Role != role {
			if err := s.repo.UpdateRole(ctx, u.ID, role); err != nil {
				return nil, err
			}
			u.Role = role
		}

		return u, nil
	}

	// Если пользователя нет — создаём нового
	return s.repo.Create(ctx, name, email, telegramID, role)
}
