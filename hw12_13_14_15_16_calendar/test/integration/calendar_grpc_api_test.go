// -----------go:build integration

package integration

import (
	"fmt"
	"slices"
	"testing"
	"time"

	c "github.com/Ageres/hw-test/hw12_13_14_15_calendar/test/integration/client"
	ch "github.com/Ageres/hw-test/hw12_13_14_15_calendar/test/integration/client/grpc"
	"github.com/Ageres/hw-test/hw12_13_14_15_calendar/test/integration/model"
	"github.com/Ageres/hw-test/hw12_13_14_15_calendar/test/integration/repo"
	"github.com/stretchr/testify/suite"
)

type CalendarGrpcAPIIntegrationSuite struct {
	suite.Suite
	apiClient c.TestCalendarApiClient
	repo      repo.Repo
}

func (s *CalendarGrpcAPIIntegrationSuite) SetupSuite() {
	s.repo = repo.NewRepo()
	s.apiClient = ch.NewGrpcAPIClient()
}

func (s *CalendarGrpcAPIIntegrationSuite) TearDownSuite() {
	userIDs := []string{
		"user-id-TestAddEventByGrpcApi",
		"user-id-TestBusyDateByGrpcApi",
		"user-id-ok-TestUserConflictErrorByGrpcApi",
		"user-id-conflict-TestUserConflictErrorByGrpcApi",
		"user-id-01-TestListDayEventsByGrpcApi",
		"user-id-02-TestListDayEventsByGrpcApi",
		"user-id-03-TestListDayEventsByGrpcApi",
		"user-id-01-TestListWeekEventsByGrpcApi",
		"user-id-02-TestListWeekEventsByGrpcApi",
		"user-id-03-TestListWeekEventsByGrpcApi",
		"user-id-01-TestListMonthEventsByGrpcApi",
		"user-id-02-TestListMonthEventsByGrpcApi",
		"user-id-03-TestListMonthEventsByGrpcApi",
	}
	for _, userID := range userIDs {
		_ = s.repo.DeleteByUserId(userID)
	}
	s.apiClient.Stop()
}

func NewTestGrpcAPISuite() *CalendarGrpcAPIIntegrationSuite {
	return &CalendarGrpcAPIIntegrationSuite{}
}

func TestGrpcAPISuite(t *testing.T) {
	suite.Run(t, NewTestGrpcAPISuite())
}

func (s *CalendarGrpcAPIIntegrationSuite) TestAddEventByGrpcAPI() {
	timeLocation, err := time.LoadLocation("Local")
	s.Require().NoError(err)
	startTime := time.Date(2030, 12, 31, 10, 0, 0, 0, timeLocation)
	grpcApiEvent := &model.TestEvent{
		Title:       "title TestAddEventByGrpcApi",
		StartTime:   startTime,
		Duration:    24 * time.Hour,
		Description: "description TestAddEventByGrpcApi",
		UserID:      "user-id-TestAddEventByGrpcApi",
		Reminder:    24 * time.Hour,
	}

	eventId, _, err := s.apiClient.AddTestEvent(grpcApiEvent)
	s.Require().NoError(err)
	s.Require().NotEqual("", eventId)
	grpcApiEvent.ID = eventId

	dbEvent, err := s.repo.Get(eventId)
	s.Require().NoError(err)
	s.Require().Equal(grpcApiEvent, dbEvent)

	err = s.repo.DeleteByUserId(grpcApiEvent.UserID)
	s.Require().NoError(err)
}

