package calendar

import (
	"errors"
	"sort"
	"sync"
	"time"

	"secondBlock/L2.18/internal/models"
)

var (
	ErrEventNotFound = errors.New("event not found")
	ErrInvalidDate   = errors.New("invalid date format")
)

type Calendar struct {
	mu     sync.RWMutex
	events []models.Event
	nextID int
}

func New() *Calendar {
	return &Calendar{
		events: make([]models.Event, 0),
		nextID: 1,
	}
}

func (c *Calendar) CreateEvent(userID int, dateStr, title string) (*models.Event, error) {
	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return nil, ErrInvalidDate
	}

	if title == "" {
		return nil, errors.New("title cannot be empty")
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	event := models.Event{
		ID:     c.nextID,
		UserID: userID,
		Date:   date,
		Title:  title,
	}

	c.events = append(c.events, event)
	c.nextID++

	return &event, nil
}

func (c *Calendar) UpdateEvent(eventID, userID int, dateStr, title string) (*models.Event, error) {
	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return nil, ErrInvalidDate
	}

	if title == "" {
		return nil, errors.New("title cannot be empty")
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	for i, event := range c.events {
		if event.ID == eventID && event.UserID == userID {
			c.events[i].Date = date
			c.events[i].Title = title
			return &c.events[i], nil
		}
	}

	return nil, ErrEventNotFound
}

func (c *Calendar) DeleteEvent(eventID, userID int) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	for i, event := range c.events {
		if event.ID == eventID && event.UserID == userID {
			c.events = append(c.events[:i], c.events[i+1:]...)
			return nil
		}
	}

	return ErrEventNotFound
}

func (c *Calendar) GetEventsForDay(userID int, dateStr string) ([]models.Event, error) {
	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return nil, ErrInvalidDate
	}

	c.mu.RLock()
	defer c.mu.RUnlock()

	var result []models.Event
	for _, event := range c.events {
		if event.UserID == userID && isSameDay(event.Date, date) {
			result = append(result, event)
		}
	}

	sortEvents(result)
	return result, nil
}

func (c *Calendar) GetEventsForWeek(userID int, dateStr string) ([]models.Event, error) {
	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return nil, ErrInvalidDate
	}

	year, week := date.ISOWeek()

	c.mu.RLock()
	defer c.mu.RUnlock()

	var result []models.Event
	for _, event := range c.events {
		if event.UserID == userID {
			eventYear, eventWeek := event.Date.ISOWeek()
			if eventYear == year && eventWeek == week {
				result = append(result, event)
			}
		}
	}

	sortEvents(result)
	return result, nil
}

func (c *Calendar) GetEventsForMonth(userID int, dateStr string) ([]models.Event, error) {
	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return nil, ErrInvalidDate
	}

	year, month := date.Year(), date.Month()

	c.mu.RLock()
	defer c.mu.RUnlock()

	var result []models.Event
	for _, event := range c.events {
		if event.UserID == userID {
			eventYear, eventMonth := event.Date.Year(), event.Date.Month()
			if eventYear == year && eventMonth == month {
				result = append(result, event)
			}
		}
	}

	sortEvents(result)
	return result, nil
}

func isSameDay(t1, t2 time.Time) bool {
	y1, m1, d1 := t1.Date()
	y2, m2, d2 := t2.Date()
	return y1 == y2 && m1 == m2 && d1 == d2
}

func sortEvents(events []models.Event) {
	sort.Slice(events, func(i, j int) bool {
		return events[i].Date.Before(events[j].Date)
	})
}
