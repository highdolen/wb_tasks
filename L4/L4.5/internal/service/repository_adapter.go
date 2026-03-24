package service

import (
	"context"

	"L0_optimize/internal/database"
	"L0_optimize/internal/models"
)

// repositoryAdapter адаптирует существующий репозиторий к интерфейсу OrderRepository
type repositoryAdapter struct {
	repo *database.OrderRepository
}

// NewRepositoryAdapter создает новый адаптер для репозитория
func NewRepositoryAdapter(repo *database.OrderRepository) OrderRepository {
	return &repositoryAdapter{
		repo: repo,
	}
}

// GetOrderByUID получает заказ из базы данных по UID
func (a *repositoryAdapter) GetOrderByUID(ctx context.Context, uid string) (*models.Order, error) {
	return a.repo.GetOrderByUID(ctx, uid)
}

// CreateOrder создает новый заказ в базе данных
func (a *repositoryAdapter) CreateOrder(ctx context.Context, order *models.Order) error {
	return a.repo.CreateOrder(ctx, order)
}

// GetAllOrders получает все заказы из базы данных
func (a *repositoryAdapter) GetAllOrders(ctx context.Context) ([]models.Order, error) {
	return a.repo.GetAllOrders(ctx)
}
