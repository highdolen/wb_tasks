package config

import (
	"log"

	wbfconfig "github.com/wb-go/wbf/config"
)

// LoadConfig - загрузка конфига
func LoadConfig(configFile string) (*AppConfig, error) {
	cfg := wbfconfig.New()

	// Загружаем YAML
	if configFile != "" {
		if err := cfg.LoadConfigFiles(configFile); err != nil {
			log.Println("Warning: could not load config file:", err)
		}
	}

	// Распаковка в структуру
	var appCfg AppConfig
	if err := cfg.Unmarshal(&appCfg); err != nil {
		return nil, err

	}

	return &appCfg, nil
}
