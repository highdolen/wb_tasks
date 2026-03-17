package calendar

import (
	"testing"
	"time"

	"calendar/internal/logger"
	"calendar/internal/reminder"
)

func newTestCalendar() *Calendar {
	logg := logger.New(10)
	rem := reminder.New(10, logg)
	return New(rem)
}

func TestCalendar_CreateAndGetEvents(t *testing.T) {
	cal := newTestCalendar()

	event, err := cal.CreateEvent(1, "2024-01-01", "New Year")
	if err != nil {
		t.Fatalf("CreateEvent failed: %v", err)
	}

	if event.Title != "New Year" {
		t.Errorf("Expected title 'New Year', got %q", event.Title)
	}

	events, err := cal.GetEventsForDay(1, "2024-01-01")
	if err != nil {
		t.Fatalf("GetEventsForDay failed: %v", err)
	}

	if len(events) != 1 {
		t.Errorf("Expected 1 event, got %d", len(events))
	}
}

func TestCalendar_UpdateEvent(t *testing.T) {
	cal := newTestCalendar()

	ev, err := cal.CreateEvent(1, "2024-01-01", "Old Title")
	if err != nil {
		t.Fatalf("CreateEvent failed: %v", err)
	}

	updated, err := cal.UpdateEvent(ev.ID, 1, "2024-01-02", "New Title")
	if err != nil {
		t.Fatalf("UpdateEvent failed: %v", err)
	}

	if updated.Title != "New Title" {
		t.Errorf("Expected 'New Title', got %q", updated.Title)
	}
}

func TestCalendar_DeleteEvent(t *testing.T) {
	cal := newTestCalendar()

	ev, err := cal.CreateEvent(1, "2024-01-01", "To delete")
	if err != nil {
		t.Fatalf("CreateEvent failed: %v", err)
	}

	err = cal.DeleteEvent(ev.ID, 1)
	if err != nil {
		t.Fatalf("DeleteEvent failed: %v", err)
	}

	err = cal.DeleteEvent(ev.ID, 1)
	if err == nil {
		t.Errorf("Expected error when deleting non-existing event")
	}
}

func TestCalendar_GetEventsForWeek(t *testing.T) {
	cal := newTestCalendar()

	if _, err := cal.CreateEvent(1, "2023-12-31", "Event 1"); err != nil {
		t.Fatalf("CreateEvent failed: %v", err)
	}
	if _, err := cal.CreateEvent(1, "2024-01-01", "Event 2"); err != nil {
		t.Fatalf("CreateEvent failed: %v", err)
	}
	if _, err := cal.CreateEvent(1, "2024-01-07", "Event 3"); err != nil {
		t.Fatalf("CreateEvent failed: %v", err)
	}
	if _, err := cal.CreateEvent(1, "2024-01-08", "Event 4"); err != nil {
		t.Fatalf("CreateEvent failed: %v", err)
	}

	events, err := cal.GetEventsForWeek(1, "2024-01-01")
	if err != nil {
		t.Fatalf("GetEventsForWeek failed: %v", err)
	}

	if len(events) != 2 {
		t.Errorf("Expected 2 events for week, got %d", len(events))
	}
}

func TestCalendar_GetEventsForMonth(t *testing.T) {
	cal := newTestCalendar()

	if _, err := cal.CreateEvent(1, "2024-01-05", "Event 1"); err != nil {
		t.Fatalf("CreateEvent failed: %v", err)
	}
	if _, err := cal.CreateEvent(1, "2024-01-20", "Event 2"); err != nil {
		t.Fatalf("CreateEvent failed: %v", err)
	}
	if _, err := cal.CreateEvent(1, "2024-02-01", "Event 3"); err != nil {
		t.Fatalf("CreateEvent failed: %v", err)
	}

	events, err := cal.GetEventsForMonth(1, "2024-01-10")
	if err != nil {
		t.Fatalf("GetEventsForMonth failed: %v", err)
	}

	if len(events) != 2 {
		t.Errorf("Expected 2 events for month, got %d", len(events))
	}
}

func TestCalendar_ArchiveOldEvents(t *testing.T) {
	cal := New(nil)

	yesterday := time.Now().Add(-24 * time.Hour).Format("2006-01-02")
	today := time.Now().Format("2006-01-02")

	_, _ = cal.CreateEvent(1, yesterday, "Old Event")
	_, _ = cal.CreateEvent(1, today, "Today Event")

	cal.ArchiveOldEvents()

	if len(cal.events) != 1 {
		t.Errorf("Expected 1 active event, got %d", len(cal.events))
	}
	if len(cal.archive) != 1 {
		t.Errorf("Expected 1 archived event, got %d", len(cal.archive))
	}
}

func TestArchiveOldEventsDirectly(t *testing.T) {
	cal := New(nil)

	yesterday := time.Now().Add(-24 * time.Hour).Format("2006-01-02")
	today := time.Now().Format("2006-01-02")
	tomorrow := time.Now().Add(24 * time.Hour).Format("2006-01-02")

	_, _ = cal.CreateEvent(1, yesterday, "Old Event")
	_, _ = cal.CreateEvent(1, today, "Today Event")
	_, _ = cal.CreateEvent(1, tomorrow, "Future Event")

	if len(cal.events) != 3 {
		t.Fatalf("Expected 3 events, got %d", len(cal.events))
	}

	cal.ArchiveOldEvents()

	if len(cal.events) != 2 {
		t.Errorf("Expected 2 active events, got %d", len(cal.events))
	}
	if len(cal.archive) != 1 {
		t.Errorf("Expected 1 archived event, got %d", len(cal.archive))
	}

	if cal.archive[0].Title != "Old Event" {
		t.Errorf("Expected archived event 'Old Event', got '%s'", cal.archive[0].Title)
	}
}
