//go:build integration

package integration

import (
	"context"
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
	apiClient c.TestCalendarAPIClient
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
		_ = s.repo.DeleteByUserID(userID)
	}
}

func NewTestRestApiSuite() *CalendarRestApiIntegrationSuite {
	return &CalendarRestApiIntegrationSuite{}
}

func TestRestApiSuite(t *testing.T) {
	suite.Run(t, NewTestRestApiSuite())
}

func (s *CalendarRestApiIntegrationSuite) TestAddEventByRestApi() {
	ctx := context.Background()
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

	eventID, _, err := s.apiClient.AddTestEvent(ctx, restApiEvent)
	s.Require().NoError(err)
	s.Require().NotEqual("", eventID)
	restApiEvent.ID = eventID

	dbEvent, err := s.repo.Get(eventID)
	s.Require().NoError(err)
	s.Require().Equal(restApiEvent, dbEvent)

	err = s.repo.DeleteByUserID(restApiEvent.UserID)
	s.Require().NoError(err)
}

func (s *CalendarRestApiIntegrationSuite) TestBusyDateErrorByRestApi() {
	ctx := context.Background()
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

	eventOkID, _, err := s.apiClient.AddTestEvent(ctx, restApiEventOk)
	s.Require().NoError(err)
	s.Require().NotEqual("", eventOkID)
	restApiEventOk.ID = eventOkID

	eventBusyID, bodyBusy, err := s.apiClient.AddTestEvent(ctx, restApiEventBusy)
	s.Require().Equal(err.Error(), "response status '409 Conflict'")
	s.Require().Contains(bodyBusy, fmt.Sprintf("add event: time is already taken by another event: %s", eventOkID))
	s.Require().Equal("", eventBusyID)
	restApiEventOk.ID = eventOkID

	dbEvents, err := s.repo.ListByUserID(userID)
	s.Require().NoError(err)
	s.Require().Equal(1, len(dbEvents))

	err = s.repo.DeleteByUserID(userID)
	s.Require().NoError(err)
}

func (s *CalendarRestApiIntegrationSuite) TestUserConflictErrorByRestApi() {
	ctx := context.Background()
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

	eventOkID, _, err := s.apiClient.AddTestEvent(ctx, restApiEventOk)
	s.Require().NoError(err)
	s.Require().NotEqual("", eventOkID)
	restApiEventOk.ID = eventOkID

	restApiEventConflict.ID = eventOkID
	bodyConflict, err := s.apiClient.UpdateTestEvent(ctx, restApiEventConflict)
	s.Require().Equal(err.Error(), "response status '409 Conflict'")
	s.Require().Contains(bodyConflict, fmt.Sprintf("'%s' user is not the owner of the event, conflict with '%s'", userIDConflict, userIDOk))
	restApiEventOk.ID = eventOkID

	dbEventOks, err := s.repo.ListByUserID(userIDOk)
	s.Require().NoError(err)
	s.Require().Equal(1, len(dbEventOks))

	dbEventConflicts, err := s.repo.ListByUserID(userIDConflict)
	s.Require().Error(err, "not found events for user_id 'user-id-conflict-TestUserConflictErrorByRestApi'")
	s.Require().Equal(0, len(dbEventConflicts))

	err = s.repo.DeleteByUserID(userIDOk)
	s.Require().NoError(err)
}

func (s *CalendarRestApiIntegrationSuite) TestListDayEventsByRestApi() {
	ctx := context.Background()
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
	eventOneID, _, err := s.apiClient.AddTestEvent(ctx, eventOne)
	s.Require().NoError(err)
	s.Require().NotEqual("", eventOneID)
	eventOne.ID = eventOneID
	dbEventOne, err := s.repo.Get(eventOneID)
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
	eventTwoID, _, err := s.apiClient.AddTestEvent(ctx, eventTwo)
	s.Require().NoError(err)
	s.Require().NotEqual("", eventTwoID)
	eventTwo.ID = eventTwoID
	dbEventTwo, err := s.repo.Get(eventTwoID)
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
	eventThreeID, _, err := s.apiClient.AddTestEvent(ctx, eventThree)
	s.Require().NoError(err)
	s.Require().NotEqual("", eventThreeID)
	eventThree.ID = eventThreeID
	dbEventThree, err := s.repo.Get(eventThreeID)
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
	eventNotInPeriodID, _, err := s.apiClient.AddTestEvent(ctx, eventNotInPeriod)
	s.Require().NoError(err)
	s.Require().NotEqual("", eventNotInPeriodID)
	eventNotInPeriod.ID = eventNotInPeriodID
	dbEventNotInPeriod, err := s.repo.Get(eventNotInPeriodID)
	s.Require().NoError(err)
	s.Require().Equal(eventNotInPeriod, dbEventNotInPeriod)

	events, _, err := s.apiClient.ListTestEvent(ctx, c.DAY, startTime)
	s.Require().NoError(err)
	s.Require().Equal(3, len(events))
	s.Require().True(slices.Contains(events, *eventOne))
	s.Require().True(slices.Contains(events, *eventTwo))
	s.Require().True(slices.Contains(events, *eventThree))
	s.Require().False(slices.Contains(events, *eventNotInPeriod))

	err = s.repo.DeleteByUserID(eventOne.UserID)
	s.Require().NoError(err)
	err = s.repo.DeleteByUserID(eventTwo.UserID)
	s.Require().NoError(err)
	err = s.repo.DeleteByUserID(eventThree.UserID)
	s.Require().NoError(err)
}

