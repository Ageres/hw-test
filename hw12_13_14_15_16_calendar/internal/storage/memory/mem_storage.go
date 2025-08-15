package memorystorage

import (
	"context"
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

func (s *MemoryStorage) Add(ctx context.Context, eventRef *storage.Event) (*storage.Event, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if err := storage.ValidateEventNotNil(eventRef); err != nil {
		return nil, err
	}

	if eventRef.ID != "" {
		if err := eventRef.FullValidate(); err != nil {
			return nil, err
		}
	} else {
		if err := eventRef.Validate(); err != nil {
			return nil, err
		}
		eventRef.GenerateId()
	}

	if _, exists := s.events[eventRef.ID]; exists {
		return nil, storage.ErrEventAllreadyExists
	}

	for _, existingEvent := range s.events {
		if existingEvent.UserID == eventRef.UserID &&
			existingEvent.Overlaps(eventRef) {
			return nil, storage.ErrDateBusy
		}
	}

	s.events[eventRef.ID] = *eventRef
	result := *eventRef
	return &result, nil
}

func (s *MemoryStorage) Update(ctx context.Context, newEventRef *storage.Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if err := storage.ValidateEventNotNil(newEventRef); err != nil {
		return err
	}

	if err := newEventRef.FullValidate(); err != nil {
		return err
	}

	id := newEventRef.ID
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

func (s *MemoryStorage) Delete(ctx context.Context, id string) error {
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

func (s *MemoryStorage) ListDay(ctx context.Context, startDay time.Time) ([]storage.Event, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	startTime := time.Date(startDay.Year(), startDay.Month(), startDay.Day(), 0, 0, 0, 0, startDay.Location())
	endTime := startTime.Add(24 * time.Hour)

	result := s.getEventsByPeriod(ctx, startTime, endTime)
	return result, nil
}

func (s *MemoryStorage) ListWeek(ctx context.Context, startDay time.Time) ([]storage.Event, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	startTime := time.Date(startDay.Year(), startDay.Month(), startDay.Day(), 0, 0, 0, 0, startDay.Location())
	endTime := startTime.AddDate(0, 0, 7)

	result := s.getEventsByPeriod(ctx, startTime, endTime)
	return result, nil
}

func (s *MemoryStorage) ListMonth(ctx context.Context, startDay time.Time) ([]storage.Event, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	startTime := time.Date(startDay.Year(), startDay.Month(), startDay.Day(), 0, 0, 0, 0, startDay.Location())
	endTime := startTime.AddDate(0, 1, 0)

	result := s.getEventsByPeriod(ctx, startTime, endTime)
	return result, nil
}

func (s *MemoryStorage) getEventsByPeriod(ctx context.Context, startTime, endTime time.Time) []storage.Event {
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
