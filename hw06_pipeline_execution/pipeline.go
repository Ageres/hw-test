package hw06pipelineexecution

import (
	"log"
	"sync"
)

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	if len(stages) == 0 {
		return in
	}

	wg := sync.WaitGroup{}

	current := in
	for i, stage := range stages {
		stageInput := make(Bi)
		wg.Add(1)
		go func(input Bi, prev Out) {
			doneA := false
			doneB := false
			defer wg.Done()
			defer close(input)
			for {
				log.Println("------------101--------------:", "i = ", i, ", doneA = ", doneA, ", doneB = ", doneB)
				if doneA || doneB {
					//close(input)
					return
				}
				select {
				case <-done:
					log.Println("------------102--------------:", "i = ", i, ", done")
					doneA = true
					_, ok := <-input
					if ok {
						close(input)
					}
					return
				case v, ok := <-prev:
					if !ok {
						return
					}
					select {
					case <-done:
						doneB = true
						_, ok := <-input
						if ok {
							close(input)
						}
						return
					case input <- v:
					}
				}
				log.Println("------------103--------------:", "i = ", i, "doneA = ", doneA, ", doneB = ", doneB)
			}
			//defer close(stageInput)
		}(stageInput, current)
		current = stage(stageInput)

	}

	go func() {
		log.Println("------------201--------------")
		wg.Wait()
		log.Println("------------202--------------")
		//close(current)
	}()

	return current
}
