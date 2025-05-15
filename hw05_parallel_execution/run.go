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
	log.Println("------------------100--------------------")
	log.Println("len(tasks):", len(tasks))
	log.Println("n:", n)
	log.Println("m:", m)

	if m <= 0 {
		return ErrErrorsLimitExceeded
	}

	taskCh := make(chan Task)
	var endCh chan struct{} = make(chan struct{})

	wg1 := &sync.WaitGroup{}

	//wg1.Add(1)
	go func() {
		defer close(taskCh)
		//defer wg1.Done()
		for i, task := range tasks {
			//for {
			log.Println("------------------200 i=", i, "--------------------")
			select {
			case taskCh <- task:
				log.Println("------------------201-------------------- put task", i)
			case <-endCh:
				log.Println("------------------202-------------------- return")
				//close(taskCh)
				//continue
				//break
				return
			}
			//log.Println("i:", i)
			log.Println("------------------203--------------------")
			//}
		}
		log.Println("------------------204 end--------------------")
	}()

	/*
		go func() {
			wg1.Wait()
		}()
	*/

	log.Println("------------------300--------------------")

	resultCh := make(chan Result)

	for j := range n {
		log.Println("------------------400 j=", j, "--------------------")
		wg1.Add(1)
		go func() {
			defer wg1.Done()
			for {
				select {
				case task := <-taskCh:

					err := task()
					log.Println("---------")
					log.Println(err)
					log.Println("---------")
					if err != nil {
						log.Println("------------------401 j=", j, "err=", err, "--------------------")
					} else {
						log.Println("------------------402 j=", j, "err=nil--------------------")
					}

					resultCh <- Result{
						Error: err,
					}
				case <-endCh:
					log.Println("------------------403-------------------- endCh return")
					//close(resultCh)
					return
				}
			}
			log.Println("------------------404 end--------------------")
		}()
	}

	var errCount atomic.Int32
	errCount.Store(0)

	wg1.Add(1)
	go func() {
		defer wg1.Done()
		for r := range resultCh {
			err := r.Error

			if err != nil {
				log.Println("------------------501 r.Error=", r.Error, "--------------------")
				errCount.Add(1)
			} else {
				log.Println("------------------502 r.Error= nil --------------------")
			}
			c := errCount.Load()
			if int(c) >= m {
				break
			}
		}
		log.Println("------------------503 c=", errCount.Load(), "--------------------")

		close(endCh)
		log.Println("------------------504 end--------------------")
	}()

	wg1.Wait()
	close(resultCh)

	log.Println("------------------600 c=", errCount.Load(), "--------------------")

	if int(errCount.Load()) >= m {
		return ErrErrorsLimitExceeded
	}
	return nil
}

//func getTaskChan()
