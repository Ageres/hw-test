package integration

import "time"

type TestEvent struct {
	ID          string        `json:"id,omitempty"` // UUID
	Title       string        `json:"title" binding:"required"`
	StartTime   time.Time     `json:"startTime" binding:"required"`
	Duration    time.Duration `json:"duration" binding:"required"`
	Description string        `json:"description,omitempty"`
	UserID      string        `json:"userId" binding:"required"`
	Reminder    time.Duration `json:"reminder" binding:"required"`
}
