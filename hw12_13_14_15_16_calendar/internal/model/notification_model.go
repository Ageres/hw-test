package model

import "time"

type Notification struct {
	ID     string    `json:"id" binding:"required"`
	Title  string    `json:"title" binding:"required"`
	Date   time.Time `json:"date" binding:"required"`
	UserID string    `json:"userId" binding:"required"`
}
