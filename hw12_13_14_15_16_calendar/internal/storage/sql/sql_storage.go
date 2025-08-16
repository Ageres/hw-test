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

var ErrDatabaseMsgTemplate = "database error: %s"

type SqlStorage struct {
	db *sqlx.DB
}

type dbResp struct {
	statusCode      int
	errorMessage    string
	conflictEventId string
	conflictUserId  string
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
	if err := storage.ValidateEvent(eventRef); err != nil {
		return nil, err
	}

	var eventID string
	dbResp := dbResp{}
	err := s.db.QueryRowContext(ctx, `
        SELECT event_id, status_code, error_message, conflict_event_id 
        FROM add_event($1, $2, $3, $4, $5, $6)`,
		eventRef.Title,
		eventRef.StartTime,
		int(eventRef.Duration.Seconds()),
		eventRef.Description,
		eventRef.UserID,
		int(eventRef.Reminder.Seconds()),
	).Scan(
		&eventID,
		&dbResp.statusCode,
		&dbResp.errorMessage,
		&dbResp.conflictEventId,
	)

	if err != nil {
		return nil, storage.NewSErrorWithCause(ErrDatabaseMsgTemplate, err)
	}

	switch dbResp.statusCode {
	case 200:
		savedEvent := *eventRef
		savedEvent.ID = eventID
		savedEvent.Duration = time.Duration(savedEvent.Duration) * time.Second
		savedEvent.Reminder = time.Duration(savedEvent.Reminder) * time.Second
		return &savedEvent, nil
	case 409:
		return nil, storage.NewSErrorWithTemplate(storage.ErrDateBusyMsgTemplate, dbResp.conflictEventId)
	default:
		return nil, storage.NewSErrorWithTemplate(ErrDatabaseMsgTemplate, dbResp.errorMessage)
	}
}

func (s *SqlStorage) Update(ctx context.Context, eventRef *storage.Event) error {
	if err := storage.FullValidateEvent(eventRef); err != nil {
		return err
	}

	var statusCode int
	dbResp := dbResp{}
	err := s.db.QueryRowContext(ctx, `
        SELECT status_code, error_message, conflict_event_id, conflict_user_id 
        FROM update_event($1, $2, $3, $4, $5, $6, $7)`,
		eventRef.ID,
		eventRef.Title,
		eventRef.StartTime,
		int(eventRef.Duration.Seconds()),
		eventRef.Description,
		eventRef.UserID,
		int(eventRef.Reminder.Seconds()),
	).Scan(
		&dbResp.statusCode,
		&dbResp.errorMessage,
		&dbResp.conflictEventId,
		&dbResp.conflictUserId,
	)

	if err != nil {
		return storage.NewSErrorWithCause(ErrDatabaseMsgTemplate, err)
	}

	switch statusCode {
	case 200:
		return nil
	case 403:
		return storage.NewSErrorWithTemplate(storage.ErrUserConflictMsgTemplate, dbResp.conflictUserId)
	case 404:
		return storage.ErrEventNotFound
	case 409:
		return storage.NewSErrorWithTemplate(storage.ErrDateBusyMsgTemplate, dbResp.conflictEventId)
	default:
		return storage.NewSErrorWithTemplate(ErrDatabaseMsgTemplate, dbResp.errorMessage)
	}
}

func (s *SqlStorage) Delete(ctx context.Context, id string) error {
	res, err := s.db.Exec("DELETE FROM events WHERE id = $1", id)
	if err != nil {
		return storage.NewSErrorWithCause(ErrDatabaseMsgTemplate, err)
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return storage.NewSErrorWithCause("failed to get rows affected when deleting event: %v", err)
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
	result := make([]storage.Event, 0, 100)
	rows, err := p.db.QueryxContext(ctx, `
        SELECT id, title, start_time, duration, description, user_id, reminder 
        FROM events 
        WHERE tstzrange(start_time, start_time + (duration * INTERVAL '1 second')) 
        && tstzrange($1::timestamptz, $2::timestamptz)`,
		start, end)
	if err != nil {
		return nil, storage.NewSErrorWithCause(ErrDatabaseMsgTemplate, err)
	}
	defer rows.Close()

	for rows.Next() {
		var e struct {
			ID          string
			Title       string
			StartTime   time.Time `db:"start_time"`
			Duration    int64
			Description string
			UserID      string `db:"user_id"`
			Reminder    int64
		}

		if err := rows.StructScan(&e); err != nil {
			return nil, storage.NewSErrorWithCause("failed to scan event: %v", err)
		}

		result = append(result, storage.Event{
			ID:          e.ID,
			Title:       e.Title,
			StartTime:   e.StartTime,
			Duration:    time.Duration(e.Duration) * time.Second,
			Description: e.Description,
			UserID:      e.UserID,
			Reminder:    time.Duration(e.Reminder) * time.Second,
		})

		select {
		case <-ctx.Done():
			return nil, storage.NewSErrorWithCause(storage.ErrContextDoneTemplate, ctx.Err())
		default:
		}

	}

	if err := rows.Err(); err != nil {
		return nil, storage.NewSErrorWithCause("rows iteration error: %v", err)
	}

	return result, nil
}
