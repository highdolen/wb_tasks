package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	"L0_optimize/internal/service"

	"github.com/gorilla/mux"
)

type OrderHandler struct {
	orderService service.OrderService
}

func NewOrderHandler(orderService service.OrderService) *OrderHandler {
	return &OrderHandler{
		orderService: orderService,
	}
}

func (h *OrderHandler) GetOrder(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	uid := vars["order_uid"]

	forceRefresh := r.URL.Query().Get("refresh") == "true"

	var result *service.OrderResult
	var err error

	if forceRefresh {
		result, err = h.orderService.GetOrderByUIDWithRefresh(r.Context(), uid)
	} else {
		result, err = h.orderService.GetOrderByUID(r.Context(), uid)
	}

	if err != nil {
		if errors.Is(err, service.ErrOrderNotFound) {
			http.Error(w, "Заказ не найден", http.StatusNotFound)
		} else {
			http.Error(w, "Ошибка при получении заказа: "+err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if result.FromCache {
		w.Header().Set("X-Cache", "HIT")
	} else {
		w.Header().Set("X-Cache", "MISS")
	}

	if err := json.NewEncoder(w).Encode(result.Order); err != nil {
		http.Error(w, "Ошибка при кодировании ответа", http.StatusInternalServerError)
	}
}

// GetCacheStats — получить статистику кеша
func (h *OrderHandler) GetCacheStats(w http.ResponseWriter, r *http.Request) {
	stats := h.orderService.GetCacheStats()
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(stats); err != nil {
		http.Error(w, "Ошибка при кодировании ответа", http.StatusInternalServerError)
	}
}

// InvalidateCache — инвалидировать весь кеш или конкретный заказ
func (h *OrderHandler) InvalidateCache(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	uid := vars["order_uid"]

	var err error
	if uid != "" {
		err = h.orderService.InvalidateCache(uid)
		if err != nil {
			http.Error(w, "Ошибка при инвалидации кеша: "+err.Error(), http.StatusInternalServerError)
			return
		}
		writeJSON(w, map[string]string{"message": fmt.Sprintf("Кеш для заказа %s инвалидирован", uid)}, http.StatusOK)
	} else {
		err = h.orderService.InvalidateAllCache()
		if err != nil {
			http.Error(w, "Ошибка при полной инвалидации кеша: "+err.Error(), http.StatusInternalServerError)
			return
		}
		writeJSON(w, map[string]string{"message": "Весь кеш инвалидирован"}, http.StatusOK)
	}
}

func writeJSON(w http.ResponseWriter, v interface{}, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		log.Printf("Ошибка кодирования JSON: %v", err)
	}
}
