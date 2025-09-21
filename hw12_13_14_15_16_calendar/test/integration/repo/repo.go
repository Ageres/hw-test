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

type Repo interface {
	Get(eventId string) (*model.TestEvent, error)
	ListByUserId(userId string) ([]model.TestEvent, error)
	Delete(eventId string) error
	DeleteByUserId(userId string) error
	CheckProcEventId(procEventId string) (bool, error)
}

type repo struct {
	db *sqlx.DB
}

func NewRepo() Repo {
	dsn := dns()
	db, err := sqlx.Connect("pgx", dsn)
	if err != nil {
		log.Fatal(err)
	}

	//db.SetMaxOpenConns(sqlConfRef.Pool.Conn.MaxOpen)
	//db.SetMaxIdleConns(sqlConfRef.Pool.Conn.MaxIdle)
	//db.SetConnMaxLifetime(time.Duration(sqlConfRef.Pool.Conn.MaxLifeTime) * time.Second)
	//db.SetConnMaxIdleTime(time.Duration(sqlConfRef.Pool.Conn.MaxLifeTime) * time.Second)

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	return &repo{db: db}
}

func dns() string {
	host := utils.GetEnvOrDefault(config.CALENDAR_DB_HOST_ENV, config.CALENDAR_DB_HOST_DEFAULT)
	port := utils.GetEnvOrDefault(config.CALENDAR_DB_PORT_ENV, config.CALENDAR_DB_PORT_DEFAULT)
	name := utils.GetEnvOrDefault(config.CALENDAR_DB_NAME_ENV, config.CALENDAR_DB_NAME_DEFAULT)
	user := utils.GetEnvOrDefault(config.CALENDAR_DB_USER_ENV, config.CALENDAR_DB_USER_DEFAULT)
	password := utils.GetEnvOrDefault(config.CALENDAR_DB_PASSWORD_ENV, config.CALENDAR_DB_PASSWORD_DEFAULT)
	return fmt.Sprintf(
		"host=%s port=%s dbname=%s user=%s password=%s sslmode=disable",
		host, port, name, user, password,
	)
}

func (r *repo) Get(eventId string) (*model.TestEvent, error) {
	rows, err := r.db.Queryx(`
        SELECT id, title, start_time, duration, description, user_id, reminder 
        FROM events 
        WHERE id = $1 `,
		eventId)
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
		return nil, fmt.Errorf("found more than one event for id '%s', len '%d'", eventId, len(result))
	}
	return &result[0], nil
}

func (r *repo) ListByUserId(userId string) ([]model.TestEvent, error) {
	rows, err := r.db.Queryx(`
        SELECT id, title, start_time, duration, description, user_id, reminder 
        FROM events 
        WHERE user_id = $1 `,
		userId)
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
		return nil, fmt.Errorf("not found events for user_id '%s'", userId)
	}
	return result, nil
}

func (r *repo) Delete(eventId string) error {
	res, err := r.db.Exec("DELETE FROM events WHERE id = $1", eventId)
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

func (r *repo) DeleteByUserId(userId string) error {
	res, err := r.db.Exec("DELETE FROM events WHERE user_id = $1", userId)
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

func (r *repo) CheckProcEventId(procEventId string) (bool, error) {
	rows, err := r.db.Queryx(`SELECT id FROM proc_events WHERE id = $1 `, procEventId)
	if err != nil {
		return false, fmt.Errorf("can't select event: %w", err)
	}
	defer rows.Close()
	result := make([]string, 0, 1)
	for rows.Next() {
		var id string
		if err := rows.StructScan(&id); err != nil {
			return false, fmt.Errorf("failed to scan proc event: %w", err)
		}
		result = append(result, id)
	}
	if len(result) > 1 {
		return false, fmt.Errorf("found more than one proc event for id '%s', len '%d'", procEventId, len(result))
	}
	if len(result) == 1 {
		return true, nil
	}
	return false, nil
}
