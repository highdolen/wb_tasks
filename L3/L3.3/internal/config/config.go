package config

import (
	"errors"
)

type AppConfig struct {
	Server struct {
		Host string `mapstructure:"host"`
		Port int    `mapstructure:"port"`
	} `mapstructure:"server"`
}

func ValidateConfig(cfg *AppConfig) error {
	if cfg.Server.Host == "" {
		return errors.New("host is required")
	}
	if cfg.Server.Port == 0 {
		return errors.New("port is required")
	}

	return nil
}
