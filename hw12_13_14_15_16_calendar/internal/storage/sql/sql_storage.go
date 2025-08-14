package sqlstorage

import (
	"context"
	"fmt"
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
func (s *SqlStorage) Add(ctx context.Context, eventRef *storage.Event) error {
	panic("unimplemented")
}

// Delete implements storage.Storage.
func (s *SqlStorage) Delete(ctx context.Context, id string) error {
	panic("unimplemented")
}

// ListDay implements storage.Storage.
func (s *SqlStorage) ListDay(ctx context.Context, start time.Time) ([]storage.Event, error) {
	panic("unimplemented")
}

// ListMonth implements storage.Storage.
func (s *SqlStorage) ListMonth(ctx context.Context, start time.Time) ([]storage.Event, error) {
	panic("unimplemented")
}

// ListWeek implements storage.Storage.
func (s *SqlStorage) ListWeek(ctx context.Context, start time.Time) ([]storage.Event, error) {
	panic("unimplemented")
}

// Update implements storage.Storage.
func (s *SqlStorage) Update(ctx context.Context, eventRef *storage.Event) error {
	panic("unimplemented")
}

func (p *SqlStorage) listEvents(ctx context.Context, start, end time.Time) ([]storage.Event, error) {
	var events []struct {
		ID          string
		Title       string
		StartTime   time.Time `db:"start_time"`
		Duration    int64
		Description string
		UserID      string `db:"user_id"`
		Reminder    int64
	}

	err := p.db.Select(&events, `
        SELECT * FROM events 
        WHERE tstzrange(start_time, start_time + (duration * INTERVAL '1 second')) 
        && 
        tstzrange($1::timestamptz, $2::timestamptz)`,
		start,
		end,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to list events: %w", err)
	}

	result := make([]storage.Event, len(events))
	for i, e := range events {
		result[i] = storage.Event{
			ID:          e.ID,
			Title:       e.Title,
			StartTime:   e.StartTime,
			Duration:    time.Duration(e.Duration) * time.Second,
			Description: e.Description,
			UserID:      e.UserID,
			Reminder:    time.Duration(e.Reminder) * time.Second,
		}
	}
	return result, nil
}
