package config

import (
	"errors"
)

// AppConfig - структура конфига
type AppConfig struct {
	Server struct {
		Host string `mapstructure:"host"`
		Port int    `mapstructure:"port"`
	} `mapstructure:"server"`

	Kafka struct {
		Brokers []string `mapstructure:"brokers"`
		Topic   string   `mapstructure:"topic"`
		GroupID string   `mapstructure:"group_id"`
	} `mapstructure:"kafka"`

	Storage struct {
		OriginalPath  string `mapstructure:"original_path"`
		ProcessedPath string `mapstructure:"processed_path"`
	} `mapstructure:"storage"`

	Processor struct {
		ResizeWidth    int    `mapstructure:"resize_width"`
		ThumbnailWidth int    `mapstructure:"thumbnail_width"`
		WatermarkText  string `mapstructure:"watermark_text"`
	} `mapstructure:"processor"`
}

// ValidateConfig - валидация конфига
func ValidateConfig(cfg *AppConfig) error {
	if cfg.Server.Host == "" || cfg.Server.Port == 0 {
		return errors.New("invalid server config")
	}

	if len(cfg.Kafka.Brokers) == 0 || cfg.Kafka.Topic == "" || cfg.Kafka.GroupID == "" {
		return errors.New("invalid kafka config")
	}

	if cfg.Storage.OriginalPath == "" || cfg.Storage.ProcessedPath == "" {
		return errors.New("invalid storage config")
	}

	return nil
}
