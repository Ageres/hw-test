package integration

import (
	"context"
	"log"
	"testing"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/stretchr/testify/suite"
)

type CalendarIntegrationSuite struct {
	suite.Suite
	restApiClient TestCalendarApiClient
	pool          *pgxpool.Pool
	repo          Repo
}

func (s *CalendarIntegrationSuite) SetupSuite() {
	pool, err := pgxpool.Connect(context.Background(), "postgres://user:password@localhost:5432/calendar")
	if err != nil {
		log.Fatal(err)
	}
	s.pool = pool

	s.repo = NewRepo()

	s.restApiClient = newRestapiClient()
}

func NewSuite() *CalendarIntegrationSuite {
	return &CalendarIntegrationSuite{}
}

func TestRepoSuite(t *testing.T) {
	suite.Run(t, NewSuite())
}

func (s *CalendarIntegrationSuite) TestAddEventByRestApi() {
	timeLocation, err := time.LoadLocation("Local")
	s.Require().NoError(err)
	startTime := time.Date(2030, 12, 31, 10, 0, 0, 0, timeLocation)
	restApiEvent := &TestEvent{
		Title:       "title TestAddEventByRestApi",
		StartTime:   startTime,
		Duration:    24 * time.Hour,
		Description: "description TestAddEventByRestApi",
		UserID:      "user-id-TestAddEventByRestApi",
		Reminder:    24 * time.Hour,
	}

	eventId, err := s.restApiClient.AddTestEvent(restApiEvent)
	s.Require().NoError(err)
	s.Require().NotEqual("", eventId)
	restApiEvent.ID = eventId

	dbEvent, err := s.repo.Get(eventId)
	s.Require().NoError(err)
	s.Require().Equal(restApiEvent, dbEvent)

	err = s.repo.Delete(eventId)
	s.Require().NoError(err)
}
