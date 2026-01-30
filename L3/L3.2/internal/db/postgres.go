package db

import (
	"context"
	"fmt"
	"time"

	"shortener/internal/config"

	"github.com/wb-go/wbf/dbpg"
)

// New инициализирует подключение к Postgres
func New(cfg *config.AppConfig) (*dbpg.DB, error) {
	pg := cfg.Postgres

	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		pg.Host,
		pg.Port,
		pg.User,
		pg.Password,
		pg.DBName,
		pg.SSLMode,
	)

	opts := &dbpg.Options{
		MaxOpenConns:    pg.MaxOpenConns,
		MaxIdleConns:    pg.MaxIdleConns,
		ConnMaxLifetime: pg.ConnMaxLifetime,
	}

	db, err := dbpg.New(dsn, nil, opts)
	if err != nil {
		return nil, err
	}

	// Проверяем соединение
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.Master.PingContext(ctx); err != nil {
		_ = db.Master.Close()
		return nil, err
	}

	return db, nil
}
