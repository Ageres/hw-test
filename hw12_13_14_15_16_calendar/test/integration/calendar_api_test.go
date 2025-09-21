////go:build integration

package integration

import (
	"fmt"
	"slices"
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

func (s *CalendarIntegrationSuite) TearDownSuite() {
	userIDs := []string{
		"user-id-TestAddEventByRestApi",
		"user-id-TestBusyDateByRestApi",
		"user-id-ok-TestUserConflictErrorByRestApi",
		"user-id-conflict-TestUserConflictErrorByRestApi",
		"user-id-01-TestListDayEventsByRestApi",
		"user-id-02-TestListDayEventsByRestApi",
		"user-id-03-TestListDayEventsByRestApi",
	}
	for _, userID := range userIDs {
		_ = s.repo.DeleteByUserId(userID)
	}
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

	eventId, _, err := s.restApiClient.AddTestEvent(restApiEvent)
	s.Require().NoError(err)
	s.Require().NotEqual("", eventId)
	restApiEvent.ID = eventId

	dbEvent, err := s.repo.Get(eventId)
	s.Require().NoError(err)
	s.Require().Equal(restApiEvent, dbEvent)

	err = s.repo.DeleteByUserId(restApiEvent.UserID)
	s.Require().NoError(err)
}

func (s *CalendarIntegrationSuite) TestBusyDateErrorByRestApi() {
	userID := "user-id-TestBusyDateByRestApi"
	timeLocation, err := time.LoadLocation("Local")
	s.Require().NoError(err)
	startTime := time.Date(2030, 12, 31, 10, 0, 0, 0, timeLocation)
	restApiEventOk := &model.TestEvent{
		Title:       "title ok TestBusyDateByRestApi",
		StartTime:   startTime,
		Duration:    24 * time.Hour,
		Description: "description ok TestBusyDateByRestApi",
		UserID:      userID,
		Reminder:    24 * time.Hour,
	}
	restApiEventBusy := &model.TestEvent{
		Title:       "title busy TestBusyDateByRestApi",
		StartTime:   startTime,
		Duration:    24 * time.Hour,
		Description: "description busy TestBusyDateByRestApi",
		UserID:      userID,
		Reminder:    24 * time.Hour,
	}

	eventOkId, _, err := s.restApiClient.AddTestEvent(restApiEventOk)
	s.Require().NoError(err)
	s.Require().NotEqual("", eventOkId)
	restApiEventOk.ID = eventOkId

	eventBusyId, bodyBusy, err := s.restApiClient.AddTestEvent(restApiEventBusy)
	s.Require().Error(err, "response status '409 Conflict'")
	s.Require().Contains(bodyBusy, fmt.Sprintf("add event: time is already taken by another event: %s", eventOkId))
	s.Require().Equal("", eventBusyId)
	restApiEventOk.ID = eventOkId

	dbEvents, err := s.repo.ListByUserId(userID)
	s.Require().NoError(err)
	s.Require().Equal(1, len(dbEvents))

	err = s.repo.DeleteByUserId(userID)
	s.Require().NoError(err)
}

func (s *CalendarIntegrationSuite) TestUserConflictErrorByRestApi() {
	userIDOk := "user-id-ok-TestUserConflictErrorByRestApi"
	userIDConflict := "user-id-conflict-TestUserConflictErrorByRestApi"
	timeLocation, err := time.LoadLocation("Local")
	s.Require().NoError(err)
	startTime := time.Date(2030, 12, 31, 10, 0, 0, 0, timeLocation)
	restApiEventOk := &model.TestEvent{
		Title:       "title ok TestUserConflictErrorByRest",
		StartTime:   startTime,
		Duration:    24 * time.Hour,
		Description: "description ok TestUserConflictErrorByRest",
		UserID:      userIDOk,
		Reminder:    24 * time.Hour,
	}
	restApiEventConflict := &model.TestEvent{
		Title:       "title conflict TestUserConflictErrorByRest",
		StartTime:   startTime,
		Duration:    24 * time.Hour,
		Description: "description conflict TestUserConflictErrorByRest",
		UserID:      userIDConflict,
		Reminder:    24 * time.Hour,
	}

	eventOkId, _, err := s.restApiClient.AddTestEvent(restApiEventOk)
	s.Require().NoError(err)
	s.Require().NotEqual("", eventOkId)
	restApiEventOk.ID = eventOkId

	restApiEventConflict.ID = eventOkId
	bodyConflict, err := s.restApiClient.UpdateTestEvent(restApiEventConflict)
	s.Require().Error(err, "response status '409 Conflict'")
	s.Require().Contains(bodyConflict, fmt.Sprintf("'%s' user is not the owner of the event, conflict with '%s'", userIDConflict, userIDOk))
	restApiEventOk.ID = eventOkId

	dbEventOks, err := s.repo.ListByUserId(userIDOk)
	s.Require().NoError(err)
	s.Require().Equal(1, len(dbEventOks))

	dbEventConflicts, err := s.repo.ListByUserId(userIDConflict)
	s.Require().Error(err, "not found events for user_id 'user-id-conflict-TestUserConflictErrorByRestApi'")
	s.Require().Equal(0, len(dbEventConflicts))

	err = s.repo.DeleteByUserId(userIDOk)
	s.Require().NoError(err)
}

func (s *CalendarIntegrationSuite) TestListDayEventsByRestApi() {
	timeLocation, err := time.LoadLocation("Local")
	s.Require().NoError(err)
	startTime := time.Date(2030, 12, 31, 10, 0, 0, 0, timeLocation)

	eventOne := &model.TestEvent{
		Title:       "title 01 TestListDayEventsByRestApi",
		StartTime:   startTime,
		Duration:    1 * time.Hour,
		Description: "description 01 TestListDayEventsByRestApi",
		UserID:      "user-id-01-TestListDayEventsByRestApi",
		Reminder:    24 * time.Hour,
	}
	eventOneId, _, err := s.restApiClient.AddTestEvent(eventOne)
	s.Require().NoError(err)
	s.Require().NotEqual("", eventOneId)
	eventOne.ID = eventOneId
	dbEventOne, err := s.repo.Get(eventOneId)
	s.Require().NoError(err)
	s.Require().Equal(eventOne, dbEventOne)

	eventTwo := &model.TestEvent{
		Title:       "title 02 TestListDayEventsByRestApi",
		StartTime:   startTime.Add(2 * time.Hour),
		Duration:    1 * time.Hour,
		Description: "description 02 TestListDayEventsByRestApi",
		UserID:      "user-id-02-TestListDayEventsByRestApi",
		Reminder:    24 * time.Hour,
	}
	eventTwoId, _, err := s.restApiClient.AddTestEvent(eventTwo)
	s.Require().NoError(err)
	s.Require().NotEqual("", eventTwoId)
	eventTwo.ID = eventTwoId
	dbEventTwo, err := s.repo.Get(eventTwoId)
	s.Require().NoError(err)
	s.Require().Equal(eventTwo, dbEventTwo)

	eventThree := &model.TestEvent{
		Title:       "title 03 TestListDayEventsByRestApi",
		StartTime:   startTime.Add(4 * time.Hour),
		Duration:    1 * time.Hour,
		Description: "description 03 TestListDayEventsByRestApi",
		UserID:      "user-id-03-TestListDayEventsByRestApi",
		Reminder:    24 * time.Hour,
	}
	eventThreeId, _, err := s.restApiClient.AddTestEvent(eventThree)
	s.Require().NoError(err)
	s.Require().NotEqual("", eventThreeId)
	eventThree.ID = eventThreeId
	dbEventThree, err := s.repo.Get(eventThreeId)
	s.Require().NoError(err)
	s.Require().Equal(eventThree, dbEventThree)

	events, _, err := s.restApiClient.ListTestEvent(c.DAY, startTime)
	s.Require().NoError(err)
	s.Require().Equal(3, len(events))
	for _, e := range events {
		s.Require().True(slices.Contains(events, e))
	}

	err = s.repo.DeleteByUserId(eventOne.UserID)
	s.Require().NoError(err)
	err = s.repo.DeleteByUserId(eventTwo.UserID)
	s.Require().NoError(err)
	err = s.repo.DeleteByUserId(eventThree.UserID)
	s.Require().NoError(err)
}
