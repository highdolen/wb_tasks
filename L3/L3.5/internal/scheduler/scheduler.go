package scheduler

import (
	"context"
	"log"
	"time"

	"eventBooker/internal/booking"
)

type ExpirationScheduler struct {
	bookingService *booking.Service
	interval       time.Duration
}

// NewExpirationScheduler - конструктор ExpirationScheduler
func NewExpirationScheduler(
	bookingService *booking.Service,
	interval time.Duration,
) *ExpirationScheduler {
	return &ExpirationScheduler{
		bookingService: bookingService,
		interval:       interval,
	}
}

// Start - запуск фонового обработчика, отменяющий неоплаченные бронирования
func (s *ExpirationScheduler) Start(ctx context.Context) {
	ticker := time.NewTicker(s.interval)

	go func() {
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				log.Println("expiration scheduler stopped")
				return

			case <-ticker.C:
				err := s.bookingService.CancelExpiredBookings(ctx)
				if err != nil {
					log.Printf("scheduler error: %v", err)
				}
			}
		}
	}()
}
