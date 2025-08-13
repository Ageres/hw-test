package sqlstorage

import (
	"context"

	"github.com/Ageres/hw-test/hw12_13_14_15_calendar/internal/model"
	"github.com/Ageres/hw-test/hw12_13_14_15_calendar/internal/storage"
)

type SqlStorage struct { // TODO
}

func NewSqlStorage(psqlConfRef *model.PSQLConfig) storage.Storage {
	return nil
}

func (s *SqlStorage) Connect(ctx context.Context) error {
	// TODO
	return nil
}

func (s *SqlStorage) Close(ctx context.Context) error {
	// TODO
	return nil
}