func (s *CalendarGrpcAPIIntegrationSuite) TestBusyDateErrorByGrpcAPI() {
	userID := "user-id-TestBusyDateByGrpcApi"
	timeLocation, err := time.LoadLocation("Local")
	s.Require().NoError(err)
	startTime := time.Date(2030, 12, 31, 10, 0, 0, 0, timeLocation)
	grpcApiEventOk := &model.TestEvent{
		Title:       "title ok TestBusyDateByGrpcApi",
		StartTime:   startTime,
		Duration:    24 * time.Hour,
		Description: "description ok TestBusyDateByGrpcApi",
		UserID:      userID,
		Reminder:    24 * time.Hour,
	}
	grpcApiEventBusy := &model.TestEvent{
		Title:       "title busy TestBusyDateByGrpcApi",
		StartTime:   startTime,
		Duration:    24 * time.Hour,
		Description: "description busy TestBusyDateByGrpcApi",
		UserID:      userID,
		Reminder:    24 * time.Hour,
	}

	eventOkId, _, err := s.apiClient.AddTestEvent(grpcApiEventOk)
	s.Require().NoError(err)
	s.Require().NotEqual("", eventOkId)
	grpcApiEventOk.ID = eventOkId

	eventBusyId, _, err := s.apiClient.AddTestEvent(grpcApiEventBusy)
	s.Require().Contains(err.Error(), fmt.Sprintf("time is already taken by another event: %s", eventOkId))
	s.Require().Equal("", eventBusyId)
	grpcApiEventOk.ID = eventOkId

	dbEvents, err := s.repo.ListByUserId(userID)
	s.Require().NoError(err)
	s.Require().Equal(1, len(dbEvents))

	err = s.repo.DeleteByUserId(userID)
	s.Require().NoError(err)
}

func (s *CalendarGrpcAPIIntegrationSuite) TestUserConflictErrorByGrpcAPI() {
	userIDOk := "user-id-ok-TestUserConflictErrorByGrpcApi"
	userIDConflict := "user-id-conflict-TestUserConflictErrorByGrpcApi"
	timeLocation, err := time.LoadLocation("Local")
	s.Require().NoError(err)
	startTime := time.Date(2030, 12, 31, 10, 0, 0, 0, timeLocation)
	grpcApiEventOk := &model.TestEvent{
		Title:       "title ok TestUserConflictErrorByGrpc",
		StartTime:   startTime,
		Duration:    24 * time.Hour,
		Description: "description ok TestUserConflictErrorByGrpc",
		UserID:      userIDOk,
		Reminder:    24 * time.Hour,
	}
	grpcApiEventConflict := &model.TestEvent{
		Title:       "title conflict TestUserConflictErrorByGrpc",
		StartTime:   startTime,
		Duration:    24 * time.Hour,
		Description: "description conflict TestUserConflictErrorByGrpc",
		UserID:      userIDConflict,
		Reminder:    24 * time.Hour,
	}

	eventOkId, _, err := s.apiClient.AddTestEvent(grpcApiEventOk)
	s.Require().NoError(err)
	s.Require().NotEqual("", eventOkId)
	grpcApiEventOk.ID = eventOkId

	grpcApiEventConflict.ID = eventOkId
	_, err = s.apiClient.UpdateTestEvent(grpcApiEventConflict)
	s.Require().Contains(err.Error(), fmt.Sprintf("'%s' user is not the owner of the event, conflict with '%s'", userIDConflict, userIDOk))
	grpcApiEventOk.ID = eventOkId

	dbEventOks, err := s.repo.ListByUserId(userIDOk)
	s.Require().NoError(err)
	s.Require().Equal(1, len(dbEventOks))

	dbEventConflicts, err := s.repo.ListByUserId(userIDConflict)
	s.Require().Equal(err.Error(), "not found events for user_id 'user-id-conflict-TestUserConflictErrorByGrpcApi'")
	s.Require().Equal(0, len(dbEventConflicts))

	err = s.repo.DeleteByUserId(userIDOk)
	s.Require().NoError(err)
}

