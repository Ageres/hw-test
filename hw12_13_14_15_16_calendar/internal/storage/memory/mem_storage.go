package memorystorage

import (
	"sync"
	"time"

	"github.com/Ageres/hw-test/hw12_13_14_15_calendar/internal/storage"
)

type MemoryStorage struct {
	mu     sync.RWMutex             //nolint:unused
	events map[string]storage.Event // key: Event.ID
}

func NewMemoryStorage() storage.Storage {
	return &MemoryStorage{
		events: make(map[string]storage.Event),
	}
}

func (s *MemoryStorage) Add(eventRef *storage.Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if err := storage.ValidateEventNotNil(eventRef); err != nil {
		return err
	}

	if eventRef.ID != "" {
		if err := eventRef.FullValidate(); err != nil {
			return err
		}
	} else {
		if err := eventRef.Validate(); err != nil {
			return err
		}
		eventRef.GenerateId()
	}

	if _, exists := s.events[eventRef.ID]; exists {
		return storage.ErrEventAllreadyExists
	}

	for _, existingEvent := range s.events {
		if existingEvent.UserID == eventRef.UserID &&
			existingEvent.Overlaps(eventRef) {
			return storage.ErrDateBusy
		}
	}

	s.events[eventRef.ID] = *eventRef
	return nil
}

func (s *MemoryStorage) Update(id string, newEventRef *storage.Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if err := storage.ValidateEventNotNil(newEventRef); err != nil {
		return err
	}

	newEventRef.ID = id
	if err := newEventRef.FullValidate(); err != nil {
		return err
	}

	oldEvent, exists := s.events[id]
	if !exists {
		return storage.ErrEventNotFound
	}
	if oldEvent.UserID != newEventRef.UserID {
		return storage.ErrUserConflict
	}

	for _, existingEvent := range s.events {
		if existingEvent.UserID == newEventRef.UserID &&
			existingEvent.ID != id &&
			existingEvent.Overlaps(newEventRef) {
			return storage.ErrDateBusy
		}
	}

	s.events[id] = *newEventRef
	return nil
}

func (s *MemoryStorage) Delete(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if err := storage.ValidateEventId(id); err != nil {
		return err
	}

	_, exists := s.events[id]
	if !exists {
		return storage.ErrEventNotFound
	}

	delete(s.events, id)
	return nil
}

func (s *MemoryStorage) ListDay(start time.Time) ([]storage.Event, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	startTime := time.Date(start.Year(), start.Month(), start.Day(), 0, 0, 0, 0, start.Location())
	endTime := startTime.Add(24 * time.Hour)

	result := s.getEventsByPeriod(startTime, endTime)
	return result, nil
}

func (s *MemoryStorage) ListWeek(start time.Time) ([]storage.Event, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	startTime := time.Date(start.Year(), start.Month(), start.Day(), 0, 0, 0, 0, start.Location())
	endTime := startTime.AddDate(0, 0, 7)

	result := s.getEventsByPeriod(startTime, endTime)
	return result, nil
}

func (s *MemoryStorage) ListMonth(start time.Time) ([]storage.Event, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	startTime := time.Date(start.Year(), start.Month(), start.Day(), 0, 0, 0, 0, start.Location())
	endTime := startTime.AddDate(0, 1, 0)

	result := s.getEventsByPeriod(startTime, endTime)
	return result, nil
}

func (s *MemoryStorage) getEventsByPeriod(startTime, endTime time.Time) []storage.Event {
	var result []storage.Event
	for _, event := range s.events {
		eventEnd := event.StartTime.Add(event.Duration)
		if (event.StartTime.After(startTime) && event.StartTime.Before(endTime)) ||
			(eventEnd.After(startTime) && eventEnd.Before(endTime)) {
			result = append(result, event)
		}
	}
	return result
}
