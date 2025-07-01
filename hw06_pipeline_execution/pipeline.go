package hw06pipelineexecution

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

	stageLen := len(stages)
	stageChans := make([]Bi, stageLen+1)
	for i := range stageChans {
		stageChans[i] = make(Bi)
	}

	//wg := sync.WaitGroup{}
	//wg.Add(1)
	go func() {
		//defer wg.Done()
		defer close(stageChans[0])
		j := 0
		for v := range in {
			stageChans[0] <- v
			j = j + 1
		}
	}()

	for i, stage := range stages {
		var out = stage(stageChans[i])
		//wg.Add(1)
		go func() {
			//defer wg.Done()
			defer close(stageChans[i+1])
			for {
				select {
				case <-done:
					//wg.Add(1)
					go func() {
						//defer wg.Done()
						for range out {
						}
					}()
					return
				case o, ok := <-out:
					if ok {
						stageChans[i+1] <- o
					} else {
						return
					}
				}
			}
		}()
	}

	//go func() {
	//wg.Wait()
	//}()

	return stageChans[stageLen]

}
