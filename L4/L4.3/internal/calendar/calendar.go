package calendar

import (
	"errors"
	"sort"
	"sync"
	"time"

	"calendar/internal/models"
	"calendar/internal/reminder"
)

var ErrEventNotFound = errors.New("event not found")
var ErrInvalidDate = errors.New("invalid date format")

// Calendar — основной сервис календаря
type Calendar struct {
	mu       sync.RWMutex
	events   []models.Event
	archive  []models.Event
	nextID   int
	reminder *reminder.Service
}

// New создает новый экземпляр календаря
func New(rem *reminder.Service) *Calendar {
	return &Calendar{
		events:   make([]models.Event, 0),
		archive:  make([]models.Event, 0),
		nextID:   1,
		reminder: rem,
	}
}

// CreateEvent - создает новое событие без напоминания
func (c *Calendar) CreateEvent(userID int, dateStr, title string) (*models.Event, error) {
	loc := time.Now().Location()
	date, err := time.ParseInLocation("2006-01-02", dateStr, loc)
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

	result := event
	return &result, nil
}

// CreateEventWithReminder - создает новое событие с напоминанием
func (c *Calendar) CreateEventWithReminder(userID int, dateStr, title, reminderStr string) (*models.Event, error) {
	loc := time.Now().Location()

	date, err := time.ParseInLocation("2006-01-02", dateStr, loc)
	if err != nil {
		return nil, ErrInvalidDate
	}

	if title == "" {
		return nil, errors.New("title cannot be empty")
	}

	var reminderTime time.Time
	if reminderStr != "" {
		reminderTime, err = time.ParseInLocation("2006-01-02T15:04:05", reminderStr, loc)
		if err != nil {
			return nil, errors.New("invalid reminder format")
		}
	}

	c.mu.Lock()

	event := models.Event{
		ID:       c.nextID,
		UserID:   userID,
		Date:     date,
		Title:    title,
		Reminder: reminderTime,
	}

	c.events = append(c.events, event)
	c.nextID++

	c.mu.Unlock()

	if !reminderTime.IsZero() && c.reminder != nil {
		c.reminder.Add(event)
	}

	result := event
	return &result, nil
}

// UpdateEvent - обновляет дату и название существующего события
func (c *Calendar) UpdateEvent(eventID, userID int, dateStr, title string) (*models.Event, error) {
	loc := time.Now().Location()
	date, err := time.ParseInLocation("2006-01-02", dateStr, loc)
	if err != nil {
		return nil, ErrInvalidDate
	}

	if title == "" {
		return nil, errors.New("title cannot be empty")
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	for i := range c.events {
		if c.events[i].ID == eventID && c.events[i].UserID == userID {
			c.events[i].Date = date
			c.events[i].Title = title

			result := c.events[i]
			return &result, nil
		}
	}

	return nil, ErrEventNotFound
}

// DeleteEvent - удаляет событие по ID и userID
func (c *Calendar) DeleteEvent(eventID, userID int) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	for i := range c.events {
		if c.events[i].ID == eventID && c.events[i].UserID == userID {
			c.events = append(c.events[:i], c.events[i+1:]...)
			return nil
		}
	}

	return ErrEventNotFound
}

// GetEventsForDay - возвращает все события пользователя за указанный день
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

// GetEventsForWeek - возвращает все события пользователя за указанную неделю
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

// GetEventsForMonth - возвращает все события пользователя за указанный месяц
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

// GetArchivedEvents - возвращает архивные события конкретного пользователя
func (c *Calendar) GetArchivedEvents(userID int) []models.Event {
	c.mu.RLock()
	defer c.mu.RUnlock()

	result := make([]models.Event, 0)
	for _, event := range c.archive {
		if event.UserID == userID {
			result = append(result, event)
		}
	}

	sortEvents(result)
	return result
}

// isSameDay - проверяет, относятся ли две даты к одному календарному дню
func isSameDay(t1, t2 time.Time) bool {
	y1, m1, d1 := t1.Date()
	y2, m2, d2 := t2.Date()
	return y1 == y2 && m1 == m2 && d1 == d2
}

// sortEvents - сортирует события по дате по возрастанию
func sortEvents(events []models.Event) {
	sort.Slice(events, func(i, j int) bool {
		return events[i].Date.Before(events[j].Date)
	})
}

// ArchiveOldEvents - переносит все события с датой раньше сегодняшней в архив
func (c *Calendar) ArchiveOldEvents() {
	c.mu.Lock()
	defer c.mu.Unlock()

	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	var active []models.Event

	for _, e := range c.events {
		eventDay := time.Date(e.Date.Year(), e.Date.Month(), e.Date.Day(), 0, 0, 0, 0, e.Date.Location())

		if eventDay.Before(today) {
			c.archive = append(c.archive, e)
		} else {
			active = append(active, e)
		}
	}

	c.events = active
}
