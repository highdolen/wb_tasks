package config

import (
	"log"
	"os"
	"strconv"
	"time"

	wbfconfig "github.com/wb-go/wbf/config"
)

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

	// Распаковка в структуру
	var appCfg AppConfig
	if err := cfg.Unmarshal(&appCfg); err != nil {
		return nil, err
	}

	// Гарантированно подтягиваем ENV для Postgres
	if v := os.Getenv("POSTGRES_HOST"); v != "" {
		appCfg.Postgres.Host = v
	}
	if v := os.Getenv("POSTGRES_PORT"); v != "" {
		port, _ := strconv.Atoi(v)
		appCfg.Postgres.Port = port
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

	// ENV для Redis
	if v := os.Getenv("REDIS_ADDR"); v != "" {
		appCfg.Redis.Addr = v
	}
	if v := os.Getenv("REDIS_TTL"); v != "" {
		dur, _ := time.ParseDuration(v + "s")
		appCfg.Redis.TTL = dur
	}

	return &appCfg, nil
}
