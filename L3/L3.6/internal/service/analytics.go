package service

import (
	"context"
	"time"

	"salesTracker/internal/models"
	"salesTracker/internal/repository"
)

// AnalyticsService - сервис бизнес-логики для аналитики
type AnalyticsService struct {
	repo *repository.ItemRepository
}

// NewAnalyticsService - создает новый сервис аналитики
func NewAnalyticsService(repo *repository.ItemRepository) *AnalyticsService {
	return &AnalyticsService{repo: repo}
}

// GetAnalytics - возвращает общую аналитику по доходам и расходам
func (s *AnalyticsService) GetAnalytics(
	ctx context.Context,
	from time.Time,
	to time.Time,
) (*models.Analytics, error) {
	return s.repo.GetAnalytics(ctx, from, to)
}

// GetGroupedAnalytics - возвращает сгруппированную аналитику
func (s *AnalyticsService) GetGroupedAnalytics(
	ctx context.Context,
	from time.Time,
	to time.Time,
	groupBy string,
	sort string,
) ([]models.GroupedAnalytics, error) {
	return s.repo.GetGroupedAnalytics(ctx, from, to, groupBy, sort)
}