func (s *CalendarGrpcAPIIntegrationSuite) TestListDayEventsByGrpcAPI() {
	timeLocation, err := time.LoadLocation("Local")
	s.Require().NoError(err)
	startTime := time.Date(2030, 12, 31, 10, 0, 0, 0, timeLocation)

	eventOne := &model.TestEvent{
		Title:       "title 01 TestListDayEventsByGrpcApi",
		StartTime:   startTime,
		Duration:    1 * time.Hour,
		Description: "description 01 TestListDayEventsByGrpcApi",
		UserID:      "user-id-01-TestListDayEventsByGrpcApi",
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
		Title:       "title 02 TestListDayEventsByGrpcApi",
		StartTime:   startTime.Add(2 * time.Hour),
		Duration:    1 * time.Hour,
		Description: "description 02 TestListDayEventsByGrpcApi",
		UserID:      "user-id-02-TestListDayEventsByGrpcApi",
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
		Title:       "title 03 TestListDayEventsByGrpcApi",
		StartTime:   startTime.Add(4 * time.Hour),
		Duration:    1 * time.Hour,
		Description: "description 03 TestListDayEventsByGrpcApi",
		UserID:      "user-id-03-TestListDayEventsByGrpcApi",
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
		Title:       "title 03 NotInPeriod TestListDayEventsByGrpcApi",
		StartTime:   startTime.AddDate(0, 0, 2),
		Duration:    1 * time.Hour,
		Description: "description 03 NotInPeriod TestListDayEventsByGrpcApi",
		UserID:      "user-id-03-TestListDayEventsByGrpcApi",
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

func (s *CalendarGrpcAPIIntegrationSuite) TestListWeekEventsByGrpcAPI() {
	timeLocation, err := time.LoadLocation("Local")
	s.Require().NoError(err)
	startTime := time.Date(2030, 12, 31, 10, 0, 0, 0, timeLocation)

	eventOne := &model.TestEvent{
		Title:       "title 01 TestListWeekEventsByGrpcApi",
		StartTime:   startTime,
		Duration:    1 * time.Hour,
		Description: "description 01 TestListWeekEventsByGrpcApi",
		UserID:      "user-id-01-TestListWeekEventsByGrpcApi",
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
		Title:       "title 02 TestListWeekEventsByGrpcApi",
		StartTime:   startTime.AddDate(0, 0, 2),
		Duration:    1 * time.Hour,
		Description: "description 02 TestListWeekEventsByGrpcApi",
		UserID:      "user-id-02-TestListWeekEventsByGrpcApi",
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
		Title:       "title 03 TestListWeekEventsByGrpcApi",
		StartTime:   startTime.AddDate(0, 0, 4),
		Duration:    1 * time.Hour,
		Description: "description 03 TestListWeekEventsByGrpcApi",
		UserID:      "user-id-03-TestListWeekEventsByGrpcApi",
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
		Title:       "title 03 NotInPeriod TestListWeekEventsByGrpcApi",
		StartTime:   startTime.AddDate(0, 0, 10),
		Duration:    1 * time.Hour,
		Description: "description 03 NotInPeriod TestListWeekEventsByGrpcApi",
		UserID:      "user-id-03-TestListWeekEventsByGrpcApi",
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

func (s *CalendarGrpcAPIIntegrationSuite) TestListMonthEventsByGrpcAPI() {
	timeLocation, err := time.LoadLocation("Local")
	s.Require().NoError(err)
	startTime := time.Date(2031, 1, 1, 10, 0, 0, 0, timeLocation)

	eventOne := &model.TestEvent{
		Title:       "title 01 TestListMonthEventsByGrpcApi",
		StartTime:   startTime,
		Duration:    1 * time.Hour,
		Description: "description 01 TestListMonthEventsByGrpcApi",
		UserID:      "user-id-01-TestListMonthEventsByGrpcApi",
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
		Title:       "title 02 TestListMonthEventsByGrpcApi",
		StartTime:   startTime.AddDate(0, 0, 8),
		Duration:    1 * time.Hour,
		Description: "description 02 TestListMonthEventsByGrpcApi",
		UserID:      "user-id-02-TestListMonthEventsByGrpcApi",
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
		Title:       "title 03 TestListMonthEventsByGrpcApi",
		StartTime:   startTime.AddDate(0, 0, 16),
		Duration:    1 * time.Hour,
		Description: "description 03 TestListMonthEventsByGrpcApi",
		UserID:      "user-id-03-TestListMonthEventsByGrpcApi",
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
		Title:       "title 03 NotInPeriod TestListMonthEventsByGrpcApi",
		StartTime:   startTime.AddDate(0, 0, 32),
		Duration:    1 * time.Hour,
		Description: "description 03 NotInPeriod TestListMonthEventsByGrpcApi",
		UserID:      "user-id-03-TestListMonthEventsByGrpcApi",
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
