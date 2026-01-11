package storage

import (
	"delayed_notifier/internal/models"
	"sync"
)

type NotificationStorage interface {
	Save(n models.Notification) error
	Get(id string) (models.Notification, bool)
	Delete(id string) bool
}

type MemoryStorage struct {
	mu   sync.RWMutex
	data map[string]models.Notification
}

func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		data: make(map[string]models.Notification),
	}
}

func (m *MemoryStorage) Save(n models.Notification) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.data[n.ID] = n
	return nil
}

func (m *MemoryStorage) Get(id string) (models.Notification, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	v, ok := m.data[id]
	return v, ok
}

func (m *MemoryStorage) Delete(id string) bool {
	m.mu.Lock()
	defer m.mu.Unlock()

	_, exists := m.data[id]
	if !exists {
		return false
	}

	delete(m.data, id)
	return true
}
