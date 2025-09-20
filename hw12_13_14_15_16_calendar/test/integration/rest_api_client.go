package integration

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
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

	log.Printf("body: '%s'", string(body))

	respEventRef := new(TestEvent)
	err = json.Unmarshal(body, respEventRef)
	if err != nil {
		return "", err
	}
	return respEventRef.ID, nil
}
