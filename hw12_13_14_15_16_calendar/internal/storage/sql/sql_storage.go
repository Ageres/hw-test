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

func (s *SqlStorage) Add(ctx context.Context, eventRef *storage.Event) (*storage.Event, error) {
	if err := storage.ValidateEvent(eventRef); err != nil {
		return nil, err
	}

	var eventID string
	serr := storage.StorageError{}
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
		&serr.StatusCode,
		&serr.ErrorMessage,
		&serr.ConflictEventId,
	)

	if err != nil {
		serr.Message = fmt.Sprintf("add failed: %v", err)
		serr.Cause = err
		return nil, &serr
	}

	switch serr.StatusCode {
	case 200:
		savedEvent := *eventRef
		savedEvent.ID = eventID
		savedEvent.Duration = time.Duration(savedEvent.Duration) * time.Second
		savedEvent.Reminder = time.Duration(savedEvent.Reminder) * time.Second
		return &savedEvent, nil
	case 409:
		serr.Message = fmt.Sprintf(storage.ErrDateBusyMsgTemplate, serr.ConflictEventId)
		return nil, &serr
	case 504:
		serr.Message = fmt.Sprintf(storage.ErrDatabaseTimeoutMsgTemplate, serr.ErrorMessage)
		return nil, &serr
	default:
		serr.Message = fmt.Sprintf(storage.ErrDatabaseMsgTemplate, serr.ErrorMessage)
		return nil, &serr
	}
}

func (s *SqlStorage) Update(ctx context.Context, eventRef *storage.Event) error {
	if err := storage.FullValidateEvent(eventRef); err != nil {
		return err
	}

	var statusCode int
	serr := storage.StorageError{}
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
		&serr.StatusCode,
		&serr.ErrorMessage,
		&serr.ConflictEventId,
		&serr.ConflictUserId,
	)

	if err != nil {
		serr.Message = fmt.Sprintf("update failed: %v", err)
		serr.Cause = err
		return &serr
	}

	switch statusCode {
	case 200:
		return nil
	case 403:
		//return storage.ErrUserConflict
		serr.Message = fmt.Sprintf(storage.ErrUserConflictMsgTemplate, eventRef.UserID, serr.ConflictUserId)
		return &serr
	case 404:
		serr.Message = storage.ErrEventNotFoundMsg
		return &serr
	case 409:
		serr.Message = fmt.Sprintf(storage.ErrDateBusyMsgTemplate, serr.ConflictEventId)
		return &serr
	case 504:
		serr.Message = fmt.Sprintf(storage.ErrDatabaseTimeoutMsgTemplate, serr.ErrorMessage)
		return &serr
	default:
		serr.Message = fmt.Sprintf(storage.ErrDatabaseMsgTemplate, serr.ErrorMessage)
		return &serr
	}
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

func (p *SqlStorage) listEvents0(ctx context.Context, start, end time.Time) ([]storage.Event, error) {
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

func (p *SqlStorage) listEvents(ctx context.Context, start, end time.Time) ([]storage.Event, error) {
	//размер пачки для выгрузки
	batchSize := 1000

	// общее количество событий
	var totalEvents int
	err := p.db.GetContext(ctx, &totalEvents, `
        SELECT COUNT(*) FROM events 
        WHERE tstzrange(start_time, start_time + (duration * INTERVAL '1 second')) 
        && tstzrange($1::timestamptz, $2::timestamptz)`,
		start, end)
	if err != nil {
		return nil, fmt.Errorf("failed to count events: %w", err)
	}

	result := make([]storage.Event, 0, totalEvents)

	// выгрузка данных пачками
	offset := 0
	for {
		var batch []struct {
			ID          string
			Title       string
			StartTime   time.Time `db:"start_time"`
			Duration    int64
			Description string
			UserID      string `db:"user_id"`
			Reminder    int64
		}

		err := p.db.SelectContext(ctx, &batch, `
            SELECT id, title, start_time, duration, description, user_id, reminder
            FROM events 
            WHERE tstzrange(start_time, start_time + (duration * INTERVAL '1 second')) 
            && tstzrange($1::timestamptz, $2::timestamptz)
            ORDER BY start_time, id
            LIMIT $3 OFFSET $4`,
			start, end, batchSize, offset)

		if err != nil {
			return nil, fmt.Errorf("failed to fetch events batch: %w", err)
		}

		// Обрабатываем полученную пачку
		for _, e := range batch {
			result = append(result, storage.Event{
				ID:          e.ID,
				Title:       e.Title,
				StartTime:   e.StartTime,
				Duration:    time.Duration(e.Duration) * time.Second,
				Description: e.Description,
				UserID:      e.UserID,
				Reminder:    time.Duration(e.Reminder) * time.Second,
			})
		}

		// Проверяем, нужно ли продолжать
		if len(batch) < batchSize {
			break
		}
		offset += batchSize
	}

	return result, nil
}
