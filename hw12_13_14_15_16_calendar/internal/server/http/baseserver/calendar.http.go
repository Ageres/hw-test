package httpserverinterface

import (
	"context"
	"time"

	"github.com/Ageres/hw-test/hw12_13_14_15_calendar/internal/storage"
)

type HTTPServer interface {
	Start(ctx context.Context) error
	Stop(ctx context.Context) error
}

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

type GetEventListResponse struct {
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
	ID string `json:"id" binding:"required"`
}

type DeleteEventStatus string

const (
	DELETE DeleteEventStatus = "Event deleted successfully"
)

type DeleteEventResponse struct {
	Status DeleteEventStatus `json:"status" binding:"required"`
}
