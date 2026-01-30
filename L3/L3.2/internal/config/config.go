package config

import (
	"errors"
	"time"
)

type AppConfig struct {
	Server struct {
		Host string `mapstructure:"host"`
		Port int    `mapstructure:"port"`
	} `mapstructure:"server"`

	Redis struct {
		Enabled  bool          `mapstructure:"enabled"`
		Addr     string        `mapstructure:"addr"`
		Password string        `mapstructure:"password"`
		DB       int           `mapstructure:"db"`
		TTL      time.Duration `mapstructure:"ttl"`
	} `mapstructure:"redis"`

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
}

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
	return nil
}
