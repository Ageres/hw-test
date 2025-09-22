//------------go:build integration

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

type CalendarRestApiIntegrationSuite struct {
	suite.Suite
	apiClient c.TestCalendarApiClient
	repo      repo.Repo
}

func (s *CalendarRestApiIntegrationSuite) SetupSuite() {
	s.repo = repo.NewRepo()
	s.apiClient = ch.NewRestAPIClient()
}

func (s *CalendarRestApiIntegrationSuite) TearDownSuite() {
	userIDs := []string{
		"user-id-TestAddEventByRestApi",
		"user-id-TestBusyDateByRestApi",
		"user-id-ok-TestUserConflictErrorByRestApi",
		"user-id-conflict-TestUserConflictErrorByRestApi",
		"user-id-01-TestListDayEventsByRestApi",
		"user-id-02-TestListDayEventsByRestApi",
		"user-id-03-TestListDayEventsByRestApi",
		"user-id-01-TestListWeekEventsByRestApi",
		"user-id-02-TestListWeekEventsByRestApi",
		"user-id-03-TestListWeekEventsByRestApi",
		"user-id-01-TestListMonthEventsByRestApi",
		"user-id-02-TestListMonthEventsByRestApi",
		"user-id-03-TestListMonthEventsByRestApi",
	}
	for _, userID := range userIDs {
		_ = s.repo.DeleteByUserId(userID)
	}
}

func NewTestRestApiSuite() *CalendarRestApiIntegrationSuite {
	return &CalendarRestApiIntegrationSuite{}
}

func TestRestApiSuite(t *testing.T) {
	suite.Run(t, NewTestRestApiSuite())
}

func (s *CalendarRestApiIntegrationSuite) TestAddEventByRestApi() {
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

	eventId, _, err := s.apiClient.AddTestEvent(restApiEvent)
	s.Require().NoError(err)
	s.Require().NotEqual("", eventId)
	restApiEvent.ID = eventId

	dbEvent, err := s.repo.Get(eventId)
	s.Require().NoError(err)
	s.Require().Equal(restApiEvent, dbEvent)

	err = s.repo.DeleteByUserId(restApiEvent.UserID)
	s.Require().NoError(err)
}

