package db

import (
	"context"
	"fmt"
	"log"
	"time"

	"eventBooker/internal/config"

	"github.com/wb-go/wbf/dbpg"
)

// New - инициализация подключения к Postgres
func New(cfg *config.AppConfig) (*dbpg.DB, error) {
	pg := cfg.Postgres

	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		pg.Host, pg.Port, pg.User, pg.Password, pg.DBName, pg.SSLMode,
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

	// Retry ping внутри storage
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	var pingErr error
	for i := 0; i < 10; i++ {
		pingErr = db.Master.PingContext(ctx)
		if pingErr == nil {
			break
		}
		log.Println("Postgres not ready, retrying...")
		time.Sleep(2 * time.Second)
	}
	if pingErr != nil {
		_ = db.Master.Close()
		return nil, fmt.Errorf("db ping error: %w", pingErr)
	}

	return db, nil
}
