package hw05parallelexecution

import (
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

type Result struct {
	Error error
}

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
// m <= 0 - максимум 0 ошибок
func Run(tasks []Task, n, m int) error {
	fmt.Println("------------------100--------------------")
	fmt.Println("len(tasks):", len(tasks))
	fmt.Println("n:", n)
	fmt.Println("m:", m)

	if m <= 0 {
		return ErrErrorsLimitExceeded
	}

	taskCh := make(chan Task)
	var endCh chan struct{} = make(chan struct{})

	wg1 := &sync.WaitGroup{}

	wg1.Add(1)
	go func() {
		defer wg1.Done()
		for i, task := range tasks {
			fmt.Println("------------------200", i, "--------------------")
			select {
			case taskCh <- task:
				fmt.Println("------------------201--------------------")
			case <-endCh:
				return
			}
			fmt.Println("i:", i)
			fmt.Println("------------------202--------------------")
		}
		close(taskCh)
	}()

	go func() {
		wg1.Wait()
	}()

	fmt.Println("------------------300--------------------")

	resultCh := make(chan Result)

	for i := 0; i < n; i++ {
		wg1.Add(1)
		go func() {
			defer wg1.Done()
			select {
			case task := <-taskCh:
				err := task()
				resultCh <- Result{
					Error: err,
				}
			case <-endCh:
				close(resultCh)
				return
			}
		}()
	}

	var errCount atomic.Int32
	errCount.Store(0)

	wg1.Add(1)
	go func() {
		defer wg1.Done()
		for r := range resultCh {
			if r.Error != nil {
				errCount.Add(1)
			}
			c := errCount.Load()
			if int(c) >= m {
				close(endCh)
			}
		}
	}()

	if int(errCount.Load()) >= m {
		return ErrErrorsLimitExceeded
	}
	return nil
}
