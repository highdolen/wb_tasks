package reminder

import (
	"fmt"
	"time"

	"calendar/internal/logger"
	"calendar/internal/models"
)

// Service — сервис напоминаний
type Service struct {
	ch     chan models.Event
	logger *logger.Logger
}

// New - создает reminder-сервис и запускает фоновый воркер
func New(buffer int, logg *logger.Logger) *Service {
	s := &Service{
		ch:     make(chan models.Event, buffer),
		logger: logg,
	}

	go s.worker()

	return s
}

// Add - добавляет событие в очередь напоминаний
func (s *Service) Add(event models.Event) {
	select {
	case s.ch <- event:
	default:
		if s.logger != nil {
			s.logger.Log(fmt.Sprintf("reminder queue is full, event dropped: id=%d user=%d", event.ID, event.UserID))
		}
	}
}

// worker — фоновая горутина, которая обрабатывает очередь напоминаний
func (s *Service) worker() {
	for event := range s.ch {
		go func(e models.Event) {
			wait := time.Until(e.Reminder)

			if wait > 0 {
				time.Sleep(wait)
			}

			if s.logger != nil {
				s.logger.Log(fmt.Sprintf("REMINDER: %s for user %d", e.Title, e.UserID))
			}
		}(event)
	}
}
