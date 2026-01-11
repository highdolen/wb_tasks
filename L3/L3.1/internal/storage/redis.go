package storage

import (
	"context"
	"encoding/json"
	"time"

	"delayed_notifier/internal/models"

	wbfredis "github.com/wb-go/wbf/redis"
)

type RedisStorage struct {
	client *wbfredis.Client
	ctx    context.Context
	ttl    time.Duration
}

func NewRedisStorage(client *wbfredis.Client, ttl time.Duration) *RedisStorage {
	return &RedisStorage{
		client: client,
		ctx:    context.Background(),
		ttl:    ttl,
	}
}

// Save сохраняет уведомление с TTL
func (r *RedisStorage) Save(n models.Notification) error {
	data, err := json.Marshal(n)
	if err != nil {
		return err
	}
	return r.client.SetWithExpiration(r.ctx, n.ID, string(data), r.ttl)
}

// Get получает уведомление по ID
func (r *RedisStorage) Get(id string) (models.Notification, bool) {
	data, err := r.client.Get(r.ctx, id)
	if err != nil {
		return models.Notification{}, false
	}
	var n models.Notification
	if err := json.Unmarshal([]byte(data), &n); err != nil {
		return models.Notification{}, false
	}
	return n, true
}

// Delete удаляет уведомление по ID
func (r *RedisStorage) Delete(id string) bool {
	if err := r.client.Del(r.ctx, id); err != nil {
		return false
	}
	return true
}
