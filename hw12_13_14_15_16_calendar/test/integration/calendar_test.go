package integration

import (
	"testing"
	"time"

	c "github.com/Ageres/hw-test/hw12_13_14_15_calendar/test/integration/client"
	ch "github.com/Ageres/hw-test/hw12_13_14_15_calendar/test/integration/client/http"
	"github.com/Ageres/hw-test/hw12_13_14_15_calendar/test/integration/model"
	"github.com/Ageres/hw-test/hw12_13_14_15_calendar/test/integration/repo"
	"github.com/stretchr/testify/suite"
)

type CalendarIntegrationSuite struct {
	suite.Suite
	restApiClient c.TestCalendarApiClient
	repo          repo.Repo
}

func (s *CalendarIntegrationSuite) SetupSuite() {
	s.repo = repo.NewRepo()
	s.restApiClient = ch.NewRestapiClient()
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
	restApiEvent := &model.TestEvent{
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

func (s *CalendarIntegrationSuite) TestBusyDateByRestApi() {
	timeLocation, err := time.LoadLocation("Local")
	s.Require().NoError(err)
	startTime := time.Date(2030, 12, 31, 10, 0, 0, 0, timeLocation)
	restApiEventOk := &model.TestEvent{
		Title:       "title TestAddEventByRestApi",
		StartTime:   startTime,
		Duration:    24 * time.Hour,
		Description: "description TestAddEventByRestApi",
		UserID:      "user-id-TestAddEventByRestApi",
		Reminder:    24 * time.Hour,
	}

	eventId, err := s.restApiClient.AddTestEvent(restApiEventOk)
	s.Require().NoError(err)
	s.Require().NotEqual("", eventId)
	restApiEventOk.ID = eventId

	dbEvent, err := s.repo.Get(eventId)
	s.Require().NoError(err)
	s.Require().Equal(restApiEventOk, dbEvent)

	err = s.repo.Delete(eventId)
	s.Require().NoError(err)
}
