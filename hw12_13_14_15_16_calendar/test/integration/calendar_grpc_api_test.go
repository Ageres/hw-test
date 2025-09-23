//------------------go:build integration

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
	apiClient c.TestCalendarAPIClient
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
		_ = s.repo.DeleteByUserID(userID)
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

	eventID, _, err := s.apiClient.AddTestEvent(grpcApiEvent)
	s.Require().NoError(err)
	s.Require().NotEqual("", eventID)
	grpcApiEvent.ID = eventID

	dbEvent, err := s.repo.Get(eventID)
	s.Require().NoError(err)
	s.Require().Equal(grpcApiEvent, dbEvent)

	err = s.repo.DeleteByUserID(grpcApiEvent.UserID)
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

	eventOkID, _, err := s.apiClient.AddTestEvent(grpcApiEventOk)
	s.Require().NoError(err)
	s.Require().NotEqual("", eventOkID)
	grpcApiEventOk.ID = eventOkID

	eventBusyID, _, err := s.apiClient.AddTestEvent(grpcApiEventBusy)
	s.Require().Contains(err.Error(), fmt.Sprintf("time is already taken by another event: %s", eventOkID))
	s.Require().Equal("", eventBusyID)
	grpcApiEventOk.ID = eventOkID

	dbEvents, err := s.repo.ListByUserID(userID)
	s.Require().NoError(err)
	s.Require().Equal(1, len(dbEvents))

	err = s.repo.DeleteByUserID(userID)
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

	eventOkID, _, err := s.apiClient.AddTestEvent(grpcApiEventOk)
	s.Require().NoError(err)
	s.Require().NotEqual("", eventOkID)
	grpcApiEventOk.ID = eventOkID

	grpcApiEventConflict.ID = eventOkID
	_, err = s.apiClient.UpdateTestEvent(grpcApiEventConflict)
	s.Require().Contains(err.Error(), fmt.Sprintf("'%s' user is not the owner of the event, conflict with '%s'", userIDConflict, userIDOk))
	grpcApiEventOk.ID = eventOkID

	dbEventOks, err := s.repo.ListByUserID(userIDOk)
	s.Require().NoError(err)
	s.Require().Equal(1, len(dbEventOks))

	dbEventConflicts, err := s.repo.ListByUserID(userIDConflict)
	s.Require().Equal(err.Error(), "not found events for user_id 'user-id-conflict-TestUserConflictErrorByGrpcApi'")
	s.Require().Equal(0, len(dbEventConflicts))

	err = s.repo.DeleteByUserID(userIDOk)
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
	eventOneID, _, err := s.apiClient.AddTestEvent(eventOne)
	s.Require().NoError(err)
	s.Require().NotEqual("", eventOneID)
	eventOne.ID = eventOneID
	dbEventOne, err := s.repo.Get(eventOneID)
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
	eventTwoID, _, err := s.apiClient.AddTestEvent(eventTwo)
	s.Require().NoError(err)
	s.Require().NotEqual("", eventTwoID)
	eventTwo.ID = eventTwoID
	dbEventTwo, err := s.repo.Get(eventTwoID)
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
	eventThreeID, _, err := s.apiClient.AddTestEvent(eventThree)
	s.Require().NoError(err)
	s.Require().NotEqual("", eventThreeID)
	eventThree.ID = eventThreeID
	dbEventThree, err := s.repo.Get(eventThreeID)
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
	eventNotInPeriodID, _, err := s.apiClient.AddTestEvent(eventNotInPeriod)
	s.Require().NoError(err)
	s.Require().NotEqual("", eventNotInPeriodID)
	eventNotInPeriod.ID = eventNotInPeriodID
	dbEventNotInPeriod, err := s.repo.Get(eventNotInPeriodID)
	s.Require().NoError(err)
	s.Require().Equal(eventNotInPeriod, dbEventNotInPeriod)

	events, _, err := s.apiClient.ListTestEvent(c.DAY, startTime)
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
	eventOneID, _, err := s.apiClient.AddTestEvent(eventOne)
	s.Require().NoError(err)
	s.Require().NotEqual("", eventOneID)
	eventOne.ID = eventOneID
	dbEventOne, err := s.repo.Get(eventOneID)
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
	eventTwoID, _, err := s.apiClient.AddTestEvent(eventTwo)
	s.Require().NoError(err)
	s.Require().NotEqual("", eventTwoID)
	eventTwo.ID = eventTwoID
	dbEventTwo, err := s.repo.Get(eventTwoID)
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
	eventThreeID, _, err := s.apiClient.AddTestEvent(eventThree)
	s.Require().NoError(err)
	s.Require().NotEqual("", eventThreeID)
	eventThree.ID = eventThreeID
	dbEventThree, err := s.repo.Get(eventThreeID)
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
	eventNotInPeriodID, _, err := s.apiClient.AddTestEvent(eventNotInPeriod)
	s.Require().NoError(err)
	s.Require().NotEqual("", eventNotInPeriodID)
	eventNotInPeriod.ID = eventNotInPeriodID
	dbEventNotInPeriod, err := s.repo.Get(eventNotInPeriodID)
	s.Require().NoError(err)
	s.Require().Equal(eventNotInPeriod, dbEventNotInPeriod)

	events, _, err := s.apiClient.ListTestEvent(c.WEEK, startTime)
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
	eventOneID, _, err := s.apiClient.AddTestEvent(eventOne)
	s.Require().NoError(err)
	s.Require().NotEqual("", eventOneID)
	eventOne.ID = eventOneID
	dbEventOne, err := s.repo.Get(eventOneID)
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
	eventTwoID, _, err := s.apiClient.AddTestEvent(eventTwo)
	s.Require().NoError(err)
	s.Require().NotEqual("", eventTwoID)
	eventTwo.ID = eventTwoID
	dbEventTwo, err := s.repo.Get(eventTwoID)
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
	eventThreeID, _, err := s.apiClient.AddTestEvent(eventThree)
	s.Require().NoError(err)
	s.Require().NotEqual("", eventThreeID)
	eventThree.ID = eventThreeID
	dbEventThree, err := s.repo.Get(eventThreeID)
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
	eventNotInPeriodID, _, err := s.apiClient.AddTestEvent(eventNotInPeriod)
	s.Require().NoError(err)
	s.Require().NotEqual("", eventNotInPeriodID)
	eventNotInPeriod.ID = eventNotInPeriodID
	dbEventNotInPeriod, err := s.repo.Get(eventNotInPeriodID)
	s.Require().NoError(err)
	s.Require().Equal(eventNotInPeriod, dbEventNotInPeriod)

	events, _, err := s.apiClient.ListTestEvent(c.MONTH, startTime)
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
