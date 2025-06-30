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
	//current := make(Bi)
	countGor := 0
	for i, stage := range stages {
		countGor = countGor + 1
		stageInput := make(Bi)
		wg.Add(1)
		go func(input Bi, prev Out) {
			defer log.Println("------------105--------------:", "i = ", i, ", end")
			defer wg.Done()
			defer close(input)
			for {
				//log.Println("------------101--------------:", "i = ", i, ", doneA = ", doneA, ", doneB = ", doneB)
				select {
				case <-done:
					log.Println("------------102--------------:", "i = ", i, ", done")
					return
				case v, ok := <-prev:
					if !ok {
						return
					}
					select {
					case <-done:
						log.Println("------------103--------------:", "i = ", i, ", done")
						return
					case input <- v:
					}
				}
				log.Println("------------104--------------:", "i = ", i)
			}
			//
		}(stageInput, current)
		current = stage(stageInput)
		log.Println("------------106--------------:", "i = ", i)
	}

	go func() {
		log.Println("------------201--------------: countGor: ", countGor)
		wg.Wait()
		//close(current)
		log.Println("------------202--------------: countGor: ", countGor)
	}()

	return current
}