func (s *CalendarRestApiIntegrationSuite) TestListWeekEventsByRestApi() {
	ctx := context.Background()
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
	eventOneID, _, err := s.apiClient.AddTestEvent(ctx, eventOne)
	s.Require().NoError(err)
	s.Require().NotEqual("", eventOneID)
	eventOne.ID = eventOneID
	dbEventOne, err := s.repo.Get(eventOneID)
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
	eventTwoID, _, err := s.apiClient.AddTestEvent(ctx, eventTwo)
	s.Require().NoError(err)
	s.Require().NotEqual("", eventTwoID)
	eventTwo.ID = eventTwoID
	dbEventTwo, err := s.repo.Get(eventTwoID)
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
	eventThreeID, _, err := s.apiClient.AddTestEvent(ctx, eventThree)
	s.Require().NoError(err)
	s.Require().NotEqual("", eventThreeID)
	eventThree.ID = eventThreeID
	dbEventThree, err := s.repo.Get(eventThreeID)
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
	eventNotInPeriodID, _, err := s.apiClient.AddTestEvent(ctx, eventNotInPeriod)
	s.Require().NoError(err)
	s.Require().NotEqual("", eventNotInPeriodID)
	eventNotInPeriod.ID = eventNotInPeriodID
	dbEventNotInPeriod, err := s.repo.Get(eventNotInPeriodID)
	s.Require().NoError(err)
	s.Require().Equal(eventNotInPeriod, dbEventNotInPeriod)

	events, _, err := s.apiClient.ListTestEvent(ctx, c.WEEK, startTime)
	s.Require().NoError(err)
	s.Require().Equal(3, len(events))
	s.Require().True(slices.Contains(events, *eventOne))
	s.Require().True(slices.Contains(events, *eventTwo))
	s.Require().True(slices.Contains(events, *eventThree))
	s.Require().False(slices.Contains(events, *eventNotInPeriod))

	err = s.repo.DeleteByUserID(eventOne.UserID)
	s.Require().NoError(err)
	err = s.repo.DeleteByUserID(eventTwo.UserID)
	s.Require().NoError(err)
	err = s.repo.DeleteByUserID(eventThree.UserID)
	s.Require().NoError(err)
}

func (s *CalendarRestApiIntegrationSuite) TestListMonthEventsByRestApi() {
	ctx := context.Background()
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
	eventOneID, _, err := s.apiClient.AddTestEvent(ctx, eventOne)
	s.Require().NoError(err)
	s.Require().NotEqual("", eventOneID)
	eventOne.ID = eventOneID
	dbEventOne, err := s.repo.Get(eventOneID)
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
	eventTwoID, _, err := s.apiClient.AddTestEvent(ctx, eventTwo)
	s.Require().NoError(err)
	s.Require().NotEqual("", eventTwoID)
	eventTwo.ID = eventTwoID
	dbEventTwo, err := s.repo.Get(eventTwoID)
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
	eventThreeID, _, err := s.apiClient.AddTestEvent(ctx, eventThree)
	s.Require().NoError(err)
	s.Require().NotEqual("", eventThreeID)
	eventThree.ID = eventThreeID
	dbEventThree, err := s.repo.Get(eventThreeID)
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
	eventNotInPeriodID, _, err := s.apiClient.AddTestEvent(ctx, eventNotInPeriod)
	s.Require().NoError(err)
	s.Require().NotEqual("", eventNotInPeriodID)
	eventNotInPeriod.ID = eventNotInPeriodID
	dbEventNotInPeriod, err := s.repo.Get(eventNotInPeriodID)
	s.Require().NoError(err)
	s.Require().Equal(eventNotInPeriod, dbEventNotInPeriod)

	events, _, err := s.apiClient.ListTestEvent(ctx, c.MONTH, startTime)
	s.Require().NoError(err)
	s.Require().Equal(3, len(events))
	s.Require().True(slices.Contains(events, *eventOne))
	s.Require().True(slices.Contains(events, *eventTwo))
	s.Require().True(slices.Contains(events, *eventThree))
	s.Require().False(slices.Contains(events, *eventNotInPeriod))

	err = s.repo.DeleteByUserID(eventOne.UserID)
	s.Require().NoError(err)
	err = s.repo.DeleteByUserID(eventTwo.UserID)
	s.Require().NoError(err)
	err = s.repo.DeleteByUserID(eventThree.UserID)
	s.Require().NoError(err)
}
