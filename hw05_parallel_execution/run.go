package hw05parallelexecution

import (
	"errors"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
// m < 0 - максимум 0 ошибок
func Run(tasks []Task, n, m int) error {
	// Place your code here.
	if m <= 0 {
		return ErrErrorsLimitExceeded
	}

	return nil
}
