package apiclient

import (
	"context"
	"time"

	"github.com/Ageres/hw-test/hw12_13_14_15_calendar/test/integration/model"
)

type ListPeriod string

const (
	DAY   ListPeriod = "day"
	WEEK  ListPeriod = "week"
	MONTH ListPeriod = "month"
)

type TestCalendarAPIClient interface {
	// return eventId, responseBody, error
	AddTestEvent(ctx context.Context, eventRef *model.TestEvent) (string, string, error)
	// return responseBody, error
	UpdateTestEvent(ctx context.Context, eventRef *model.TestEvent) (string, error)
	// return events, responseBody, error
	ListTestEvent(ctx context.Context, period ListPeriod, startDay time.Time) ([]model.TestEvent, string, error)
	Stop()
}
