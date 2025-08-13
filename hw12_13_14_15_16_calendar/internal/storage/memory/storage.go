package memorystorage

import (
	"sync"

	"github.com/Ageres/hw-test/hw12_13_14_15_calendar/internal/model"
)

type Storage struct {
	// TODO
	mu sync.RWMutex //nolint:unused
}

func New(cfgRef *model.StorageConf) *Storage {
	return &Storage{}
}

// TODO
