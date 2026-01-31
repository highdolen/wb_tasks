package cache

import (
	"errors"
	"shortener/internal/config"

	wbfredis "github.com/wb-go/wbf/redis"
)

// New инициализирует Redis
func New(cfg *config.AppConfig) (*wbfredis.Client, error) {
	if !cfg.Redis.Enabled {
		return nil, nil
	}

	if cfg.Redis.Addr == "" {
		return nil, errors.New("redis.addr is required when redis is enabled")
	}

	client, err := wbfredis.Connect(wbfredis.Options{
		Address:  cfg.Redis.Addr,
		Password: cfg.Redis.Password,
	})
	if err != nil {
		return nil, err
	}
	return client, nil
}
