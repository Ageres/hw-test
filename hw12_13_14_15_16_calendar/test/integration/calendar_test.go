package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/stretchr/testify/suite"
)

type TestCalendarApiClient interface {
	AddTestEvent(eventRef *TestEvent) (string, error)
}

type restApiClient struct {
	url        string
	httpClient *http.Client
}

func newRestapiClient() TestCalendarApiClient {
	restApiHost, isSet := os.LookupEnv("CALENDAR_REST_API_HOST")
	if !isSet {
		restApiHost = "localhost"
		log.Println("not found calendar rest api host, set default 'localhost'")
	}
	restApiPort, isSet := os.LookupEnv("CALENDAR_REST_API_PORT")
	if !isSet {
		restApiPort = "8888"
		log.Println("not found calendar rest api port, set default '8888'")
	}
	return &restApiClient{
		url:        fmt.Sprintf("http://%s:%s/v1/event", restApiHost, restApiPort),
		httpClient: http.DefaultClient,
	}
}

func (c *restApiClient) AddTestEvent(eventRef *TestEvent) (string, error) {
	jsonBody, err := json.Marshal(eventRef)
	if err != nil {
		return "", err
	}
	req, err := http.NewRequest(http.MethodPost, c.url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return "", err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("response status '%s'", resp.Status)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	log.Printf("----------------- body: '%s'", string(body))
	respEventRef := new(TestEvent)
	err = json.Unmarshal(body, respEventRef)
	if err != nil {
		return "", err
	}
	return respEventRef.ID, nil
}

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

	dbEvent, err := s.repo.Get(context.Background(), eventId)
	s.Require().NoError(err)
	s.Require().Equal(restApiEvent, dbEvent)
}
