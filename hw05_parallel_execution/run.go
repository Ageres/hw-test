package hw05parallelexecution

import (
	"errors"
	"fmt"
	"sync"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

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

	fmt.Println("------------------200--------------------")

	tasksize := min(n, len(tasks))
	fmt.Println("tasksize:", tasksize)

	taskCh := make(chan Task, tasksize)
	//defer close(taskCh)

	go func() {
		for i, task := range tasks {
			fmt.Println("------------------300--------------------")
			taskCh <- task
			fmt.Println("i:", i)
			fmt.Println("------------------301--------------------")
		}
	}()

	ch := make(chan error)
	//defer close(ch)

	go func() {
		for _, task := range tasks {
			taskCh <- task
		}
	}()

	wg := new(sync.WaitGroup)
	//wgCount := 0

	errCount := 0
	var out error
	wg.Add(1)
	go func() {
		defer wg.Done()
		j := 0
		for task := range taskCh {
			wg.Add(1)
			go func() {
				defer wg.Done()
				fmt.Println("------------------400--------------------")
				err := task()
				//if err != nil {
				ch <- err
				//}
				fmt.Println("err:", err)
				fmt.Println("j:", j)
				j++
				fmt.Println("------------------401--------------------")
			}()
			result, ok := <-ch
			if ok && result != nil {
				errCount++
			}
			if errCount == m {
				out = ErrErrorsLimitExceeded
				return
			}
		}

		/*
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
			} */
	}()

	/*
		go func() {
			wg.Wait()
		}()
	*/

	/*
		//errCount := 0
		for _ = range ch {
			errCount++
			if errCount == m {
				return ErrErrorsLimitExceeded
			}
		}
	*/

	wg.Wait()
	return out
}
