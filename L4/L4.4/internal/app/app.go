package app

import (
	"fmt"

	"gcMetrics/internal/config"
	"gcMetrics/internal/httpserver"
	"gcMetrics/internal/metrics"
)

// App хранит зависимости приложения: конфиг и HTTP-сервер.
type App struct {
	cfg    config.Config
	server *httpserver.Server
}

// New собирает приложение: загружает конфиг, инициализирует сборщик метрик,
// Prometheus handler и HTTP-сервер.
func New() (*App, error) {
	cfg, err := config.LoadConfig("config.yaml")
	if err != nil {
		return nil, err
	}

	collector := metrics.NewCollector()
	handler := metrics.NewPrometheusHandler(collector)

	server := httpserver.New(cfg, handler)

	return &App{
		cfg:    *cfg,
		server: server,
	}, nil
}

// Run запускает HTTP-сервер и выводит в консоль информацию о доступных endpoints.
func (a *App) Run() error {
	fmt.Printf("starting server on %s\n", a.cfg.Address())
	fmt.Printf("metrics endpoint: http://%s%s\n", a.cfg.Address(), a.cfg.MetricsPath)
	fmt.Printf("pprof endpoint:   http://%s/debug/pprof/\n", a.cfg.Address())

	return a.server.Run()
}
