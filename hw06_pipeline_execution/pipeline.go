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
	// Place your code here.
	outCh := make(chan interface{})
	wg := sync.WaitGroup{}
	for _, stage := range stages {
		wg.Add(1)
		go func() {
			defer wg.Done()
			o := stage(in)
			for r := range o {
				log.Println("----201---- r:", r)
				//ri := r.(int)
				//rs := strconv.Itoa(ri)
				outCh <- r
			}
		}()
	}
	go func() {
		wg.Wait()
	}()
	return outCh
}
