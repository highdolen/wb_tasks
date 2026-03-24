package service

import (
	"context"
	"errors"
	"log"
)

// orderService реализует интерфейс OrderService
type orderService struct {
	repo  OrderRepository
	cache CacheService
}

var ErrOrderNotFound = errors.New("заказ не найден")

// NewOrderService создает новый экземпляр сервиса заказов
func NewOrderService(repo OrderRepository, cache CacheService) OrderService {
	return &orderService{
		repo:  repo,
		cache: cache,
	}
}

// LoadFromDB загружает все заказы из БД в кэш
func (s *orderService) LoadFromDB(ctx context.Context) error {
	orders, err := s.repo.GetAllOrders(ctx)
	if err != nil {
		return err
	}

	for _, order := range orders {
		s.cache.Set(order.OrderUID, order)
	}

	log.Printf("Загружено %d заказов в кэш", len(orders))
	return nil
}

// Refresh обновляет конкретный заказ в кэше из базы данных
func (s *orderService) Refresh(ctx context.Context, uid string) error {
	order, err := s.repo.GetOrderByUID(ctx, uid)
	if err != nil {
		return err
	}

	if order != nil {
		s.cache.Set(uid, *order)
	} else {
		s.cache.Delete(uid)
	}

	return nil
}

// GetOrderByUID получает заказ по UID с использованием кеша
func (s *orderService) GetOrderByUID(ctx context.Context, uid string) (*OrderResult, error) {
	if cachedOrder, found := s.cache.Get(uid); found {
		return &OrderResult{
			Order:     &cachedOrder,
			FromCache: true,
		}, nil
	}

	order, err := s.repo.GetOrderByUID(ctx, uid)
	if err != nil {
		return nil, err
	}

	if order == nil {
		return nil, ErrOrderNotFound
	}

	s.cache.Set(uid, *order)

	return &OrderResult{
		Order:     order,
		FromCache: false,
	}, nil
}

// GetOrderByUIDWithRefresh принудительно обновляет заказ из БД и возвращает его
func (s *orderService) GetOrderByUIDWithRefresh(ctx context.Context, uid string) (*OrderResult, error) {
	order, err := s.repo.GetOrderByUID(ctx, uid)
	if err != nil {
		return nil, err
	}

	if order == nil {
		s.cache.Delete(uid)
		return nil, ErrOrderNotFound
	}

	s.cache.Set(uid, *order)

	return &OrderResult{
		Order:     order,
		FromCache: false,
	}, nil
}

// GetCacheStats возвращает статистику кеша
func (s *orderService) GetCacheStats() map[string]interface{} {
	return s.cache.GetStats()
}

// InvalidateCache инвалидирует конкретный заказ в кеше
func (s *orderService) InvalidateCache(uid string) error {
	s.cache.Delete(uid)
	return nil
}

// InvalidateAllCache полностью очищает кеш
func (s *orderService) InvalidateAllCache() error {
	s.cache.InvalidateAll()
	return nil
}
