package integration

import (
	"errors"
	"fmt"
	"log"
	"os"
	"time"

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
	//Add(ctx context.Context, testEventRef *TestEvent) (string, error)
	Get(eventId string) (*TestEvent, error)
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

func (r *repo) Get(eventId string) (*TestEvent, error) {
	rows, err := r.db.Queryx(`
        SELECT id, title, start_time, duration, description, user_id, reminder 
        FROM events 
        WHERE id = $1 `,
		eventId)
	if err != nil {
		return nil, fmt.Errorf("can't select event: %w", err)
	}
	defer rows.Close()
	result := make([]TestEvent, 0, 1)
	for rows.Next() {
		var e dbEvent
		if err := rows.StructScan(&e); err != nil {
			return nil, fmt.Errorf("failed to scan event: %w", err)
		}
		result = append(result, TestEvent{
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

/*
func (r *repo) Add(ctx context.Context, event *TestEvent) (string, error) {
	tx, err := r.db.BeginTxx(ctx, &sql.TxOptions{})
	if err != nil {
		return "", fmt.Errorf("can't create tx: %w", err)
	}

	defer func() {
		if err != nil {
			rollbackErr := tx.Rollback()
			if rollbackErr != nil {
				err = errors.Join(err, rollbackErr)
			}
		}
	}()

	query, args, err := sq.
		Insert(eventsTable).
		Columns("name", "description", "created_at", "updated_at").
		Values(
			event.Name,
			event.Description,
			time.Now().Format(time.RFC3339),
			time.Now().Format(time.RFC3339),
		).
		Suffix("RETURNING id").
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return 0, fmt.Errorf("can't build sql: %w", err)
	}

	rows, err := tx.Query(ctx, query, args...)
	if err != nil {
		return "", fmt.Errorf("tx err: %w", err)
	}
	defer rows.Close()

	var itemID string
	for rows.Next() {
		if scanErr := rows.Scan(&itemID); scanErr != nil {
			return "", fmt.Errorf("can't scan itemID: %w", scanErr)
		}
	}

	if err := tx.Commit(); err != nil {
		return "", fmt.Errorf("can't commit tx: %w", err)
	}

	return itemID, nil
}
*/
