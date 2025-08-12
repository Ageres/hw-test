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

	out := stages[0](in)

	for _, stage := range stages[1:] {
		out = stage(out)
	}

	result := make(Bi)
	go func() {
		defer close(result)
		for {
			select {
			case <-done:
				go func() {
					for v := range out {
						_ = v
					}
				}()
				return
			case o, ok := <-out:
				if !ok {
					return
				}
				select {
				case <-done:
					return
				case result <- o:
				}
			}
		}
	}()

	return result
}
