package config

import (
	"log"
	"os"
	"strconv"

	wbfconfig "github.com/wb-go/wbf/config"
)

// LoadConfig - загрузка конфига
func LoadConfig(configFile string) (*AppConfig, error) {
	cfg := wbfconfig.New()

	// Загружаем .env
	_ = cfg.LoadEnvFiles(".env")
	cfg.EnableEnv("")

	// Загружаем YAML
	if configFile != "" {
		if err := cfg.LoadConfigFiles(configFile); err != nil {
			log.Println("Warning: could not load config file:", err)
		}
	}

	var appCfg AppConfig
	if err := cfg.Unmarshal(&appCfg); err != nil {
		return nil, err
	}

	// ===== PostgreSQL ENV =====
	if v := os.Getenv("POSTGRES_HOST"); v != "" {
		appCfg.Postgres.Host = v
	}
	if v := os.Getenv("POSTGRES_PORT"); v != "" {
		if port, err := strconv.Atoi(v); err == nil {
			appCfg.Postgres.Port = port
		}
	}
	if v := os.Getenv("POSTGRES_USER"); v != "" {
		appCfg.Postgres.User = v
	}
	if v := os.Getenv("POSTGRES_PASSWORD"); v != "" {
		appCfg.Postgres.Password = v
	}
	if v := os.Getenv("POSTGRES_DB"); v != "" {
		appCfg.Postgres.DBName = v
	}
	if appCfg.Postgres.DBName == "" {
		appCfg.Postgres.DBName = "control_events" // дефолт на случай пустого env
	}
	if v := os.Getenv("POSTGRES_SSLMODE"); v != "" {
		appCfg.Postgres.SSLMode = v
	}

	return &appCfg, nil
}
