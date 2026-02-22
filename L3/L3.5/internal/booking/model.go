package booking

import "time"

type Status string

const (
	StatusPending   Status = "pending"
	StatusConfirmed Status = "confirmed"
	StatusCanceled  Status = "canceled"
)

//Booking - структура бронирования
type Booking struct {
	ID        int64
	EventID   int64
	Status    Status
	CreatedAt time.Time
	ExpiresAt time.Time
}
