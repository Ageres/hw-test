package apiclient

import (
	"time"

	"github.com/Ageres/hw-test/hw12_13_14_15_calendar/test/integration/model"
)

type ListPeriod string

const (
	DAY   ListPeriod = "day"
	WEEK  ListPeriod = "week"
	MONTH ListPeriod = "month"
)

type TestCalendarApiClient interface {
	AddTestEvent(eventRef *model.TestEvent) (string, string, error)                         // return eventId, responseBody, error
	UpdateTestEvent(eventRef *model.TestEvent) (string, error)                              // return responseBody, error
	ListTestEvent(period ListPeriod, startDay time.Time) ([]model.TestEvent, string, error) // return events, responseBody, error
}
