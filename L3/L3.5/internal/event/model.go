package event

import "time"

//Event - структура события
type Event struct {
	ID             int64         `json:"id"`
	Name           string        `json:"name"`
	Date           time.Time     `json:"date"`
	TotalSeats     int           `json:"total_seats"`
	AvailableSeats int           `json:"available_seats"`
	BookingTTL     time.Duration `json:"booking_ttl"` // Duration, в секундах конвертируем при записи в репозиторий
	CreatedAt      time.Time     `json:"created_at"`
}
