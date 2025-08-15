package sqlstorage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/Ageres/hw-test/hw12_13_14_15_calendar/internal/model"
	"github.com/Ageres/hw-test/hw12_13_14_15_calendar/internal/storage"
	"github.com/jackc/pgx/v5/pgconn"
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

func (s *SqlStorage) Add(ctx context.Context, eventRef *storage.Event) (*storage.Event, error) {
	// Создаем копию события, чтобы не модифицировать исходный объект
	savedEvent := *eventRef

	err := s.db.QueryRowContext(ctx, `
        INSERT INTO events 
        (title, start_time, duration, description, user_id, reminder)
        VALUES ($1, $2, $3, $4, $5, $6)
        RETURNING id, title, start_time, duration, description, user_id, reminder`,
		eventRef.Title,
		eventRef.StartTime,
		int(eventRef.Duration.Seconds()),
		eventRef.Description,
		eventRef.UserID,
		int(eventRef.Reminder.Seconds()),
	).Scan(
		&savedEvent.ID,
		&savedEvent.Title,
		&savedEvent.StartTime,
		&savedEvent.Duration,
		&savedEvent.Description,
		&savedEvent.UserID,
		&savedEvent.Reminder,
	)

	if err != nil {
		if pqErr, ok := err.(*pgconn.PgError); ok && pqErr.Code == "23P01" {
			return nil, storage.ErrDateBusy
		}
		return nil, fmt.Errorf("failed to insert event: %w", err)
	}

	savedEvent.Duration = time.Duration(savedEvent.Duration) * time.Second
	savedEvent.Reminder = time.Duration(savedEvent.Reminder) * time.Second

	return &savedEvent, nil
}

func (s *SqlStorage) Update(ctx context.Context, eventRef *storage.Event) error {
	tx, err := s.db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	var currentUserID string
	err = tx.GetContext(ctx, &currentUserID,
		"SELECT user_id FROM events WHERE id = $1",
		eventRef.ID)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return storage.ErrEventNotFound
		}
		return fmt.Errorf("failed to check event owner: %w", err)
	}

	if currentUserID != eventRef.UserID {
		return storage.ErrUserConflict
	}

	_, err = tx.ExecContext(ctx, `
        UPDATE events SET
            title = $1,
            start_time = $2,
            duration = $3,
            description = $4,
            reminder = $5
        WHERE id = $6`,
		eventRef.Title,
		eventRef.StartTime,
		int(eventRef.Duration.Seconds()),
		eventRef.Description,
		int(eventRef.Reminder.Seconds()),
		eventRef.ID,
	)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch pgErr.Code {
			case "23P01":
				return storage.ErrDateBusy
			case "23514":
				return fmt.Errorf("invalid event data: %w", err)
			}
		}
		return fmt.Errorf("failed to update event: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (s *SqlStorage) Delete(ctx context.Context, id string) error {
	res, err := s.db.Exec("DELETE FROM events WHERE id = $1", id)
	if err != nil {
		return fmt.Errorf("failed to delete event: %w", err)
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected when deleting event: %w", err)
	} else if rows == 0 {
		return storage.ErrEventNotFound
	}
	return nil
}

func (s *SqlStorage) ListDay(ctx context.Context, startDay time.Time) ([]storage.Event, error) {
	start := s.getStartDayTime(startDay)
	end := start.AddDate(0, 0, 1)
	return s.listEvents(ctx, start, end)
}

func (s *SqlStorage) ListWeek(ctx context.Context, startDay time.Time) ([]storage.Event, error) {
	start := s.getStartDayTime(startDay)
	end := start.AddDate(0, 0, 7)
	return s.listEvents(ctx, start, end)
}

func (s *SqlStorage) ListMonth(ctx context.Context, startDay time.Time) ([]storage.Event, error) {
	start := s.getStartDayTime(startDay)
	end := start.AddDate(0, 1, 0)
	return s.listEvents(ctx, start, end)
}

func (p *SqlStorage) getStartDayTime(start time.Time) time.Time {
	return time.Date(start.Year(), start.Month(), start.Day(), 0, 0, 0, 0, start.Location())
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
