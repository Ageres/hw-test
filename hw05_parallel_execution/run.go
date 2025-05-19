package hw05parallelexecution

import (
	"errors"
	"sync"
	"sync/atomic"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

type Result struct {
	Error error
}

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
// m <= 0 - максимум 0 ошибок.
func Run(tasks []Task, n, m int) error {
	if m <= 0 {
		return ErrErrorsLimitExceeded
	}

	resultCh := make(chan Result)
	endCh, errCountRef := countErr(resultCh, m)
	taskCh := getTaskChan(tasks, endCh)
	work(n, resultCh, endCh, taskCh)

	if int(errCountRef.Load()) >= m {
		return ErrErrorsLimitExceeded
	}
	return nil
}

func countErr(resultCh <-chan Result, m int) (<-chan struct{}, *atomic.Int32) {
	endCh := make(chan struct{})
	var errCount atomic.Int32
	go func() {
		for r := range resultCh {
			err := r.Error
			if err != nil {
				errCount.Add(1)
			}
			c := errCount.Load()
			if int(c) >= m {
				break
			}
		}
		close(endCh)
	}()
	return endCh, &errCount
}

func getTaskChan(tasks []Task, endCh <-chan struct{}) <-chan Task {
	taskCh := make(chan Task)
	go func() {
		defer close(taskCh)
		for _, task := range tasks {
			select {
			case taskCh <- task:
			case <-endCh:
				return
			}
		}
	}()
	return taskCh
}

func work(n int, resultCh chan Result, endCh <-chan struct{}, taskCh <-chan Task) {
	wg := &sync.WaitGroup{}
	for range n {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				task, ok := <-taskCh
				if ok {
					err := task()
					select {
					case <-endCh:
						return
					case resultCh <- Result{
						Error: err,
					}:
					}
				} else {
					return
				}
			}
		}()
	}
	wg.Wait()
	close(resultCh)
}
