package cache

import (
	"sync"
	"time"

	"L0_optimize/internal/models"
)

// CacheEntry представляет запись в кеше с временной меткой
type CacheEntry struct {
	Order     models.Order
	Timestamp time.Time
}

type OrderCache struct {
	mu       sync.RWMutex
	cache    map[string]CacheEntry
	ttl      time.Duration
	stopChan chan bool
}

// New создает новый кэш с указанным TTL
func New(ttl time.Duration) *OrderCache {
	c := &OrderCache{
		cache:    make(map[string]CacheEntry),
		ttl:      ttl,
		stopChan: make(chan bool),
	}

	go c.cleanupExpired()

	return c
}

// Get — получить заказ по UID с проверкой TTL
func (c *OrderCache) Get(uid string) (models.Order, bool) {
	c.mu.RLock()
	entry, exists := c.cache[uid]
	c.mu.RUnlock()

	if !exists {
		return models.Order{}, false
	}

	if time.Since(entry.Timestamp) <= c.ttl {
		return entry.Order, true
	}

	c.mu.Lock()
	currentEntry, stillExists := c.cache[uid]
	if stillExists && currentEntry.Timestamp.Equal(entry.Timestamp) {
		delete(c.cache, uid)
	}
	c.mu.Unlock()

	return models.Order{}, false
}

// Set — добавить или обновить заказ с текущей временной меткой
func (c *OrderCache) Set(uid string, order models.Order) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.cache[uid] = CacheEntry{
		Order:     order,
		Timestamp: time.Now(),
	}
}

// Delete — удалить заказ по UID
func (c *OrderCache) Delete(uid string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.cache, uid)
}

// Invalidate — инвалидировать конкретный заказ
func (c *OrderCache) Invalidate(uid string) {
	c.Delete(uid)
}

// InvalidateAll — очистить весь кэш
func (c *OrderCache) InvalidateAll() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.cache = make(map[string]CacheEntry)
}

// GetStats — получить статистику кэша
func (c *OrderCache) GetStats() map[string]interface{} {
	c.mu.RLock()
	defer c.mu.RUnlock()

	totalEntries := len(c.cache)
	expiredEntries := 0
	now := time.Now()

	for _, entry := range c.cache {
		if now.Sub(entry.Timestamp) > c.ttl {
			expiredEntries++
		}
	}

	return map[string]interface{}{
		"total_entries":   totalEntries,
		"expired_entries": expiredEntries,
		"valid_entries":   totalEntries - expiredEntries,
		"ttl_minutes":     c.ttl.Minutes(),
	}
}

// cleanupExpired — горутина для очистки устаревших записей
func (c *OrderCache) cleanupExpired() {
	ticker := time.NewTicker(c.ttl / 2)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			c.mu.Lock()
			now := time.Now()

			for key, entry := range c.cache {
				if now.Sub(entry.Timestamp) >= c.ttl {
					delete(c.cache, key)
				}
			}

			c.mu.Unlock()

		case <-c.stopChan:
			return
		}
	}
}

// Close — остановить горутину очистки
func (c *OrderCache) Close() {
	close(c.stopChan)
}
