package hw05parallelexecution

import (
	"errors"
	"sync"
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

	tasksize := max(n, len(tasks))

	taskCh := make(chan Task, tasksize)

	ch := make(chan error)

	wg := new(sync.WaitGroup)
	wgCount := 0

	go func() {
		for {
			if wgCount <= tasksize {
				wg.Add(1)
				wgCount++
				go func() {
					defer wg.Done()
					task := <-taskCh
					err := task()
					if err == nil {
						ch <- err
					}
				}()
			}

		}
	}()

	errCount := 0
	for _ = range ch {
		errCount++
		if errCount == m {
			return ErrErrorsLimitExceeded
		}
	}

	return nil
}
