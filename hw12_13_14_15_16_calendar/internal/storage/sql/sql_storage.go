package sqlstorage

import (
	"context"
	"os"
	"time"

	lg "github.com/Ageres/hw-test/hw12_13_14_15_calendar/internal/logger"
	"github.com/Ageres/hw-test/hw12_13_14_15_calendar/internal/model"
	"github.com/Ageres/hw-test/hw12_13_14_15_calendar/internal/storage"
	// регистрация драйвера PostgreSQL.
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

var ErrDatabaseMsg = "database error"

type SQLStorage struct {
	db *sqlx.DB
}

type dbResp struct {
	statusCode      int
	errorMessage    string
	conflictEventID string
	conflictUserID  string
}

func NewSQLStorage(ctx context.Context, storageConfRef *model.StorageConf) storage.Storage {
	logger := lg.GetLogger(ctx)

	sqlConfRef := storageConfRef.SQL
	dsn := sqlConfRef.DB.DSN()
	db, err := sqlx.Connect("pgx", dsn)
	if err != nil {
		logger.WithError(err).Error("failed to load driver")
		os.Exit(1)
	}

	db.SetMaxOpenConns(sqlConfRef.Pool.Conn.MaxOpen)
	db.SetMaxIdleConns(sqlConfRef.Pool.Conn.MaxIdle)
	db.SetConnMaxLifetime(time.Duration(sqlConfRef.Pool.Conn.MaxLifeTime) * time.Second)
	db.SetConnMaxIdleTime(time.Duration(sqlConfRef.Pool.Conn.MaxLifeTime) * time.Second)

	err = db.Ping()
	if err != nil {
		logger.WithError(err).Error("failed to connect to db")
		os.Exit(1)
	}

	return &SQLStorage{
		db: db,
	}
}

func (s *SQLStorage) Close() error {
	return s.db.Close()
}

func (s *SQLStorage) Add(ctx context.Context, eventRef *storage.Event) (*storage.Event, error) {
	logger := lg.GetLogger(ctx)
	logger.Info("add event", map[string]any{"event": lg.MarshalAny(eventRef)})

	if err := storage.ValidateEvent(eventRef); err != nil {
		logger.WithError(err).Error("add event")
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
		&dbResp.conflictEventID,
	)
	if err != nil {
		err = storage.NewSError(ErrDatabaseMsg, err)
		logger.WithError(err).Error("add event", map[string]any{"error": err})
		return nil, err
	}

	switch dbResp.statusCode {
	case 200:
		savedEvent := *eventRef
		savedEvent.ID = eventID
		return &savedEvent, nil
	case 409:
		err := storage.NewSErrorWithTemplate(storage.ErrDateBusyMsgTemplate, dbResp.conflictEventID)
		logger.WithError(err).Error("add event", map[string]any{"databaseResponse": dbResp})
		return nil, err
	default:
		err := storage.NewSErrorWithTemplate(ErrDatabaseMsg, dbResp.errorMessage)
		logger.WithError(err).Error("add event", map[string]any{"databaseResponse": dbResp})
		return nil, err
	}
}

func (s *SQLStorage) Update(ctx context.Context, eventRef *storage.Event) error {
	logger := lg.GetLogger(ctx)
	logger.Info("update event", map[string]any{"event": lg.MarshalAny(eventRef)})

	if err := storage.FullValidateEvent(eventRef); err != nil {
		logger.WithError(err).Error("update event")
		return err
	}

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
		&dbResp.conflictEventID,
		&dbResp.conflictUserID,
	)
	if err != nil {
		err = storage.NewSError(ErrDatabaseMsg, err)
		logger.WithError(err).Error("update event")
		return err
	}

	switch dbResp.statusCode {
	case 200:
		return nil
	case 403:
		err := storage.NewSErrorWithTemplate(storage.ErrUserConflictMsgTemplate, eventRef.UserID, dbResp.conflictUserID)
		logger.WithError(err).Error("update event", map[string]any{"databaseResponse": dbResp})
		return err
	case 404:
		err := storage.ErrEventNotFound
		logger.WithError(err).Error("update event", map[string]any{"databaseResponse": dbResp})
		return err
	case 409:
		err := storage.NewSErrorWithTemplate(storage.ErrDateBusyMsgTemplate, dbResp.conflictEventID)
		logger.WithError(err).Error("update event", map[string]any{"databaseResponse": dbResp})
		return err
	default:
		err := storage.NewSErrorWithTemplate(ErrDatabaseMsg, dbResp.errorMessage)
		logger.WithError(err).Error("add event")
		return err
	}
}