func (s *CalendarRestApiIntegrationSuite) TestBusyDateErrorByRestApi() {
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

	eventOkId, _, err := s.apiClient.AddTestEvent(restApiEventOk)
	s.Require().NoError(err)
	s.Require().NotEqual("", eventOkId)
	restApiEventOk.ID = eventOkId

	eventBusyId, bodyBusy, err := s.apiClient.AddTestEvent(restApiEventBusy)
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

func (s *CalendarRestApiIntegrationSuite) TestUserConflictErrorByRestApi() {
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

	eventOkId, _, err := s.apiClient.AddTestEvent(restApiEventOk)
	s.Require().NoError(err)
	s.Require().NotEqual("", eventOkId)
	restApiEventOk.ID = eventOkId

	restApiEventConflict.ID = eventOkId
	bodyConflict, err := s.apiClient.UpdateTestEvent(restApiEventConflict)
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

func (s *CalendarRestApiIntegrationSuite) TestListDayEventsByRestApi() {
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
	eventOneId, _, err := s.apiClient.AddTestEvent(eventOne)
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
	eventTwoId, _, err := s.apiClient.AddTestEvent(eventTwo)
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
	eventThreeId, _, err := s.apiClient.AddTestEvent(eventThree)
	s.Require().NoError(err)
	s.Require().NotEqual("", eventThreeId)
	eventThree.ID = eventThreeId
	dbEventThree, err := s.repo.Get(eventThreeId)
	s.Require().NoError(err)
	s.Require().Equal(eventThree, dbEventThree)

	eventNotInPeriod := &model.TestEvent{
		Title:       "title 03 NotInPeriod TestListDayEventsByRestApi",
		StartTime:   startTime.AddDate(0, 0, 2),
		Duration:    1 * time.Hour,
		Description: "description 03 NotInPeriod TestListDayEventsByRestApi",
		UserID:      "user-id-03-TestListDayEventsByRestApi",
		Reminder:    24 * time.Hour,
	}
	eventNotInPeriodId, _, err := s.apiClient.AddTestEvent(eventNotInPeriod)
	s.Require().NoError(err)
	s.Require().NotEqual("", eventNotInPeriodId)
	eventNotInPeriod.ID = eventNotInPeriodId
	dbEventNotInPeriod, err := s.repo.Get(eventNotInPeriodId)
	s.Require().NoError(err)
	s.Require().Equal(eventNotInPeriod, dbEventNotInPeriod)

	events, _, err := s.apiClient.ListTestEvent(c.DAY, startTime)
	s.Require().NoError(err)
	s.Require().Equal(3, len(events))
	s.Require().True(slices.Contains(events, *eventOne))
	s.Require().True(slices.Contains(events, *eventTwo))
	s.Require().True(slices.Contains(events, *eventThree))
	s.Require().False(slices.Contains(events, *eventNotInPeriod))

	err = s.repo.DeleteByUserId(eventOne.UserID)
	s.Require().NoError(err)
	err = s.repo.DeleteByUserId(eventTwo.UserID)
	s.Require().NoError(err)
	err = s.repo.DeleteByUserId(eventThree.UserID)
	s.Require().NoError(err)
}

func (s *CalendarRestApiIntegrationSuite) TestListWeekEventsByRestApi() {
	timeLocation, err := time.LoadLocation("Local")
	s.Require().NoError(err)
	startTime := time.Date(2030, 12, 31, 10, 0, 0, 0, timeLocation)

	eventOne := &model.TestEvent{
		Title:       "title 01 TestListWeekEventsByRestApi",
		StartTime:   startTime,
		Duration:    1 * time.Hour,
		Description: "description 01 TestListWeekEventsByRestApi",
		UserID:      "user-id-01-TestListWeekEventsByRestApi",
		Reminder:    24 * time.Hour,
	}
	eventOneId, _, err := s.apiClient.AddTestEvent(eventOne)
	s.Require().NoError(err)
	s.Require().NotEqual("", eventOneId)
	eventOne.ID = eventOneId
	dbEventOne, err := s.repo.Get(eventOneId)
	s.Require().NoError(err)
	s.Require().Equal(eventOne, dbEventOne)

	eventTwo := &model.TestEvent{
		Title:       "title 02 TestListWeekEventsByRestApi",
		StartTime:   startTime.AddDate(0, 0, 2),
		Duration:    1 * time.Hour,
		Description: "description 02 TestListWeekEventsByRestApi",
		UserID:      "user-id-02-TestListWeekEventsByRestApi",
		Reminder:    24 * time.Hour,
	}
	eventTwoId, _, err := s.apiClient.AddTestEvent(eventTwo)
	s.Require().NoError(err)
	s.Require().NotEqual("", eventTwoId)
	eventTwo.ID = eventTwoId
	dbEventTwo, err := s.repo.Get(eventTwoId)
	s.Require().NoError(err)
	s.Require().Equal(eventTwo, dbEventTwo)

	eventThree := &model.TestEvent{
		Title:       "title 03 TestListWeekEventsByRestApi",
		StartTime:   startTime.AddDate(0, 0, 4),
		Duration:    1 * time.Hour,
		Description: "description 03 TestListWeekEventsByRestApi",
		UserID:      "user-id-03-TestListWeekEventsByRestApi",
		Reminder:    24 * time.Hour,
	}
	eventThreeId, _, err := s.apiClient.AddTestEvent(eventThree)
	s.Require().NoError(err)
	s.Require().NotEqual("", eventThreeId)
	eventThree.ID = eventThreeId
	dbEventThree, err := s.repo.Get(eventThreeId)
	s.Require().NoError(err)
	s.Require().Equal(eventThree, dbEventThree)

	eventNotInPeriod := &model.TestEvent{
		Title:       "title 03 NotInPeriod TestListWeekEventsByRestApi",
		StartTime:   startTime.AddDate(0, 0, 10),
		Duration:    1 * time.Hour,
		Description: "description 03 NotInPeriod TestListWeekEventsByRestApi",
		UserID:      "user-id-03-TestListWeekEventsByRestApi",
		Reminder:    24 * time.Hour,
	}
	eventNotInPeriodId, _, err := s.apiClient.AddTestEvent(eventNotInPeriod)
	s.Require().NoError(err)
	s.Require().NotEqual("", eventNotInPeriodId)
	eventNotInPeriod.ID = eventNotInPeriodId
	dbEventNotInPeriod, err := s.repo.Get(eventNotInPeriodId)
	s.Require().NoError(err)
	s.Require().Equal(eventNotInPeriod, dbEventNotInPeriod)

	events, _, err := s.apiClient.ListTestEvent(c.WEEK, startTime)
	s.Require().NoError(err)
	s.Require().Equal(3, len(events))
	s.Require().True(slices.Contains(events, *eventOne))
	s.Require().True(slices.Contains(events, *eventTwo))
	s.Require().True(slices.Contains(events, *eventThree))
	s.Require().False(slices.Contains(events, *eventNotInPeriod))

	err = s.repo.DeleteByUserId(eventOne.UserID)
	s.Require().NoError(err)
	err = s.repo.DeleteByUserId(eventTwo.UserID)
	s.Require().NoError(err)
	err = s.repo.DeleteByUserId(eventThree.UserID)
	s.Require().NoError(err)
}

func (s *CalendarRestApiIntegrationSuite) TestListMonthEventsByRestApi() {
	timeLocation, err := time.LoadLocation("Local")
	s.Require().NoError(err)
	startTime := time.Date(2031, 1, 1, 10, 0, 0, 0, timeLocation)

	eventOne := &model.TestEvent{
		Title:       "title 01 TestListMonthEventsByRestApi",
		StartTime:   startTime,
		Duration:    1 * time.Hour,
		Description: "description 01 TestListMonthEventsByRestApi",
		UserID:      "user-id-01-TestListMonthEventsByRestApi",
		Reminder:    24 * time.Hour,
	}
	eventOneId, _, err := s.apiClient.AddTestEvent(eventOne)
	s.Require().NoError(err)
	s.Require().NotEqual("", eventOneId)
	eventOne.ID = eventOneId
	dbEventOne, err := s.repo.Get(eventOneId)
	s.Require().NoError(err)
	s.Require().Equal(eventOne, dbEventOne)

	eventTwo := &model.TestEvent{
		Title:       "title 02 TestListMonthEventsByRestApi",
		StartTime:   startTime.AddDate(0, 0, 8),
		Duration:    1 * time.Hour,
		Description: "description 02 TestListMonthEventsByRestApi",
		UserID:      "user-id-02-TestListMonthEventsByRestApi",
		Reminder:    24 * time.Hour,
	}
	eventTwoId, _, err := s.apiClient.AddTestEvent(eventTwo)
	s.Require().NoError(err)
	s.Require().NotEqual("", eventTwoId)
	eventTwo.ID = eventTwoId
	dbEventTwo, err := s.repo.Get(eventTwoId)
	s.Require().NoError(err)
	s.Require().Equal(eventTwo, dbEventTwo)

	eventThree := &model.TestEvent{
		Title:       "title 03 TestListMonthEventsByRestApi",
		StartTime:   startTime.AddDate(0, 0, 16),
		Duration:    1 * time.Hour,
		Description: "description 03 TestListMonthEventsByRestApi",
		UserID:      "user-id-03-TestListMonthEventsByRestApi",
		Reminder:    24 * time.Hour,
	}
	eventThreeId, _, err := s.apiClient.AddTestEvent(eventThree)
	s.Require().NoError(err)
	s.Require().NotEqual("", eventThreeId)
	eventThree.ID = eventThreeId
	dbEventThree, err := s.repo.Get(eventThreeId)
	s.Require().NoError(err)
	s.Require().Equal(eventThree, dbEventThree)

	eventNotInPeriod := &model.TestEvent{
		Title:       "title 03 NotInPeriod TestListMonthEventsByRestApi",
		StartTime:   startTime.AddDate(0, 0, 32),
		Duration:    1 * time.Hour,
		Description: "description 03 NotInPeriod TestListMonthEventsByRestApi",
		UserID:      "user-id-03-TestListMonthEventsByRestApi",
		Reminder:    24 * time.Hour,
	}
	eventNotInPeriodId, _, err := s.apiClient.AddTestEvent(eventNotInPeriod)
	s.Require().NoError(err)
	s.Require().NotEqual("", eventNotInPeriodId)
	eventNotInPeriod.ID = eventNotInPeriodId
	dbEventNotInPeriod, err := s.repo.Get(eventNotInPeriodId)
	s.Require().NoError(err)
	s.Require().Equal(eventNotInPeriod, dbEventNotInPeriod)

	events, _, err := s.apiClient.ListTestEvent(c.MONTH, startTime)
	s.Require().NoError(err)
	s.Require().Equal(3, len(events))
	s.Require().True(slices.Contains(events, *eventOne))
	s.Require().True(slices.Contains(events, *eventTwo))
	s.Require().True(slices.Contains(events, *eventThree))
	s.Require().False(slices.Contains(events, *eventNotInPeriod))

	err = s.repo.DeleteByUserId(eventOne.UserID)
	s.Require().NoError(err)
	err = s.repo.DeleteByUserId(eventTwo.UserID)
	s.Require().NoError(err)
	err = s.repo.DeleteByUserId(eventThree.UserID)
	s.Require().NoError(err)
}
