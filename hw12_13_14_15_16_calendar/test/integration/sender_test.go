//go:build integration

package integration

import (
	"context"
	"testing"
	"time"

	c "github.com/Ageres/hw-test/hw12_13_14_15_calendar/test/integration/client"
	ch "github.com/Ageres/hw-test/hw12_13_14_15_calendar/test/integration/client/http"
	"github.com/Ageres/hw-test/hw12_13_14_15_calendar/test/integration/model"
	"github.com/Ageres/hw-test/hw12_13_14_15_calendar/test/integration/repo"
	"github.com/stretchr/testify/suite"
)

const testTimeOut = 15

type SenderIntegrationSuite struct {
	suite.Suite
	restApiClient c.TestCalendarAPIClient
	repo          repo.Repo
}

func (s *SenderIntegrationSuite) SetupSuite() {
	s.repo = repo.NewRepo()
	s.restApiClient = ch.NewRestAPIClient()
}

func (s *SenderIntegrationSuite) TearDownSuite() {
	userIDs := []string{
		"user-id-01-TestSender",
		"user-id-02-TestSender",
		"user-id-03-TestSender",
	}
	for _, userID := range userIDs {
		_ = s.repo.DeleteByUserID(userID)
		_ = s.repo.DeleteProcEventByUserID(userID)
	}
}

func NewTestSenderSuite() *SenderIntegrationSuite {
	return &SenderIntegrationSuite{}
}

func TestSenderSuite(t *testing.T) {
	suite.Run(t, NewTestSenderSuite())
}

func (s *SenderIntegrationSuite) TestSender() {
	ctx := context.Background()
	now := time.Now().Local()

	eventOne := &model.TestEvent{
		Title:       "title 01 TestSender",
		StartTime:   now.Add(24 * time.Hour).Add(11 * time.Second),
		Duration:    1 * time.Hour,
		Description: "description 01 TestSender",
		UserID:      "user-id-01-TestSender",
		Reminder:    24 * time.Hour,
	}
	eventOneID, _, err := s.restApiClient.AddTestEvent(ctx, eventOne)
	s.Require().NoError(err)
	s.Require().NotEqual("", eventOne)

	eventTwo := &model.TestEvent{
		Title:       "title 02 TestSender",
		StartTime:   now.Add(48 * time.Hour).Add(11 * time.Second),
		Duration:    1 * time.Hour,
		Description: "description 02 TestSender",
		UserID:      "user-id-02-TestSender",
		Reminder:    48 * time.Hour,
	}
	eventTwoID, _, err := s.restApiClient.AddTestEvent(ctx, eventTwo)
	s.Require().NoError(err)
	s.Require().NotEqual("", eventTwoID)

	eventThree := &model.TestEvent{
		Title:       "title 03 TestSender",
		StartTime:   now.Add(96 * time.Hour).Add(11 * time.Second),
		Duration:    1 * time.Hour,
		Description: "description 03 TestSender",
		UserID:      "user-id-03-TestSender",
		Reminder:    96 * time.Hour,
	}
	eventThreeID, _, err := s.restApiClient.AddTestEvent(ctx, eventThree)
	s.Require().NoError(err)
	s.Require().NotEqual("", eventThreeID)

	eventNotInPeriod := &model.TestEvent{
		Title:       "title 03 NotInPeriod TestSender",
		StartTime:   now.AddDate(0, 0, 32),
		Duration:    1 * time.Hour,
		Description: "description 03 NotInPeriod TestSender",
		UserID:      "user-id-03-TestSender",
		Reminder:    24 * time.Hour,
	}
	eventNotInPeriodID, _, err := s.restApiClient.AddTestEvent(ctx, eventNotInPeriod)
	s.Require().NoError(err)
	s.Require().NotEqual("", eventNotInPeriodID)

	time.Sleep(time.Duration(testTimeOut * time.Second))

	isExistProcEventOne, err := s.repo.CheckProcEvent(eventOneID)
	s.Require().NoError(err)
	s.Require().True(isExistProcEventOne)

	isExistProcEventTwo, err := s.repo.CheckProcEvent(eventTwoID)
	s.Require().NoError(err)
	s.Require().True(isExistProcEventTwo)

	isExistProcEventThree, err := s.repo.CheckProcEvent(eventThreeID)
	s.Require().NoError(err)
	s.Require().True(isExistProcEventThree)

	isExistProcEventNotInPeriodID, err := s.repo.CheckProcEvent(eventNotInPeriodID)
	s.Require().NoError(err)
	s.Require().False(isExistProcEventNotInPeriodID)

	err = s.repo.DeleteByUserID(eventOne.UserID)
	s.Require().NoError(err)
	err = s.repo.DeleteByUserID(eventTwo.UserID)
	s.Require().NoError(err)
	err = s.repo.DeleteByUserID(eventThree.UserID)
	s.Require().NoError(err)

	err = s.repo.DeleteProcEventByUserID(eventOne.UserID)
	s.Require().NoError(err)
	err = s.repo.DeleteProcEventByUserID(eventTwo.UserID)
	s.Require().NoError(err)
	err = s.repo.DeleteProcEventByUserID(eventThree.UserID)
	s.Require().NoError(err)
}
