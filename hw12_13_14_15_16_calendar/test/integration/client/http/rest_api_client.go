package integration

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	c "github.com/Ageres/hw-test/hw12_13_14_15_calendar/test/integration/client"
	"github.com/Ageres/hw-test/hw12_13_14_15_calendar/test/integration/config"
	"github.com/Ageres/hw-test/hw12_13_14_15_calendar/test/integration/model"
	"github.com/Ageres/hw-test/hw12_13_14_15_calendar/test/integration/utils"
)

type restApiClient struct {
	url        string
	httpClient *http.Client
}

func NewRestapiClient() c.TestCalendarApiClient {
	restApiHost := utils.GetEnvOrDefault(config.CALENDAR_REST_API_HOST_ENV, config.CALENDAR_REST_API_HOST_DEFAULT)
	restApiPort := utils.GetEnvOrDefault(config.CALENDAR_REST_API_PORT_ENV, config.CALENDAR_REST_API_PORT_DEFAULT)
	return &restApiClient{
		url:        fmt.Sprintf("http://%s:%s/v1/event", restApiHost, restApiPort),
		httpClient: http.DefaultClient,
	}
}

func (c *restApiClient) AddTestEvent(eventRef *model.TestEvent) (string, error) {
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

	log.Printf("body: '%s'", string(body))

	respEventRef := new(model.TestEvent)
	err = json.Unmarshal(body, respEventRef)
	if err != nil {
		return "", err
	}
	return respEventRef.ID, nil
}