func (s *SQLStorage) Delete(ctx context.Context, id string) error {
	logger := lg.GetLogger(ctx)
	logger.Info("delete event", map[string]any{"eventId": id})

	res, err := s.db.Exec("DELETE FROM events WHERE id = $1", id)
	if err != nil {
		err = storage.NewSError(ErrDatabaseMsg, err)
		logger.WithError(err).Error("delete event")
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		err = storage.NewSError("failed to get rows affected when deleting event", err)
		logger.WithError(err).Error("delete event")
		return err
	} else if rows == 0 {
		err := storage.ErrEventNotFound
		logger.WithError(err).Error("delete event")
		return err
	}
	return nil
}

func (s *SQLStorage) ListDay(ctx context.Context, startDay time.Time) ([]storage.Event, error) {
	logger := lg.GetLogger(ctx)
	logger.Info("list day events", map[string]any{"startDay": startDay})

	start := s.getStartDayTime(startDay)
	end := start.AddDate(0, 0, 1)
	return s.listEvents(ctx, start, end)
}

func (s *SQLStorage) ListWeek(ctx context.Context, startDay time.Time) ([]storage.Event, error) {
	logger := lg.GetLogger(ctx)
	logger.Info("list week events", map[string]any{"startDay": startDay})

	start := s.getStartDayTime(startDay)
	end := start.AddDate(0, 0, 7)
	return s.listEvents(ctx, start, end)
}

func (s *SQLStorage) ListMonth(ctx context.Context, startDay time.Time) ([]storage.Event, error) {
	logger := lg.GetLogger(ctx)
	logger.Info("list month events", map[string]any{"startDay": startDay})

	start := s.getStartDayTime(startDay)
	end := start.AddDate(0, 1, 0)
	return s.listEvents(ctx, start, end)
}

func (s *SQLStorage) getStartDayTime(start time.Time) time.Time {
	return time.Date(start.Year(), start.Month(), start.Day(), 0, 0, 0, 0, start.Location())
}

func (s *SQLStorage) listEvents(ctx context.Context, startTime, endTime time.Time) ([]storage.Event, error) {
	logger := lg.GetLogger(ctx)
	logger.Info("list events", map[string]any{
		"startTime": startTime,
		"endTime":   endTime,
	})

	result := make([]storage.Event, 0, 100)
	rows, err := s.db.QueryxContext(ctx, `
        SELECT id, title, start_time, duration, description, user_id, reminder 
        FROM events 
        WHERE tstzrange(start_time, start_time + (duration * INTERVAL '1 second')) 
        && tstzrange($1::timestamptz, $2::timestamptz)`,
		startTime, endTime)
	if err != nil {
		err = storage.NewSError(ErrDatabaseMsg, err)
		logger.WithError(err).Error("list events")
		return nil, err
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
			err = storage.NewSError("failed to scan event", err)
			logger.WithError(err).Error("list events")
			return nil, err
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
			err = storage.NewSError(storage.ErrContextDone, ctx.Err())
			logger.WithError(err).Error("list events")
			return nil, err
		default:
		}
	}

	if err := rows.Err(); err != nil {
		err = storage.NewSError("rows iteration error", err)
		logger.WithError(err).Error("list events")
		return nil, err
	}
	logger.Info("list events", map[string]any{"found": len(result)})
	return result, nil
}
