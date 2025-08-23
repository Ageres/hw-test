package model

import (
	"fmt"
	"time"
)

type ServiceName string

const CalendarServiceName ServiceName = "calendar"

type ServiceError struct {
	ServiceName ServiceName `json:"serviceName" binding:"required"`
	Status      int         `json:"status" binding:"required"`
	Message     string      `json:"message" binding:"required"`
	RequestID   string      `json:"requestId" binding:"required"`
	Timestamp   time.Time   `json:"timestamp" binding:"required"`
	Cause       error
}

func (s *ServiceError) Unwrap() error {
	return s.Cause
}

func (s *ServiceError) Error() string {
	return s.Message
}

func (s *ServiceError) ToJSON() string {
	if s.Cause == nil {
		return fmt.Sprintf(
			`{"serviceName":"%s","status":"%d","message":"%s","requestId":"%s","timestamp":"%v"}`,
			s.ServiceName, s.Status, s.Message, s.RequestID, s.Timestamp,
		)
	}
	return fmt.Sprintf(
		`{"serviceName":"%s","status":"%d","message":"%s","requestId":"%s","timestamp":"%v","cause":"%s"}`,
		s.ServiceName, s.Status, s.Message, s.RequestID, s.Timestamp, s.Cause.Error(),
	)
}

func NewCalendarServiceError(status int, message, requestId string, cause error) error {
	return NewCalendarServiceErrorAsIs(status, message, requestId, cause)
}

func NewCalendarServiceErrorAsIs(status int, message, requestId string, cause error) *ServiceError {
	return &ServiceError{
		ServiceName: CalendarServiceName,
		Status:      status,
		Message:     message,
		RequestID:   requestId,
		Timestamp:   time.Now(),
		Cause:       cause,
	}
}
