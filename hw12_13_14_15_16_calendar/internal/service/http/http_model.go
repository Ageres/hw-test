package httpservice

import (
	"time"

	"github.com/Ageres/hw-test/hw12_13_14_15_calendar/internal/storage"
)

// ---------------------------------------------------------
// get event  list model
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
