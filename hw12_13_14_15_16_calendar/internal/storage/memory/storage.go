package memorystorage

import (
	"sync"
	"time"

	"github.com/Ageres/hw-test/hw12_13_14_15_calendar/internal/model"
	"github.com/Ageres/hw-test/hw12_13_14_15_calendar/internal/storage"
)

type MemoryStorage struct {
	mu     sync.RWMutex             //nolint:unused
	events map[string]storage.Event // key: Event.ID
}

func NewMemoryStorage(cfgRef *model.StorageConf) storage.Storage {
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

	if err := eventRef.Validate(); err != nil {
		return err
	}
	eventRef.CheckAndGenerateId()

	if _, exists := s.events[eventRef.ID]; exists {
		return storage.ErrEventAllreadyCreated
	}

	for _, createdEvent := range s.events {
		if createdEvent.UserID == eventRef.UserID &&
			checkTimeOverlap(createdEvent.StartTime, createdEvent.Duration, eventRef.StartTime, eventRef.Duration) {
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
			checkTimeOverlap(existingEvent.StartTime, existingEvent.Duration, newEventRef.StartTime, newEventRef.Duration) {
			return storage.ErrDateBusy
		}
	}

	s.events[id] = *newEventRef
	return nil
}

func (s *MemoryStorage) Delete(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if err := storage.ValidateId(id); err != nil {
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

	var result []storage.Event
	startOfDay := time.Date(start.Year(), start.Month(), start.Day(), 0, 0, 0, 0, start.Location())
	endOfDay := startOfDay.Add(24 * time.Hour)

	for _, event := range s.events {
		if !event.StartTime.Before(startOfDay) && event.StartTime.Before(endOfDay) {
			result = append(result, event)
			continue
		}
		eventEndTime := event.StartTime.Add(event.Duration)
		if !eventEndTime.Before(startOfDay) && eventEndTime.Before(endOfDay) {
			result = append(result, event)
		}
	}

	if len(result) == 0 {
		return nil, storage.ErrEventNotFound
	}

	return result, nil
}

func (s *MemoryStorage) ListWeek(start time.Time) ([]storage.Event, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var result []storage.Event
	startOfDay := time.Date(start.Year(), start.Month(), start.Day(), 0, 0, 0, 0, start.Location())
	endOfDay := startOfDay.Add(168 * time.Hour)

	for _, event := range s.events {
		if !event.StartTime.Before(startOfDay) && event.StartTime.Before(endOfDay) {
			result = append(result, event)
			continue
		}
		eventEndTime := event.StartTime.Add(event.Duration)
		if !eventEndTime.Before(startOfDay) && eventEndTime.Before(endOfDay) {
			result = append(result, event)
		}
	}

	if len(result) == 0 {
		return nil, storage.ErrEventNotFound
	}

	return result, nil
}

func (s *MemoryStorage) ListMonth(start time.Time) ([]storage.Event, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var result []storage.Event
	startOfDay := time.Date(start.Year(), start.Month(), start.Day(), 0, 0, 0, 0, start.Location())
	endOfDay := startOfDay.Add(5040 * time.Hour)

	for _, event := range s.events {
		if !event.StartTime.Before(startOfDay) && event.StartTime.Before(endOfDay) {
			result = append(result, event)
			continue
		}
		eventEndTime := event.StartTime.Add(event.Duration)
		if !eventEndTime.Before(startOfDay) && eventEndTime.Before(endOfDay) {
			result = append(result, event)
		}
	}

	if len(result) == 0 {
		return nil, storage.ErrEventNotFound
	}

	return result, nil
}

func checkTimeOverlap(start1 time.Time, duration1 time.Duration, start2 time.Time, duration2 time.Duration) bool {
	end1 := start1.Add(duration1)
	end2 := start2.Add(duration2)
	return start1.Before(end2) && end1.After(start2)
}
