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

	log.Println("------------------700 c=", errCountRef.Load(), "--------------------")

	if int(errCountRef.Load()) >= m {
		return ErrErrorsLimitExceeded
	}
	return nil
}

func work(n int, resultCh chan Result, endCh <-chan struct{}, taskCh <-chan Task) {
	wg := &sync.WaitGroup{}
	//log.Println("------------------300--------------------")

	for j := range n {
		wg.Add(1)
		go func() {
			log.Println(">>> work start gourutine", j)
			defer log.Println("<<< work end gourutine", j)
			defer wg.Done()
			for {
				select {
				case task, ok := <-taskCh:
					if ok {
						log.Println("--- work process task gourutine", j)
						err := task()
						/*
							if err != nil {
								log.Println("------------------401 j=", j, "err=", err, "--------------------")
							} else {
								log.Println("------------------402 j=", j, "err=nil--------------------")
							}
						*/
						select {
						case <-endCh:
							return
						case resultCh <- Result{
							Error: err,
						}:
						}

					} else {
						//log.Println(">>>---------------403-------------------- endCh return j=", j)
						return
					}

				case <-endCh:
					//log.Println(">>>---------------404-------------------- endCh return j=", j)
					return
				}
			}
		}()
	}

	wg.Wait()
	close(resultCh)
}

func countErr(resultCh <-chan Result, m int) (<-chan struct{}, *atomic.Int32) {
	endCh := make(chan struct{})
	var errCount atomic.Int32
	go func() {
		log.Println(">>>---countErr start gourutine--------------------")
		defer log.Println("<<<---countErr end gourutine--------------------")
		for r := range resultCh {
			err := r.Error

			if err != nil {
				//log.Println("------------------501 r.Error=", r.Error, "--------------------")
				errCount.Add(1)
			} else {
				//log.Println("------------------502 r.Error= nil --------------------")
			}
			c := errCount.Load()
			if int(c) >= m {
				break
			}
		}
		//log.Println("------------------503 c=", errCount.Load(), "--------------------")

		close(endCh)
		//log.Println("------------------504 end--------------------")
	}()
	return endCh, &errCount
}

func getTaskChan(tasks []Task, endCh <-chan struct{}) <-chan Task {
	taskCh := make(chan Task)
	go func() {
		log.Println(">>>---getTaskChan start gourutine--------------------")
		defer log.Println("<<<---getTaskChan end gourutine--------------------")
		defer close(taskCh)
		for _, task := range tasks {
			//log.Println("------------------200 i=", i, "--------------------")
			select {
			case taskCh <- task:
				//log.Println("------------------201-------------------- put task", i)
			case <-endCh:
				//log.Println("------------------202-------------------- return")
				return
			}
			//log.Println("------------------203--------------------")
		}
		//log.Println("------------------204 end--------------------")
	}()
	return taskCh
}
