package memorystorage

import (
	"sync"

	"github.com/Ageres/hw-test/hw12_13_14_15_calendar/internal/model"
	"github.com/Ageres/hw-test/hw12_13_14_15_calendar/internal/storage"
)

type MemoryStorage struct {
	storage.Storage
	mu sync.RWMutex //nolint:unused
}

func New(cfgRef *model.StorageConf) *MemoryStorage {
	return &MemoryStorage{}
}

// TODO
