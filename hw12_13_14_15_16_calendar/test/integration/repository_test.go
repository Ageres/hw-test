////go:build integration

package integration

/*
import (
	"context"
	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/stretchr/testify/suite"
	"integration_testing/internal/domain"
	"integration_testing/internal/repository/postgres"
	"log"
	"os"
	"testing"
	"time"
)

type MyNewIntegrationSuite struct {
	suite.Suite
	pool *pgxpool.Pool
	r    Repo
}

func TestMyNewIntegrationSuite(t *testing.T) {
	suite.Run(t, new(MyNewIntegrationSuite))
}

func (s *MyNewIntegrationSuite) SetupSuite() {
	pool, err := pgxpool.Connect(context.Background(), "postgres://postgres:changeme@localhost:5432/postgres")
	if err != nil {
		log.Fatal(err)
	}
	s.pool = pool
}

func (s *MyNewIntegrationSuite) SetupTest() {
	s.r = postgres.NewRepo(s.pool)
	//s.r = fake.NewRepo(s.pool)
}

func (s *MyNewIntegrationSuite) TearDownTest() {
	_, _ = s.pool.Exec(context.Background(), "TRUNCATE items CASCADE")
}

func (s *MyNewIntegrationSuite) TestSave() {
	startTime := time.Now()
	item := domain.Item{
		Name:        "test",
		Description: "description",
		CreatedAt:   startTime,
		UpdatedAt:   startTime,
	}
	id, err := s.r.Save(context.Background(), item)
	s.Require().NoError(err)
	s.Require().NotEqual(0, id)

	dbItem := s.getDirectItem("test")

	s.Require().Equal(id, dbItem.ID)
	s.Require().Equal(item.Name, dbItem.Name)
	s.Require().Equal(item.Description, dbItem.Description)
	s.Require().Equal(item.CreatedAt.Format(time.RFC1123), dbItem.CreatedAt.Format(time.RFC1123))
	s.Require().Equal(item.UpdatedAt.Format(time.RFC1123), dbItem.UpdatedAt.Format(time.RFC1123))
}

func (s *MyNewIntegrationSuite) TestSaveDuplicateWithFixtures() {
	data, err := os.ReadFile("fixtures/orders.sql")
	if err != nil {
		s.Fail("could not read fixtures.")
	}
	_, err = s.pool.Exec(context.Background(), string(data))
	if err != nil {
		s.Fail("could not execute fixtures.")
	}

	startTime := time.Now()
	item := domain.Item{
		Name:        "item_1",
		Description: "description",
		CreatedAt:   startTime,
		UpdatedAt:   startTime,
	}
	_, err = s.r.Save(context.Background(), item)
	s.Require().Error(err)
}

func (s *MyNewIntegrationSuite) TestGetExist() {
	data, err := os.ReadFile("fixtures/orders.sql")
	if err != nil {
		s.Fail("could not read fixtures.")
	}
	_, err = s.pool.Exec(context.Background(), string(data))
	if err != nil {
		s.Fail("could not execute fixtures.")
	}

	item, err := s.r.Get(context.Background(), "item_1")
	s.Require().NoError(err)
	s.Require().NotEqual(0, item.ID)
}

func (s *MyNewIntegrationSuite) TestGet() {
	startTime := time.Now()
	existItem := domain.Item{
		Name:        "test",
		Description: "description",
		CreatedAt:   startTime,
		UpdatedAt:   startTime,
	}
	id := s.saveDirectItem(existItem)

	item, err := s.r.Get(context.Background(), "test")
	s.Require().NoError(err)
	s.Require().NotEqual(0, item.ID)

	s.Require().Equal(id, item.ID)
	s.Require().Equal(item.Name, existItem.Name)
	s.Require().Equal(item.Description, existItem.Description)
	s.Require().Equal(item.CreatedAt.Format(time.RFC1123), existItem.CreatedAt.Format(time.RFC1123))
	s.Require().Equal(item.UpdatedAt.Format(time.RFC1123), existItem.UpdatedAt.Format(time.RFC1123))
}

func (s *MyNewIntegrationSuite) TestGetNotFound() {
	_, err := s.r.Get(context.Background(), "test")
	s.Require().Error(err)
	s.Require().ErrorIs(err, pgx.ErrNoRows)
}

func (s *MyNewIntegrationSuite) saveDirectItem(item domain.Item) uint64 {
	query, args, err := sq.
		Insert("items").
		Columns("name", "description", "created_at", "updated_at").
		Values(
			item.Name,
			item.Description,
			time.Now().Format(time.RFC3339),
			time.Now().Format(time.RFC3339),
		).
		Suffix("RETURNING id").
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		s.Fail(err.Error())
	}

	rows, err := s.pool.Query(context.Background(), query, args...)
	if err != nil {
		s.Fail(err.Error())
	}
	defer rows.Close()

	var itemID uint64
	for rows.Next() {
		if scanErr := rows.Scan(&itemID); scanErr != nil {
			s.Fail(scanErr.Error())
		}
	}

	return itemID
}

func (s *MyNewIntegrationSuite) getDirectItem(name string) domain.Item {
	query, args, err := sq.
		Select("id", "name", "description", "created_at", "updated_at").
		From("items").
		Where(sq.Eq{"name": name}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		s.Fail(err.Error())
	}

	rows, err := s.pool.Query(context.Background(), query, args...)
	if err != nil {
		s.Fail(err.Error())
	}
	defer rows.Close()
	item := domain.Item{}

	for rows.Next() {
		scanErr := rows.Scan(&item.ID, &item.Name, &item.Description, &item.CreatedAt, &item.UpdatedAt)
		if scanErr != nil {
			s.Fail(scanErr.Error())
		}
	}

	return item
}
*/
