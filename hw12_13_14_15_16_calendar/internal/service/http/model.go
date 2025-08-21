package httpservice

import (
	"context"
	"time"

	"github.com/Ageres/hw-test/hw12_13_14_15_calendar/internal/storage"
	"github.com/Ageres/hw-test/hw12_13_14_15_calendar/internal/utils"
)

// ---------------------------------------------------------
// get event list models
type GetEventListPeriod string

const (
	DAY   GetEventListPeriod = "day"
	WEEK  GetEventListPeriod = "week"
	MONTH GetEventListPeriod = "month"
)

type GetEventListRequest struct {
	Period   GetEventListPeriod `json:"period" binding:"required"`
	StartDay *time.Time         `json:"startDay" binding:"required"`
}

type GetEventListStatus string

const (
	LISTDAY   GetEventListStatus = "Day event list successfully retrieved"
	LISTWEEK  GetEventListStatus = "Week event list successfully retrieved"
	LISTMONTH GetEventListStatus = "Month event list successfully retrieved"
)

type GetListResponse struct {
	Status GetEventListStatus `json:"status" binding:"required"`
	Events []storage.Event    `json:"events,omitempty"`
}

// ---------------------------------------------------------
// post event models (add event)
type AddEventRequest storage.Event

type AddEventStatus string

const (
	ADD AddEventStatus = "Event added successfully"
)

type AddEventResponse struct {
	Status AddEventStatus `json:"status" binding:"required"`
	Event  *storage.Event `json:"events,omitempty"`
}

// ---------------------------------------------------------
// put event models (update event)
type UpdateEventRequest storage.Event

type UpdateEventStatus string

const (
	UPDATE UpdateEventStatus = "Event updated successfully"
)

type UpdateEventResponse struct {
	Status UpdateEventStatus `json:"status" binding:"required"`
}

// ---------------------------------------------------------
// delete event models
type DeleteEventRequest struct {
	Id string `json:"id" binding:"required"`
}

type DeleteEventStatus string

const (
	Delete DeleteEventStatus = "Event deleted successfully"
)

type DeleteEventResponse struct {
	Status DeleteEventStatus `json:"status" binding:"required"`
}

// ---------------------------------------------------------
// error models

type ServiceName string

const CalendarServiceName ServiceName = "calendar"

type HttpError struct {
	ServiceName `json:"serviceName" binding:"required"`
	Message     string    `json:"message" binding:"required"`
	RequestID   string    `json:"requestId" binding:"required"`
	Timestamp   time.Time `json:"timestamp" binding:"required"`
}

func (he *HttpError) Error() string {
	return he.Message
}

func NewHttpError(ctx context.Context, message string) error {
	return &HttpError{
		ServiceName: CalendarServiceName,
		Message:     message,
		RequestID:   utils.GetRequestID(ctx),
		Timestamp:   time.Now(),
	}
}
