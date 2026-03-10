package storage

import (
	"fmt"
	"time"

	"warehouseControl/internal/config"

	"github.com/wb-go/wbf/dbpg"
)

// NewPostgres - создает новое подключение к базе данных PostgreSQL
func NewPostgres(cfg *config.AppConfig) (*dbpg.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.Postgres.Host,
		cfg.Postgres.Port,
		cfg.Postgres.User,
		cfg.Postgres.Password,
		cfg.Postgres.DBName,
		cfg.Postgres.SSLMode,
	)

	opts := &dbpg.Options{
		MaxOpenConns:    cfg.Postgres.MaxOpenConns,
		MaxIdleConns:    cfg.Postgres.MaxIdleConns,
		ConnMaxLifetime: cfg.Postgres.ConnMaxLifetime * time.Second,
	}

	return dbpg.New(dsn, nil, opts)
}
