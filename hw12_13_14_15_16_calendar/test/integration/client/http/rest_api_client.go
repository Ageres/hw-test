package integration

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

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

func (c *restApiClient) AddTestEvent(eventRef *model.TestEvent) (string, string, error) {
	jsonBody, err := json.Marshal(eventRef)
	if err != nil {
		return "", "", err
	}
	req, err := http.NewRequest(http.MethodPost, c.url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return "", "", err
	}

	resp, err := http.DefaultClient.Do(req)
	body, bodyStr, err := parseHTTPResponce(resp, err)
	if err != nil {
		return "", bodyStr, err
	}
	if resp.StatusCode != http.StatusOK {
		return "", bodyStr, fmt.Errorf("response status '%s'", resp.Status)
	}

	respEventRef := new(model.TestEvent)
	err = json.Unmarshal(body, respEventRef)
	if err != nil {
		return "", bodyStr, err
	}
	return respEventRef.ID, bodyStr, nil
}

// UpdateTestEvent implements apiclient.TestCalendarApiClient.
func (c *restApiClient) UpdateTestEvent(eventRef *model.TestEvent) (string, error) {
	jsonBody, err := json.Marshal(eventRef)
	if err != nil {
		return "", err
	}
	req, err := http.NewRequest(http.MethodPut, c.url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return "", err
	}
	resp, err := http.DefaultClient.Do(req)
	body, bodyStr, err := parseHTTPResponce(resp, err)
	if err != nil {
		return bodyStr, err
	}
	if resp.StatusCode != http.StatusOK {
		return bodyStr, fmt.Errorf("response status '%s'", resp.Status)
	}

	respEventRef := new(model.TestEvent)
	err = json.Unmarshal(body, respEventRef)
	if err != nil {
		return bodyStr, err
	}
	return bodyStr, nil
}

type ListTestEventRequestBody struct {
	Period    c.ListPeriod
	StartDate time.Time
}

type ListTestEventResponseBody struct {
	Status string
	Events []model.TestEvent
}

func (c *restApiClient) ListTestEvent(period c.ListPeriod, startDate time.Time) ([]model.TestEvent, string, error) {
	reqBody := ListTestEventRequestBody{
		Period:    period,
		StartDate: startDate,
	}
	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return nil, "", err
	}
	req, err := http.NewRequest(http.MethodGet, c.url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, "", err
	}
	resp, err := http.DefaultClient.Do(req)
	body, bodyStr, err := parseHTTPResponce(resp, err)
	if err != nil {
		return nil, bodyStr, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, bodyStr, fmt.Errorf("response status '%s'", resp.Status)
	}

	var result ListTestEventResponseBody
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, bodyStr, err
	}
	return result.Events, bodyStr, nil
}

func parseHTTPResponce(resp *http.Response, err error) ([]byte, string, error) {
	var body []byte
	var bodyStr string
	var parseRespErr error
	if resp != nil {
		defer resp.Body.Close()
		body, parseRespErr = io.ReadAll(resp.Body)
		if parseRespErr != nil {
			if err != nil {
				return nil, "", fmt.Errorf("parse response error '%w', response error '%w'", parseRespErr, err)
			}
			return nil, "", fmt.Errorf("parse response error '%w'", parseRespErr)
		}
		bodyStr = string(body)
	}
	log.Printf("body: '%s'", string(body))
	if err != nil {
		return body, bodyStr, err
	}
	return body, bodyStr, nil
}
