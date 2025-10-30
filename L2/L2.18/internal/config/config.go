package config

import (
	"fmt"
	"os"
)

type Config struct {
	Server ServerConfig
}

type ServerConfig struct {
	Port string
}

func (c *Config) Validate() error {
	if c.Server.Port == "" {
		return fmt.Errorf("server port is missing")
	}
	return nil
}

func LoadConfig() (*Config, error) {
	port := os.Getenv("CALENDAR_PORT")
	if port == "" {
		port = "8080" // дефолтный порт
	}

	cfg := &Config{
		Server: ServerConfig{
			Port: port,
		},
	}

	return cfg, nil
}
