package memorystorage

import (
	"context"
	"sync"
	"time"

	"github.com/Ageres/hw-test/hw12_13_14_15_calendar/internal/storage"
)

type MemoryStorage struct {
	mu     sync.RWMutex
	events map[string]storage.Event // key: Event.ID
}

func NewMemoryStorage() storage.Storage {
	return &MemoryStorage{
		events: make(map[string]storage.Event),
	}
}

func (s *MemoryStorage) Add(ctx context.Context, eventRef *storage.Event) (*storage.Event, error) {
	if err := storage.ValidateEvent(eventRef); err != nil {
		return nil, err
	}
	eventRef.GenerateEventId()

	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.events[eventRef.ID]; exists {
		return nil, storage.NewSErrorWithTemplate("failed to add event, event with this id already exists: %s", eventRef.ID)
	}

	for _, existingEvent := range s.events {
		if existingEvent.UserID == eventRef.UserID &&
			existingEvent.Overlaps(eventRef) {
			return nil, storage.NewSErrorWithTemplate(storage.ErrDateBusyMsgTemplate, existingEvent.ID)
		}
	}

	s.events[eventRef.ID] = *eventRef
	result := *eventRef
	return &result, nil
}

func (s *MemoryStorage) Update(ctx context.Context, newEventRef *storage.Event) error {

	if err := storage.FullValidateEvent(newEventRef); err != nil {
		return err
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	id := newEventRef.ID
	oldEvent, exists := s.events[id]
	if !exists {
		return storage.ErrEventNotFound
	}
	if oldEvent.UserID != newEventRef.UserID {
		return storage.NewSErrorWithTemplate(storage.ErrUserConflictMsgTemplate, newEventRef.UserID, oldEvent.UserID)
	}

	for _, existingEvent := range s.events {
		if existingEvent.UserID == newEventRef.UserID &&
			existingEvent.ID != id &&
			existingEvent.Overlaps(newEventRef) {
			return storage.NewSErrorWithTemplate(storage.ErrDateBusyMsgTemplate, existingEvent.ID)
		}
	}

	s.events[id] = *newEventRef
	return nil
}

func (s *MemoryStorage) Delete(ctx context.Context, id string) error {
	if err := storage.ValidateEventId(id); err != nil {
		return err
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	_, exists := s.events[id]
	if !exists {
		return storage.ErrEventNotFound
	}

	delete(s.events, id)
	return nil
}

func (s *MemoryStorage) ListDay(ctx context.Context, startDay time.Time) ([]storage.Event, error) {
	startTime := time.Date(startDay.Year(), startDay.Month(), startDay.Day(), 0, 0, 0, 0, startDay.Location())
	endTime := startTime.Add(24 * time.Hour)
	return s.listEvents(ctx, startTime, endTime)
}

func (s *MemoryStorage) ListWeek(ctx context.Context, startDay time.Time) ([]storage.Event, error) {
	startTime := time.Date(startDay.Year(), startDay.Month(), startDay.Day(), 0, 0, 0, 0, startDay.Location())
	endTime := startTime.AddDate(0, 0, 7)
	return s.listEvents(ctx, startTime, endTime)
}

func (s *MemoryStorage) ListMonth(ctx context.Context, startDay time.Time) ([]storage.Event, error) {
	startTime := time.Date(startDay.Year(), startDay.Month(), startDay.Day(), 0, 0, 0, 0, startDay.Location())
	endTime := startTime.AddDate(0, 1, 0)
	return s.listEvents(ctx, startTime, endTime)
}

func (s *MemoryStorage) listEvents(ctx context.Context, startTime, endTime time.Time) ([]storage.Event, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var result []storage.Event
	for _, event := range s.events {

		select {
		case <-ctx.Done():
			return nil, storage.NewSErrorWithCause(storage.ErrContextDoneTemplate, ctx.Err())
		default:
		}

		eventEnd := event.StartTime.Add(event.Duration)
		if (event.StartTime.After(startTime) && event.StartTime.Before(endTime)) ||
			(eventEnd.After(startTime) && eventEnd.Before(endTime)) {
			result = append(result, event)
		}
	}
	return result, nil
}
