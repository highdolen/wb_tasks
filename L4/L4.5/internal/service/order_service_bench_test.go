package service

import (
	"context"
	"testing"
	"time"

	"L0_optimize/internal/cache"
	"L0_optimize/internal/models"
)

type benchRepo struct {
	order *models.Order
}

func (r *benchRepo) GetOrderByUID(ctx context.Context, uid string) (*models.Order, error) {
	if r.order == nil || r.order.OrderUID != uid {
		return nil, nil
	}
	return r.order, nil
}

func (r *benchRepo) CreateOrder(ctx context.Context, order *models.Order) error {
	return nil
}

func (r *benchRepo) GetAllOrders(ctx context.Context) ([]models.Order, error) {
	if r.order == nil {
		return nil, nil
	}
	return []models.Order{*r.order}, nil
}

func benchmarkOrder() models.Order {
	return models.Order{
		OrderUID:    "bench-order-uid",
		TrackNumber: "WBILMTESTTRACK",
		Entry:       "WBIL",
		Delivery: models.Delivery{
			ID:      1,
			Name:    "Test User",
			Phone:   "+79990000000",
			Zip:     "123456",
			City:    "Moscow",
			Address: "Lenina 1",
			Region:  "Moscow",
			Email:   "test@example.com",
		},
		Payment: models.Payment{
			ID:           1,
			Transaction:  "bench-order-uid",
			RequestID:    "",
			Currency:     "USD",
			Provider:     "wbpay",
			Amount:       1000,
			PaymentDt:    1637907727,
			Bank:         "alpha",
			DeliveryCost: 200,
			GoodsTotal:   800,
			CustomFee:    0,
		},
		Items: []models.Item{
			{
				ID:          1,
				ChrtID:      123,
				TrackNumber: "WBILMTESTTRACK",
				Price:       800,
				Rid:         "rid-1",
				Name:        "Test Item",
				Sale:        0,
				Size:        "M",
				TotalPrice:  800,
				NmID:        111,
				Brand:       "Brand",
				Status:      202,
				OrderUID:    "bench-order-uid",
			},
		},
		Locale:            "en",
		InternalSignature: "",
		CustomerID:        "customer-1",
		DeliveryService:   "meest",
		Shardkey:          "9",
		SmID:              99,
		DateCreated:       time.Unix(1637907727, 0),
		OofShard:          "1",
	}
}

func BenchmarkOrderService_GetOrderByUID_CacheHit(b *testing.B) {
	order := benchmarkOrder()
	repo := &benchRepo{order: &order}
	c := cache.New(30 * time.Minute)
	defer c.Close()

	c.Set(order.OrderUID, order)

	svc := NewOrderService(repo, NewCacheAdapter(c))
	ctx := context.Background()

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		result, err := svc.GetOrderByUID(ctx, order.OrderUID)
		if err != nil {
			b.Fatalf("unexpected error: %v", err)
		}
		if result == nil || result.Order == nil {
			b.Fatal("expected order, got nil")
		}
	}
}

func BenchmarkOrderService_GetOrderByUID_CacheMiss(b *testing.B) {
	order := benchmarkOrder()
	repo := &benchRepo{order: &order}
	c := cache.New(30 * time.Minute)
	defer c.Close()

	svc := NewOrderService(repo, NewCacheAdapter(c))
	ctx := context.Background()

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		c.InvalidateAll()

		result, err := svc.GetOrderByUID(ctx, order.OrderUID)
		if err != nil {
			b.Fatalf("unexpected error: %v", err)
		}
		if result == nil || result.Order == nil {
			b.Fatal("expected order, got nil")
		}
	}
}

func BenchmarkOrderService_GetOrderByUIDWithRefresh(b *testing.B) {
	order := benchmarkOrder()
	repo := &benchRepo{order: &order}
	c := cache.New(30 * time.Minute)
	defer c.Close()

	svc := NewOrderService(repo, NewCacheAdapter(c))
	ctx := context.Background()

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		result, err := svc.GetOrderByUIDWithRefresh(ctx, order.OrderUID)
		if err != nil {
			b.Fatalf("unexpected error: %v", err)
		}
		if result == nil || result.Order == nil {
			b.Fatal("expected order, got nil")
		}
	}
}
