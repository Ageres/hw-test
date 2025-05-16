package hw05parallelexecution

import (
	"errors"
	"log"
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
	log.Println("len(tasks) =", len(tasks), ", n =", n, ", m =", m)

	if m <= 0 {
		return ErrErrorsLimitExceeded
	}

	resultCh := make(chan Result)
	endCh, errCountRef := countErr(resultCh, m)
	taskCh := getTaskChan(tasks, endCh)
	work(n, resultCh, endCh, taskCh)

	log.Println("errCount =", errCountRef.Load())

	if int(errCountRef.Load()) >= m {
		return ErrErrorsLimitExceeded
	}
	return nil
}

func countErr(resultCh <-chan Result, m int) (<-chan struct{}, *atomic.Int32) {
	endCh := make(chan struct{})
	var errCount atomic.Int32
	go func() {
		log.Println(">>> countErr start gourutine")
		defer log.Println("<<< countErr end gourutine")
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
		log.Println("!!! close endCh")
		close(endCh)
	}()
	return endCh, &errCount
}

func getTaskChan(tasks []Task, endCh <-chan struct{}) <-chan Task {
	taskCh := make(chan Task)
	go func() {
		log.Println(">>> getTaskChan start gourutine")
		defer close(taskCh)
		defer log.Println("!!! close taskCh")

		for i, task := range tasks {
			select {
			case taskCh <- task:
				log.Println("+++ getTaskChan process task", i)
			case <-endCh:
				log.Println("<<< getTaskChan end 1 gourutine")
				return
			}
		}
		log.Println("<<< getTaskChan end 2 gourutine")
	}()
	return taskCh
}

func work(n int, resultCh chan Result, endCh <-chan struct{}, taskCh <-chan Task) {
	wg := &sync.WaitGroup{}
	for j := range n {
		wg.Add(1)
		go func() {
			log.Println(">>> work start gourutine", j)
			defer wg.Done()
			for {
				task, ok := <-taskCh
				if ok {
					log.Println("--- work process task gourutine", j)
					err := task()
					select {
					case <-endCh:
						log.Println("<<< work end 1 gourutine", j)
						return
					case resultCh <- Result{
						Error: err,
					}:
					}
				} else {
					log.Println("<<< work end 2 gourutine", j)
					return
				}
			}
		}()
	}

	wg.Wait()

	log.Println("!!! close resultCh")
	close(resultCh)
}
