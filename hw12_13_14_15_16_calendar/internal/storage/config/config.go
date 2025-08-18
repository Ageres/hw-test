package storage_config

import (
	"context"
	"os"

	l "github.com/Ageres/hw-test/hw12_13_14_15_calendar/internal/logger"
	"github.com/Ageres/hw-test/hw12_13_14_15_calendar/internal/model"
	s "github.com/Ageres/hw-test/hw12_13_14_15_calendar/internal/storage"
	memorystorage "github.com/Ageres/hw-test/hw12_13_14_15_calendar/internal/storage/memory"
	sqlstorage "github.com/Ageres/hw-test/hw12_13_14_15_calendar/internal/storage/sql"
)

var ErrUnknowTypeStorageMsgTemplate = "unknow type storage: %s"

type StorageType string

const (
	IN_MEMORY StorageType = "IN_MEMORY"
	SQL       StorageType = "SQL"
)

func NewStorage(ctx context.Context, storageConfRef *model.StorageConf) s.Storage {
	logger := l.GetLogger(ctx)
	sType := storageConfRef.Type
	var storage s.Storage
	switch StorageType(sType) {
	case IN_MEMORY:
		storage = memorystorage.NewMemoryStorage(ctx, storageConfRef)
	case SQL:
		storage = sqlstorage.NewSQLStorage(ctx, storageConfRef)
	default:
		logger.Error("unknow type storage", map[string]any{
			"storageType": sType,
		})
		os.Exit(1)
	}
	logger.Info("storage configured", map[string]any{
		"storageType": sType,
	})
	return storage
}
