package apiclient

import "github.com/Ageres/hw-test/hw12_13_14_15_calendar/test/integration/model"

type TestCalendarApiClient interface {
	AddTestEvent(eventRef *model.TestEvent) (string, string, error) // eventId, responseBody, error
	UpdateTestEvent(eventRef *model.TestEvent) (string, error)      // responseBody, error
}
