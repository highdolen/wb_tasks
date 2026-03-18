package httpserver

import (
	"fmt"
	"net/http"

	"gcMetrics/internal/config"
	"gcMetrics/internal/debug"
	"gcMetrics/internal/profiler"
)

// Server - описывает HTTP-сервер приложения.
type Server struct {
	cfg     config.Config
	handler http.Handler
}

// New - создает HTTP-сервер и подготавливает маршруты приложения
func New(cfg *config.Config, metricsHandler http.Handler) *Server {
	return &Server{
		cfg:     *cfg,
		handler: buildMux(*cfg, metricsHandler),
	}
}

// Run - запускает HTTP-сервер
func (s *Server) Run() error {
	server := &http.Server{
		Addr:    s.cfg.Address(),
		Handler: s.handler,
	}

	return server.ListenAndServe()
}

// buildMux - собирает маршруты приложения
func buildMux(cfg config.Config, metricsHandler http.Handler) http.Handler {
	mux := http.NewServeMux()

	mux.Handle(cfg.MetricsPath, metricsHandler)
	mux.HandleFunc(cfg.HealthPath, healthHandler)
	debug.Register(mux)

	profiler.Register(mux)

	return logMiddleware(mux)
}

// healthHandler - обрабатывает healthcheck endpoint
func healthHandler(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("ok\n"))
}

// logMiddleware - логирует входящие HTTP-запросы
func logMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("%s %s\n", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}
