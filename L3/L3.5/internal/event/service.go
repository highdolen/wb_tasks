package event

import (
	"context"
	"time"
)

// Repository - интерфейс
type Repository interface {
	Create(ctx context.Context, e *Event) error
	GetByID(ctx context.Context, id int64) (*Event, error)
	GetAll(ctx context.Context) ([]*Event, error)
}

type Service struct {
	repo Repository
}

// NewService - конструктор сервиса
func NewService(r Repository) *Service {
	return &Service{repo: r}
}

// CreateEvent - создание события
func (s *Service) CreateEvent(
	ctx context.Context,
	name string,
	date time.Time,
	totalSeats int,
	bookingTTL time.Duration,
) (*Event, error) {
	e := &Event{
		Name:           name,
		Date:           date,
		TotalSeats:     totalSeats,
		AvailableSeats: totalSeats,
		BookingTTL:     bookingTTL,
	}

	if err := s.repo.Create(ctx, e); err != nil {
		return nil, err
	}

	return e, nil
}

func (s *Service) GetEvent(ctx context.Context, id int64) (*Event, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *Service) GetAll(ctx context.Context) ([]*Event, error) {
	return s.repo.GetAll(ctx)
}
