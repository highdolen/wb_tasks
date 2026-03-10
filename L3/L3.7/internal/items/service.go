package items

import "context"

// Service - бизнес логика работы с товарами
type Service struct {
	repo *Repository
}

// NewService - создание сервиса товаров
func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

// CreateItem - создание товара
func (s *Service) CreateItem(ctx context.Context, item *Item, username string) error {
	return s.repo.CreateItem(ctx, username, item)
}

// GetItems - получение списка товаров
func (s *Service) GetItems(ctx context.Context) ([]Item, error) {
	return s.repo.GetItems(ctx)
}

// UpdateItem - обновление товара
func (s *Service) UpdateItem(ctx context.Context, id int, item *Item, username string) error {
	return s.repo.UpdateItem(ctx, id, item, username)
}

// DeleteItem - удаление товара
func (s *Service) DeleteItem(ctx context.Context, id int, username string) error {
	return s.repo.DeleteItem(ctx, id, username)
}
