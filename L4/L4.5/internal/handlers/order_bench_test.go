package handlers

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"L0_optimize/internal/models"
	"L0_optimize/internal/service"

	"github.com/gorilla/mux"
)

type benchOrderService struct {
	result *service.OrderResult
	stats  map[string]interface{}
	err    error
}

func (s *benchOrderService) LoadFromDB(ctx context.Context) error {
	return nil
}

func (s *benchOrderService) Refresh(ctx context.Context, uid string) error {
	return nil
}

func (s *benchOrderService) GetOrderByUID(ctx context.Context, uid string) (*service.OrderResult, error) {
	return s.result, s.err
}

func (s *benchOrderService) GetOrderByUIDWithRefresh(ctx context.Context, uid string) (*service.OrderResult, error) {
	return s.result, s.err
}

func (s *benchOrderService) GetCacheStats() map[string]interface{} {
	if s.stats != nil {
		return s.stats
	}
	return map[string]interface{}{
		"total_entries": 1,
	}
}

func (s *benchOrderService) InvalidateCache(uid string) error {
	return nil
}

func (s *benchOrderService) InvalidateAllCache() error {
	return nil
}

func benchmarkHandlerOrder() *models.Order {
	return &models.Order{
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

func BenchmarkOrderHandler_GetOrder_CacheHit(b *testing.B) {
	order := benchmarkHandlerOrder()
	svc := &benchOrderService{
		result: &service.OrderResult{
			Order:     order,
			FromCache: true,
		},
	}
	h := NewOrderHandler(svc)

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		req := httptest.NewRequest(http.MethodGet, "/order/"+order.OrderUID, nil)
		req = mux.SetURLVars(req, map[string]string{"order_uid": order.OrderUID})
		rr := httptest.NewRecorder()

		h.GetOrder(rr, req)

		if rr.Code != http.StatusOK {
			b.Fatalf("unexpected status: %d", rr.Code)
		}
	}
}

func BenchmarkOrderHandler_GetOrder_CacheMiss(b *testing.B) {
	order := benchmarkHandlerOrder()
	svc := &benchOrderService{
		result: &service.OrderResult{
			Order:     order,
			FromCache: false,
		},
	}
	h := NewOrderHandler(svc)

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		req := httptest.NewRequest(http.MethodGet, "/order/"+order.OrderUID, nil)
		req = mux.SetURLVars(req, map[string]string{"order_uid": order.OrderUID})
		rr := httptest.NewRecorder()

		h.GetOrder(rr, req)

		if rr.Code != http.StatusOK {
			b.Fatalf("unexpected status: %d", rr.Code)
		}
	}
}

func BenchmarkOrderHandler_GetOrder_WithRefresh(b *testing.B) {
	order := benchmarkHandlerOrder()
	svc := &benchOrderService{
		result: &service.OrderResult{
			Order:     order,
			FromCache: false,
		},
	}
	h := NewOrderHandler(svc)

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		req := httptest.NewRequest(http.MethodGet, "/order/"+order.OrderUID+"?refresh=true", nil)
		req = mux.SetURLVars(req, map[string]string{"order_uid": order.OrderUID})
		rr := httptest.NewRecorder()

		h.GetOrder(rr, req)

		if rr.Code != http.StatusOK {
			b.Fatalf("unexpected status: %d", rr.Code)
		}
	}
}
