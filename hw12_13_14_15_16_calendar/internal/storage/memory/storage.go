package memorystorage

import (
	"sync"
	"time"

	"github.com/Ageres/hw-test/hw12_13_14_15_calendar/internal/model"
	"github.com/Ageres/hw-test/hw12_13_14_15_calendar/internal/storage"
	"github.com/google/uuid"
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
	eventId := event.ID
	if event.ID != "" {
		_, exists := s.events[eventId]
		if exists {
			return storage.ErrEventAllreadyCreated
		}
	} else {
		event.ID = uuid.New().String()
	}
	if event.Title == "" {
		return storage.ErrEmptyTitle
	}

	return nil

}

func ListPeriodByUserId(start time.Time, duration time.Duration, userId string) ([]storage.Event, error) {
	return nil, nil
}

/*

 */
