package config

import (
	"errors"
	"time"
)

// AppConfig - структура конфига приложения
type AppConfig struct {
	Server struct {
		Host string `mapstructure:"host"`
		Port int    `mapstructure:"port"`
	} `mapstructure:"server"`

	Postgres struct {
		Host            string        `mapstructure:"host"`
		Port            int           `mapstructure:"port"`
		User            string        `mapstructure:"user"`
		Password        string        `mapstructure:"password"`
		DBName          string        `mapstructure:"db"`
		SSLMode         string        `mapstructure:"sslmode"`
		MaxOpenConns    int           `mapstructure:"max_open_conns"`
		MaxIdleConns    int           `mapstructure:"max_idle_conns"`
		ConnMaxLifetime time.Duration `mapstructure:"conn_max_lifetime"`
	} `mapstructure:"postgres"`

	Telegram struct {
		BotToken string `mapstructure:"bot_token"`
	} `mapstructure:"telegram"`

	Scheduler struct {
		Interval time.Duration `mapstructure:"interval"`
	} `mapstructure:"scheduler"`
}

// ValidateConfig - валидация конфига
func ValidateConfig(cfg *AppConfig) error {
	if cfg.Postgres.Host == "" {
		return errors.New("postgres.host is required")
	}
	if cfg.Postgres.Port == 0 {
		return errors.New("postgres.port is required")
	}
	if cfg.Postgres.User == "" {
		return errors.New("postgres.user is required")
	}
	if cfg.Postgres.DBName == "" {
		return errors.New("postgres.db is required")
	}
	if cfg.Telegram.BotToken == "" {
		return errors.New("telegram.bot_token is required")
	}
	if cfg.Scheduler.Interval <= 0 {
		return errors.New("scheduler.interval must be > 0")
	}
	return nil
}
