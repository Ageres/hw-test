package integration

import (
	"context"
	"fmt"
	"integration_testing/internal/domain"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v4/pgxpool"
)

const (
	eventsTable = "events"
)

type Repo interface {
	//Save(ctx context.Context, testEventRef *TestEvent) (string, error)
	Get(ctx context.Context, eventId string) (*TestEvent, error)
	Delete(ctx context.Context, eventId string) error
}

type repo struct {
	db *pgxpool.Pool
}

// Get implements Repo.
func (r *repo) Get(ctx context.Context, eventId string) (*TestEvent, error) {
	panic("unimplemented")
}

func NewRepo(db *pgxpool.Pool) Repo {
	return &repo{
		db: db,
	}
}

/*
func (r *repo) Save(ctx context.Context, event *TestEvent) (string, error) {
	tx, err := r.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return "", fmt.Errorf("can't create tx: %w", err)
	}

	defer func() {
		if err != nil {
			rollbackErr := tx.Rollback(ctx)
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
		return 0, fmt.Errorf("tx err: %w", err)
	}
	defer rows.Close()

	var itemID uint64
	for rows.Next() {
		if scanErr := rows.Scan(&itemID); scanErr != nil {
			return 0, fmt.Errorf("can't scan itemID: %w", scanErr)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return 0, fmt.Errorf("can't commit tx: %w", err)
	}

	return itemID, nil
}
*/

func (r *repo) Get(ctx context.Context, eventId string) (domain.Item, error) {

	// build
	query, args, err := sq.
		Select("id", "name", "description", "created_at", "updated_at").
		From(eventsTable).
		Where(sq.Eq{"name": eventId}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return domain.Item{}, fmt.Errorf("can't build query: %w", err)
	}

	// get
	item := domain.Item{}
	err = r.db.QueryRow(ctx, query, args...).Scan(&item.ID, &item.Name, &item.Description, &item.CreatedAt, &item.UpdatedAt)
	if err != nil {
		return domain.Item{}, fmt.Errorf("can't select orders: %w", err)
	}

	//for rows.Next() {
	//	scanErr := rows.Scan()
	//	if scanErr != nil {
	//		return domain.Item{}, fmt.Errorf("can't scan order: %w", scanErr)
	//	}
	//}

	return item, nil
}

// Delete implements Repo.
func (r *repo) Delete(ctx context.Context, eventId string) error {
	panic("unimplemented")
}
