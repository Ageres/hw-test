////go:build integration

package integration

/*
import (
	"context"
	"integration_testing/internal/domain"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/stretchr/testify/suite"

	//myRepo "integration_testing/internal/repository/fake"
	"log"
	"testing"
	"time"
)

type Repo interface {
	Save(ctx context.Context, item domain.Item) (uint64, error)
	Get(ctx context.Context, name string) (domain.Item, error)
}

type MyFirstIntegrationSuite struct {
	suite.Suite
	pool *pgxpool.Pool
	r    Repo
}

func NewRepoSuite() *MyFirstIntegrationSuite {
	return &MyFirstIntegrationSuite{}
}

type CalendarIntegrationSuite struct {
	suite.Suite
	pool *pgxpool.Pool
	r    Repo
}

func (s *MyFirstIntegrationSuite) SetupSuite() {
	pool, err := pgxpool.Connect(context.Background(), "postgres://postgres:changeme@localhost:5432/postgres")
	if err != nil {
		log.Fatal(err)
	}
	s.pool = pool
}

func (s *MyFirstIntegrationSuite) SetupTest() {
	s.r = myRepo.NewRepo(s.pool)
}

// Сохранить товар
func (s *MyFirstIntegrationSuite) TestSaveItem() {
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
}

// Получить товар
func (s *MyFirstIntegrationSuite) TestGetItem() {
	startTime := time.Now()
	item := domain.Item{
		Name:        "test",
		Description: "description",
		CreatedAt:   startTime,
		UpdatedAt:   startTime,
	}
	id, err := s.r.Save(context.Background(), item)
	s.Require().NoError(err)

	location, err := time.LoadLocation("Europe/Moscow")
	s.Require().NoError(err)

	dbEvent, err := s.r.Get(context.Background(), "test")
	s.Require().NoError(err)
	s.Require().Equal(startTime.Format(time.RFC1123), dbEvent.CreatedAt.In(location).Format(time.RFC1123))
	s.Require().Equal(id, dbEvent.ID)
	s.Require().Equal("description", dbEvent.Description)
}

func TestRepoSuite(t *testing.T) {
	suite.Run(t, NewRepoSuite())
}
*/
