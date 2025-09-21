package repo

import (
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/Ageres/hw-test/hw12_13_14_15_calendar/test/integration/model"
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
	Delete(eventId string) error
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
	host, isSet := os.LookupEnv("DB_HOST")
	if !isSet {
		host = "localhost"
		log.Println("not found calendar db host, set default 'localhost'")
	}
	port, isSet := os.LookupEnv("DB_PORT")
	if !isSet {
		port = "5432"
		log.Println("not found calendar db port, set default '5432'")
	}
	name, isSet := os.LookupEnv("DB_NAME")
	if !isSet {
		name = "calendar"
		log.Println("not found calendar db name, set default 'calendar'")
	}
	user, isSet := os.LookupEnv("DB_USER")
	if !isSet {
		user = "user"
		log.Println("not found calendar db user, set default '5432'")
	}
	password, isSet := os.LookupEnv("DB_PASSWORD")
	if !isSet {
		password = "password"
		log.Println("not found calendar db password, set default 'password'")
	}
	return fmt.Sprintf(
		"host=%s port=%s dbname=%s user=%s password=%s sslmode=disable",
		host, port, name, user, password,
	)
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
			StartTime:   e.StartTime,
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
