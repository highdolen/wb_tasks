package config

import (
	"log"

	wbfconfig "github.com/wb-go/wbf/config"
)

// LoadConfig загружает конфигурацию
func LoadConfig(configFile string) (*AppConfig, error) {
	cfg := wbfconfig.New()

	// Загрузка основной YAML конфиг
	if configFile != "" {
		if err := cfg.LoadConfigFiles(configFile); err != nil {
			log.Println("Warning: could not load config file:", err)
		}
	}

	// Распаковка конфига в структуру
	var appCfg AppConfig
	if err := cfg.Unmarshal(&appCfg); err != nil {
		log.Fatal("Failed to unmarshal config:", err)
	}

	return &appCfg, nil
}
