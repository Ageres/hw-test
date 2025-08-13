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

func NewMemoryStorage(cfgRef *model.StorageConf) *MemoryStorage {
	return &MemoryStorage{
		events: make(map[string]storage.Event),
	}
}

func (s *MemoryStorage) Add(event storage.Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if err := event.Validate(); err != nil {
		return err
	}
	event.CheckAndGenerateId()

	if _, exists := s.events[event.ID]; exists {
		return storage.ErrEventAllreadyCreated
	}

	for _, createdEvent := range s.events {
		if createdEvent.UserID == event.UserID &&
			checkTimeOverlap(createdEvent.StartTime, createdEvent.Duration, event.StartTime, event.Duration) {
			return storage.ErrDateBusy
		}
	}

	s.events[event.ID] = event

	return nil

}

func ListPeriodByUserId(start time.Time, duration time.Duration, userId string) ([]storage.Event, error) {
	return nil, nil
}

func checkTimeOverlap(start1 time.Time, duration1 time.Duration, start2 time.Time, duration2 time.Duration) bool {
	end1 := start1.Add(duration1)
	end2 := start2.Add(duration2)
	return start1.Before(end2) && end1.After(start2)
}
