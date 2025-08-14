package sqlstorage

import (
	"context"
	"log"

	"github.com/Ageres/hw-test/hw12_13_14_15_calendar/internal/model"
	"github.com/Ageres/hw-test/hw12_13_14_15_calendar/internal/storage"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

type SqlStorage struct {
	db *sqlx.DB
}

func NewSqlStorage(psqlConfRef *model.PSQLConfig) storage.Storage {
	dsn := psqlConfRef.DB.DSN()
	log.Println("dsc:", dsn)

	db, err := sqlx.Connect("pgx", dsn)
	if err != nil {
		log.Fatalf("failed to load driver: %v", err)
	}
	err = db.Ping()
	if err != nil {
		log.Fatalf("failed to connect to db: %w", err)
	}

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
