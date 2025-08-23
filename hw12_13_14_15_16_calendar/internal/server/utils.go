package server

import (
	"net/http"
	"strings"
)

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
