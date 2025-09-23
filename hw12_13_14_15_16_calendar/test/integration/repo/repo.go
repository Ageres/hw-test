package repo

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/Ageres/hw-test/hw12_13_14_15_calendar/test/integration/config"
	"github.com/Ageres/hw-test/hw12_13_14_15_calendar/test/integration/model"
	"github.com/Ageres/hw-test/hw12_13_14_15_calendar/test/integration/utils"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

type dbEvent struct {
	ID          string
	Title       string
	StartTime   time.Time `db:"start_time"`
	Duration    int64
	Description string
	UserID      string `db:"user_id"`
	Reminder    int64
}

type dbProcEvent struct {
	ID     string
	UserID string `db:"user_id"`
}

type Repo interface {
	Get(eventID string) (*model.TestEvent, error)
	ListByUserID(userID string) ([]model.TestEvent, error)
	Delete(eventID string) error
	DeleteByUserID(userID string) error
	CheckProcEvent(procEventID string) (bool, error)
	DeleteProcEventByUserID(userID string) error
}

const (
	reconnectAttempt = 6
	reconectTimeout  = 10
)

type repo struct {
	db *sqlx.DB
}

func NewRepo() Repo {
	dsn := dns()

	var db *sqlx.DB
	var err error
	for i := range reconnectAttempt {
		db, err = sqlx.Connect("pgx", dsn)
		if err != nil {
			log.Println(err)
			if i < reconnectAttempt-1 {
				db = nil
				err = nil
				time.Sleep(reconectTimeout * time.Second)
				continue
			}
		}
	}
	if err != nil {
		log.Fatal(err)
	}

	return &repo{db: db}
}

func dns() string {
	host := utils.GetEnvOrDefault(config.CalendarDBHostEnv, config.CalendarDBHostDefault)
	port := utils.GetEnvOrDefault(config.CalendarDBPortEnv, config.CalendarDBPortDefault)
	name := utils.GetEnvOrDefault(config.CalendarDBNameEnv, config.CalendarDBNameDefault)
	user := utils.GetEnvOrDefault(config.CalendarDBUserEnv, config.CalendarDBUserDefault)
	password := utils.GetEnvOrDefault(config.CalendarDBPasswordEnv, config.CalendarDBPasswordDefault)
	return fmt.Sprintf(
		"host=%s port=%s dbname=%s user=%s password=%s sslmode=disable",
		host, port, name, user, password,
	)
}

func (r *repo) Get(eventID string) (*model.TestEvent, error) {
	rows, err := r.db.Queryx(`
        SELECT id, title, start_time, duration, description, user_id, reminder 
        FROM events 
        WHERE id = $1 `,
		eventID)
	if err != nil {
		return nil, fmt.Errorf("can't select event: %w", err)
	}
	defer rows.Close()
	result := make([]model.TestEvent, 0, 1)
	for rows.Next() {
		var e dbEvent
		if err := rows.StructScan(&e); err != nil {
			return nil, fmt.Errorf("failed to scan event: %w", err)
		}
		result = append(result, model.TestEvent{
			ID:          e.ID,
			Title:       e.Title,
			StartTime:   e.StartTime.Local(),
			Duration:    time.Duration(e.Duration) * time.Second,
			Description: e.Description,
			UserID:      e.UserID,
			Reminder:    time.Duration(e.Reminder) * time.Second,
		})
	}
	if len(result) > 1 {
		return nil, fmt.Errorf("found more than one event for id '%s', len '%d'", eventID, len(result))
	}
	return &result[0], nil
}

func (r *repo) ListByUserID(userID string) ([]model.TestEvent, error) {
	rows, err := r.db.Queryx(`
        SELECT id, title, start_time, duration, description, user_id, reminder 
        FROM events 
        WHERE user_id = $1 `,
		userID)
	if err != nil {
		return nil, fmt.Errorf("can't select event: %w", err)
	}
	defer rows.Close()
	result := make([]model.TestEvent, 0, 1)
	for rows.Next() {
		var e dbEvent
		if err := rows.StructScan(&e); err != nil {
			return nil, fmt.Errorf("failed to scan event: %w", err)
		}
		result = append(result, model.TestEvent{
			ID:          e.ID,
			Title:       e.Title,
			StartTime:   e.StartTime.Local(),
			Duration:    time.Duration(e.Duration) * time.Second,
			Description: e.Description,
			UserID:      e.UserID,
			Reminder:    time.Duration(e.Reminder) * time.Second,
		})
	}
	if len(result) == 0 {
		return nil, fmt.Errorf("not found events for user_id '%s'", userID)
	}
	return result, nil
}

func (r *repo) Delete(eventID string) error {
	res, err := r.db.Exec("DELETE FROM events WHERE id = $1", eventID)
	if err != nil {
		return fmt.Errorf("delete event: %s", err.Error())
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("can't delete event: %w", err)
	} else if rows == 0 {
		return errors.New("not found event for deleting")
	}
	return nil
}

func (r *repo) DeleteByUserID(userID string) error {
	res, err := r.db.Exec("DELETE FROM events WHERE user_id = $1", userID)
	if err != nil {
		return fmt.Errorf("delete event: %s", err.Error())
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("can't delete event: %w", err)
	} else if rows == 0 {
		return errors.New("not found event for deleting")
	}
	return nil
}

func (r *repo) CheckProcEvent(procEventID string) (bool, error) {
	rows, err := r.db.Queryx(`SELECT id FROM proc_events WHERE id = $1 `, procEventID)
	if err != nil {
		return false, fmt.Errorf("can't select event: %w", err)
	}
	defer rows.Close()
	result := make([]dbProcEvent, 0, 1)
	for rows.Next() {
		var p dbProcEvent
		if err := rows.StructScan(&p); err != nil {
			return false, fmt.Errorf("failed to scan proc event: %w", err)
		}
		result = append(result, p)
	}
	if len(result) > 1 {
		return false, fmt.Errorf("found more than one proc event for id '%s', len '%d'", procEventID, len(result))
	}
	if len(result) == 1 {
		return true, nil
	}
	return false, nil
}

func (r *repo) DeleteProcEventByUserID(userID string) error {
	res, err := r.db.Exec("DELETE FROM proc_events WHERE user_id = $1", userID)
	if err != nil {
		return fmt.Errorf("delete proc event: %s", err.Error())
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("can't delete proc event: %w", err)
	} else if rows == 0 {
		return errors.New("not found proc event for deleting")
	}
	return nil
}
