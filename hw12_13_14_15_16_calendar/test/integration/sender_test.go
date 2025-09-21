package integration

import (
	c "github.com/Ageres/hw-test/hw12_13_14_15_calendar/test/integration/client"
	"github.com/Ageres/hw-test/hw12_13_14_15_calendar/test/integration/repo"
	"github.com/stretchr/testify/suite"
)

type SenderIntegrationSuite struct {
	suite.Suite
	restApiClient c.TestCalendarApiClient
	repo          repo.Repo
}
