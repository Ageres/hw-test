package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	c "github.com/Ageres/hw-test/hw12_13_14_15_calendar/test/integration/client"
	"github.com/Ageres/hw-test/hw12_13_14_15_calendar/test/integration/config"
	"github.com/Ageres/hw-test/hw12_13_14_15_calendar/test/integration/model"
	"github.com/Ageres/hw-test/hw12_13_14_15_calendar/test/integration/utils"
)

type restAPIClient struct {
	url        string
	httpClient *http.Client
}

func NewRestAPIClient() c.TestCalendarAPIClient {
	restAPIHost := utils.GetEnvOrDefault(config.CalendarRestAPIHostEnv, config.CalendarRestAPIHostDefault)
	restAPIPort := utils.GetEnvOrDefault(config.CalendarRestAPIPortEnv, config.CalendarRestAPIPortDefault)
	return &restAPIClient{
		url:        fmt.Sprintf("http://%s:%s/v1/event", restAPIHost, restAPIPort),
		httpClient: http.DefaultClient,
	}
}

func (c *restAPIClient) AddTestEvent(ctx context.Context, eventRef *model.TestEvent) (string, string, error) {
	jsonBody, err := json.Marshal(eventRef)
	if err != nil {
		return "", "", err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.url, bytes.NewBuffer(jsonBody))
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
func (c *restAPIClient) UpdateTestEvent(ctx context.Context, eventRef *model.TestEvent) (string, error) {
	jsonBody, err := json.Marshal(eventRef)
	if err != nil {
		return "", err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPut, c.url, bytes.NewBuffer(jsonBody))
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
	Period   c.ListPeriod `json:"period,omitempty"`
	StartDay time.Time    `json:"startDay,omitempty"`
}

type ListTestEventResponseBody struct {
	Status string            `json:"status,omitempty"`
	Events []model.TestEvent `json:"events,omitempty"`
}

func (c *restAPIClient) ListTestEvent(
	ctx context.Context,
	period c.ListPeriod,
	startDay time.Time,
) ([]model.TestEvent, string, error) {
	reqBody := ListTestEventRequestBody{
		Period:   period,
		StartDay: startDay,
	}
	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return nil, "", err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.url, bytes.NewBuffer(jsonBody))
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
	for i := range result.Events {
		result.Events[i].StartTime = result.Events[i].StartTime.Local()
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
	if err != nil {
		return body, bodyStr, err
	}
	return body, bodyStr, nil
}

// Stop implements apiclient.TestCalendarApiClient.
func (c *restAPIClient) Stop() {
	panic("unimplemented")
}
