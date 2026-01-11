package config

import (
	"errors"
	"log"

	"github.com/wb-go/wbf/rabbitmq"
)

type AppConfig struct {
	Server struct {
		Host string `mapstructure:"host"`
		Port int    `mapstructure:"port"`
	} `mapstructure:"server"`

	RabbitMQ struct {
		URL string `mapstructure:"url"`
	} `mapstructure:"rabbitmq"`

	Redis struct {
		URL string `mapstructure:"url"`
	} `mapstructure:"redis"`

	SMTP struct {
		Host     string `mapstructure:"host"`
		Port     int    `mapstructure:"port"`
		Username string `mapstructure:"username"`
		Password string `mapstructure:"password"`
	} `mapstructure:"smtp"`

	Telegram struct {
		Token string `mapstructure:"token"`
	} `mapstructure:"telegram"`
}

func ValidateConfig(cfg *AppConfig) error {
	if cfg.RabbitMQ.URL == "" {
		return rabbitmq.ErrMissingURL
	}

	if cfg.Server.Port == 0 {
		return errors.New("server port cannot be zero")

	}

	if cfg.Redis.URL == "" {
		log.Println("warning: redis URL is empty — redis features disabled")
	}

	if cfg.SMTP.Host == "" {
		log.Println("smpt host cannot be zero")
	}

	if cfg.Telegram.Token == "" {
		log.Println("telegram token cannot be zero")
	}
	return nil
}
