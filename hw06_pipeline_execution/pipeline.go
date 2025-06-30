package hw06pipelineexecution

import "sync"

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
	for _, stage := range stages {
		stageInput := make(Bi)
		wg.Add(1)
		go func(input Bi, prev Out) {
			defer wg.Done()
			defer close(input)
			for {
				select {
				case <-done:
					return
				case v, ok := <-prev:
					if !ok {
						return
					}
					select {
					case <-done:
						return
					case input <- v:
					}
				}
			}
		}(stageInput, current)
		current = stage(stageInput)
	}

	go func() {
		wg.Wait()
	}()

	return current
}
