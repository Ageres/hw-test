package hw06pipelineexecution

import (
	"log"
	"strconv"
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
				log.Printf("----201---- r: %v, type: %T", r, r)

				ri, ok := r.(int)
				log.Println("----202---- ok:", ok)
				if ok {
					log.Printf("----203---- ok: true, r: %v, type_r: %T, ri: %v", r, r, ri)
					rs := strconv.Itoa(ri)
					outCh <- rs
				} else {
					log.Printf("----204---- ok: false, r: %v, type_r: %T, ri: %v", r, r, ri)
					outCh <- r.(string)
				}

				log.Println("----205---- ri:", ri)
				//
				//outCh <- ri
			}
		}()
	}
	go func() {
		wg.Wait()
	}()
	return outCh
}
