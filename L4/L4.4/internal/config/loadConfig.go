package config

import (
	"fmt"
	"os"
	"runtime/debug"

	"gopkg.in/yaml.v3"
)

// LoadConfig - загружает конфигурацию из YAML-файла
func LoadConfig(configFile string) (*Config, error) {
	cfg := &Config{}

	data, err := os.ReadFile(configFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, fmt.Errorf("failed to parse yaml: %w", err)
	}

	// defaults (на всякий случай)
	setDefaults(cfg)

	// применяем GC
	debug.SetGCPercent(cfg.GCPercent)

	return cfg, nil
}

// setDefaults - устанавливает дефолтные значения
func setDefaults(cfg *Config) {
	if cfg.Host == "" {
		cfg.Host = "0.0.0.0"
	}
	if cfg.Port == "" {
		cfg.Port = "8080"
	}
	if cfg.MetricsPath == "" {
		cfg.MetricsPath = "/metrics"
	}
	if cfg.HealthPath == "" {
		cfg.HealthPath = "/healthz"
	}
	if cfg.GCPercent == 0 {
		cfg.GCPercent = 100
	}
}
