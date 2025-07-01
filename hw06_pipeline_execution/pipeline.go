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

	go func() {
		defer close(stageChans[0])
		j := 0
		for v := range in {
			stageChans[0] <- v
			j++
		}
	}()

	for i, stage := range stages {
		var out = stage(stageChans[i])
		go func() {
			defer close(stageChans[i+1])
			for {
				select {
				case <-done:
					go func() {
						for _ = range out {
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

	return stageChans[stageLen]
}
