package users

import "context"

// Service - сервис работы с пользователями
type Service struct {
	repo *Repository
}

// NewService - создание сервиса пользователей
func NewService(repo *Repository) *Service {
	return &Service{
		repo: repo,
	}
}

// GetByUsername - получение пользователя по имени
func (s *Service) GetByUsername(ctx context.Context, username string) (*User, error) {
	return s.repo.GetByUsername(ctx, username)
}
