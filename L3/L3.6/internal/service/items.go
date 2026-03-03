package service

import (
	"context"
	"errors"
	"time"

	"salesTracker/internal/models"
	"salesTracker/internal/repository"
)

// ItemService - сервис бизнес-логики для работы с операциями
type ItemService struct {
	repo *repository.ItemRepository
}

// NewItemService - создает новый сервис операций
func NewItemService(r *repository.ItemRepository) *ItemService {
	return &ItemService{repo: r}
}

// Create - создает новую операцию (доход или расход)
func (s *ItemService) Create(ctx context.Context, item *models.Item) (int64, error) {
	if item.Amount < 0 {
		return 0, errors.New("amount cannot be negative")
	}

	if item.CreatedAt.IsZero() {
		return 0, errors.New("created_at is required")
	}

	if item.CreatedAt.After(time.Now()) {
		return 0, errors.New("transaction date cannot be in the future")
	}

	return s.repo.Create(ctx, item)
}

// GetAll - возвращает список операций с учетом фильтров
func (s *ItemService) GetAll(ctx context.Context, filter models.ItemFilter) ([]models.Item, error) {
	return s.repo.GetAll(ctx, filter)
}

// Delete - удаляет операцию по ID
func (s *ItemService) Delete(ctx context.Context, id int64) error {
	return s.repo.Delete(ctx, id)
}

// Update - обновляет существующую операцию
func (s *ItemService) Update(ctx context.Context, id int64, item *models.Item) error {
	if item.Amount < 0 {
		return errors.New("amount cannot be negative")
	}

	if !item.CreatedAt.IsZero() && item.CreatedAt.After(time.Now()) {
		return errors.New("transaction date cannot be in the future")
	}

	return s.repo.Update(ctx, id, item)
}
