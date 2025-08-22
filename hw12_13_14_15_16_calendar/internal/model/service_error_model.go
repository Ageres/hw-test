package model

import (
	"fmt"
	"net/http"
	"strings"
	"time"
)

type ServiceError struct {
	ServiceName string    `json:"serviceName" binding:"required"`
	Status      int       `json:"status" binding:"required"`
	Message     string    `json:"message" binding:"required"`
	RequestID   string    `json:"requestId" binding:"required"`
	Timestamp   time.Time `json:"timestamp" binding:"required"`
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
		ServiceName: "calendar",
		Status:      status,
		Message:     message,
		RequestID:   requestId,
		Timestamp:   time.Now(),
		Cause:       cause,
	}
}

func DefineStatusCode(errMsg string) int {
	if strings.Contains(errMsg, "user is not the owner of the event, conflict with") || strings.Contains(errMsg, "time is already taken by another event") {
		return http.StatusConflict
	}
	if strings.Contains(errMsg, "event not found") {
		return http.StatusNotFound
	}
	if strings.Contains(errMsg, "event is nil") ||
		strings.Contains(errMsg, "failed to validate event id") ||
		strings.Contains(errMsg, "title is empty") ||
		strings.Contains(errMsg, "event time is expired") ||
		strings.Contains(errMsg, "duration must be positive") ||
		strings.Contains(errMsg, "user id is empty") {
		return http.StatusBadRequest
	}
	return http.StatusInternalServerError
}
