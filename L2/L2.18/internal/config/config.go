package config

import "fmt"

type Config struct {
	ServerConfig
}

type ServerConfig struct {
	Port string
}

// Validate проверяет, что все обязательные поля заполнены
func (c *Config) Validate() error {
	if c.Port == "" {
		return fmt.Errorf("server port is missing")
	}

	return nil
}
