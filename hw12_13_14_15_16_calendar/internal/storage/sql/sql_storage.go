package sqlstorage

import (
	"context"
	"log"
	"time"

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

	db.SetMaxOpenConns(psqlConfRef.Pool.Conn.MaxOpen)
	db.SetMaxIdleConns(psqlConfRef.Pool.Conn.MaxIdle)
	db.SetConnMaxLifetime(time.Duration(psqlConfRef.Pool.Conn.MaxLifeTime) * time.Second)
	db.SetConnMaxIdleTime(time.Duration(psqlConfRef.Pool.Conn.MaxLifeTime) * time.Second)

	err = db.Ping()
	if err != nil {
		log.Fatalf("failed to connect to db: %v", err)
	}

	return &SqlStorage{
		db: db,
	}
}

func (s *SqlStorage) Close(ctx context.Context) error {
	return s.db.Close()
}

// Add implements storage.Storage.
func (s *SqlStorage) Add(eventRef *storage.Event) error {
	panic("unimplemented")
}

// Delete implements storage.Storage.
func (s *SqlStorage) Delete(id string) error {
	panic("unimplemented")
}

// ListDay implements storage.Storage.
func (s *SqlStorage) ListDay(start time.Time) ([]storage.Event, error) {
	panic("unimplemented")
}

// ListMonth implements storage.Storage.
func (s *SqlStorage) ListMonth(start time.Time) ([]storage.Event, error) {
	panic("unimplemented")
}

// ListWeek implements storage.Storage.
func (s *SqlStorage) ListWeek(start time.Time) ([]storage.Event, error) {
	panic("unimplemented")
}

// Update implements storage.Storage.
func (s *SqlStorage) Update(id string, eventRef *storage.Event) error {
	panic("unimplemented")
}
