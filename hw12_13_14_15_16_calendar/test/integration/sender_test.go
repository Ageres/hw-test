//go:build integration

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

const testTimeOut = 15

type SenderIntegrationSuite struct {
	suite.Suite
	restApiClient c.TestCalendarApiClient
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
		_ = s.repo.DeleteByUserId(userID)
		_ = s.repo.DeleteProcEventByUserId(userID)
	}
}

func NewTestSenderSuite() *SenderIntegrationSuite {
	return &SenderIntegrationSuite{}
}

func TestSenderSuite(t *testing.T) {
	suite.Run(t, NewTestSenderSuite())
}

func (s *SenderIntegrationSuite) TestSender() {
	now := time.Now().Local()

	eventOne := &model.TestEvent{
		Title:       "title 01 TestSender",
		StartTime:   now.Add(24 * time.Hour).Add(11 * time.Second),
		Duration:    1 * time.Hour,
		Description: "description 01 TestSender",
		UserID:      "user-id-01-TestSender",
		Reminder:    24 * time.Hour,
	}
	eventOneId, _, err := s.restApiClient.AddTestEvent(eventOne)
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
	eventTwoId, _, err := s.restApiClient.AddTestEvent(eventTwo)
	s.Require().NoError(err)
	s.Require().NotEqual("", eventTwoId)

	eventThree := &model.TestEvent{
		Title:       "title 03 TestSender",
		StartTime:   now.Add(96 * time.Hour).Add(11 * time.Second),
		Duration:    1 * time.Hour,
		Description: "description 03 TestSender",
		UserID:      "user-id-03-TestSender",
		Reminder:    96 * time.Hour,
	}
	eventThreeId, _, err := s.restApiClient.AddTestEvent(eventThree)
	s.Require().NoError(err)
	s.Require().NotEqual("", eventThreeId)

	eventNotInPeriod := &model.TestEvent{
		Title:       "title 03 NotInPeriod TestSender",
		StartTime:   now.AddDate(0, 0, 32),
		Duration:    1 * time.Hour,
		Description: "description 03 NotInPeriod TestSender",
		UserID:      "user-id-03-TestSender",
		Reminder:    24 * time.Hour,
	}
	eventNotInPeriodId, _, err := s.restApiClient.AddTestEvent(eventNotInPeriod)
	s.Require().NoError(err)
	s.Require().NotEqual("", eventNotInPeriodId)

	time.Sleep(time.Duration(testTimeOut * time.Second))

	isExistProcEventOne, err := s.repo.CheckProcEvent(eventOneId)
	s.Require().NoError(err)
	s.Require().True(isExistProcEventOne)

	isExistProcEventTwo, err := s.repo.CheckProcEvent(eventTwoId)
	s.Require().NoError(err)
	s.Require().True(isExistProcEventTwo)

	isExistProcEventThree, err := s.repo.CheckProcEvent(eventThreeId)
	s.Require().NoError(err)
	s.Require().True(isExistProcEventThree)

	isExistProcEventNotInPeriodId, err := s.repo.CheckProcEvent(eventNotInPeriodId)
	s.Require().NoError(err)
	s.Require().False(isExistProcEventNotInPeriodId)

	err = s.repo.DeleteByUserId(eventOne.UserID)
	s.Require().NoError(err)
	err = s.repo.DeleteByUserId(eventTwo.UserID)
	s.Require().NoError(err)
	err = s.repo.DeleteByUserId(eventThree.UserID)
	s.Require().NoError(err)

	err = s.repo.DeleteProcEventByUserId(eventOne.UserID)
	s.Require().NoError(err)
	err = s.repo.DeleteProcEventByUserId(eventTwo.UserID)
	s.Require().NoError(err)
	err = s.repo.DeleteProcEventByUserId(eventThree.UserID)
	s.Require().NoError(err)
}
