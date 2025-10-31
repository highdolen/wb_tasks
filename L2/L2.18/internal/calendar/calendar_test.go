package calendar

import (
	"testing"
)

func TestCalendar_CreateAndGetEvents(t *testing.T) {
	cal := New()

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
	cal := New()
	ev, _ := cal.CreateEvent(1, "2024-01-01", "Old Title")

	updated, err := cal.UpdateEvent(ev.ID, 1, "2024-01-02", "New Title")
	if err != nil {
		t.Fatalf("UpdateEvent failed: %v", err)
	}

	if updated.Title != "New Title" {
		t.Errorf("Expected 'New Title', got %q", updated.Title)
	}
}

func TestCalendar_DeleteEvent(t *testing.T) {
	cal := New()
	ev, _ := cal.CreateEvent(1, "2024-01-01", "To delete")

	err := cal.DeleteEvent(ev.ID, 1)
	if err != nil {
		t.Fatalf("DeleteEvent failed: %v", err)
	}

	err = cal.DeleteEvent(ev.ID, 1)
	if err == nil {
		t.Errorf("Expected error when deleting non-existing event")
	}
}

func TestCalendar_GetEventsForWeek(t *testing.T) {
	cal := New()

	// События вокруг Нового года, ISO-недели пересекаются
	cal.CreateEvent(1, "2023-12-31", "Event 1") // неделя 52
	cal.CreateEvent(1, "2024-01-01", "Event 2") // неделя 1
	cal.CreateEvent(1, "2024-01-07", "Event 3") // неделя 1
	cal.CreateEvent(1, "2024-01-08", "Event 4") // неделя 2

	events, err := cal.GetEventsForWeek(1, "2024-01-01")
	if err != nil {
		t.Fatalf("GetEventsForWeek failed: %v", err)
	}

	if len(events) != 2 {
		t.Errorf("Expected 2 events for week, got %d", len(events))
	}
}

func TestCalendar_GetEventsForMonth(t *testing.T) {
	cal := New()

	cal.CreateEvent(1, "2024-01-05", "Event 1")
	cal.CreateEvent(1, "2024-01-20", "Event 2")
	cal.CreateEvent(1, "2024-02-01", "Event 3")

	events, err := cal.GetEventsForMonth(1, "2024-01-10")
	if err != nil {
		t.Fatalf("GetEventsForMonth failed: %v", err)
	}

	if len(events) != 2 {
		t.Errorf("Expected 2 events for month, got %d", len(events))
	}
}
