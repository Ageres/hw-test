package storage_config

import (
	"log"
	"os"

	"github.com/Ageres/hw-test/hw12_13_14_15_calendar/internal/model"
	"github.com/Ageres/hw-test/hw12_13_14_15_calendar/internal/storage"
	memorystorage "github.com/Ageres/hw-test/hw12_13_14_15_calendar/internal/storage/memory"
	sqlstorage "github.com/Ageres/hw-test/hw12_13_14_15_calendar/internal/storage/sql"
)

var ErrUnknowTypeStorageMsgTemplate = "unknow type storage: %s"

type StorageType string

const (
	IN_MEMORY StorageType = "IN_MEMORY"
	SQL       StorageType = "SQL"
)

func NewStorage(storageConfRef *model.StorageConf) storage.Storage {
	sType := storageConfRef.Type
	switch StorageType(sType) {
	case IN_MEMORY:
		return memorystorage.NewMemoryStorage()
	case SQL:
		return sqlstorage.NewSqlStorage(storageConfRef.PSQL)
	default:
		log.Printf(ErrUnknowTypeStorageMsgTemplate, sType)
		os.Exit(1)
	}
	return nil
}
