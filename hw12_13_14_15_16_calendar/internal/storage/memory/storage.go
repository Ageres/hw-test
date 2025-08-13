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

	if _, exists := s.events[event.ID]; exists {
		return storage.ErrEventAllreadyCreated
	}

	return nil

}

func ListPeriodByUserId(start time.Time, duration time.Duration, userId string) ([]storage.Event, error) {
	return nil, nil
}

/*

 */
