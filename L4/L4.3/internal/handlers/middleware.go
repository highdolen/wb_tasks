package handlers

import (
	"fmt"
	"net/http"
	"time"

	"calendar/internal/logger"
)

// statusRecorder - оборачивает ResponseWriter
type statusRecorder struct {
	http.ResponseWriter
	status int
	size   int
}

// WriteHeader - сохраняет HTTP-статус ответа
func (r *statusRecorder) WriteHeader(code int) {
	r.status = code
	r.ResponseWriter.WriteHeader(code)
}

// Write - записывает тело ответа и считает его размер в байтах
func (r *statusRecorder) Write(b []byte) (int, error) {
	if r.status == 0 {
		r.status = http.StatusOK
	}
	n, err := r.ResponseWriter.Write(b)
	r.size += n
	return n, err
}

// Middleware - хранит зависимости middleware-слоя, в том числе логгер
type Middleware struct {
	logger *logger.Logger
}

// NewMiddleware - создает новый экземпляр middleware
func NewMiddleware(l *logger.Logger) *Middleware {
	return &Middleware{logger: l}
}

// Logging — middleware для асинхронного логирования HTTP-запросов
func (m *Middleware) Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		rec := &statusRecorder{ResponseWriter: w}

		next.ServeHTTP(rec, r)

		if rec.status == 0 {
			rec.status = http.StatusOK
		}

		m.logger.Log(
			fmt.Sprintf(
				"%s %s %d %s",
				r.Method,
				r.URL.Path,
				rec.status,
				time.Since(start),
			),
		)
	})
}
