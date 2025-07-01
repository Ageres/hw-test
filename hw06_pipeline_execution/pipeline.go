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

	/*
		in := make(Bi)
		for v := range inn {
			in <- v
		}
	*/

	wg := sync.WaitGroup{}

	current := in
	//current := make(Out)
	countGor := 0
	for i, stage := range stages {
		countGor = countGor + 1
		stageInput := make(Bi)
		wg.Add(1)
		go func(input Bi, prev Out) {
			//doneS := false
			defer log.Println("------------105--------------:", "i = ", i, ", end")
			defer wg.Done()
			//defer closeCurrent(current, doneS, &wg)
			defer close(input)
			//defer closeInput(input, doneS)
			//defer closePrev(doneS, &wg, prev)

			for {
				//log.Println("------------101--------------:", "i = ", i, ", doneA = ", doneA, ", doneB = ", doneB)
				select {
				case <-done:
					log.Println("------------102--------------:", "i = ", i, ", done")

					/*
						go func() {
							for range current {
							}
						}()
					*/

					//
					for range prev {
						//stage(prev)
					}

					//doneS = true
					/*
						wg.Add(1)
						go func() {
							defer wg.Done()
							//for range prev {}

							countPrev := 0

							timer := time.NewTimer(1 * time.Second)

							for {
								countPrev = countPrev + 1
								select {
								case v, ok := <-prev:
									log.Println("------------1021-------------: prev  ", "i = ", i, ", countPrev =", countPrev)
									if !ok {
										return
									} else {
										input <- v
									}
								case <-timer.C:
									log.Println("------------1021-------------: timer ", "i = ", i, ", countPrev =", countPrev)
									return
								}
							}

							/*
								for range prev {
									countPrev = countPrev + 1
									log.Println("------------1021-------------:", "i = ", i, ", countPrev =", countPrev)
								}
					*/
					//}()

					//<-prev
					return
				case v, ok := <-prev:
					if !ok {
						return
					}
					select {
					case <-done:
						log.Println("------------103--------------:", "i = ", i, ", done")
						//for range input {}
						return
					case input <- v:
						log.Println("------------104--------------:", "i = ", i, ", v = ", v, ", prev")
					}
				}
				//log.Println("------------104--------------:", "i = ", i)
			}
			//
		}(stageInput, current)
		current = stage(stageInput)
		//log.Println("------------106--------------:", "i = ", i)
	}

	go func() {
		log.Println("------------201--------------: countGor: ", countGor)
		wg.Wait()
		log.Println("------------202--------------: countGor: ", countGor)

		for v := range current {
			log.Println("------------203--------------: v: ", v)
		}
		/*
			for {
				_, ok := <-current
				log.Println("------------203--------------: ok: ", ok)
				if !ok {
					break
				}
			}
		*/
		//log.Println("------------204--------------: countGor: ", countGor)
	}()

	wg2 := sync.WaitGroup{}
	wg2.Add(1)
	go func() {
		defer wg2.Done()
		for {
			//log.Println("------------205--------------")
			select {
			case <-done:
				//log.Println("------------206-------------- done")
				for range current {
				}
				return
				/*
					case _, ok := <-in:
						if !ok {
							return
						} else {
							in <- v
						}*/

			}
		}
	}()

	go func() {
		log.Println("------------301--------------: countGor: ", countGor)
		wg2.Wait()
		log.Println("------------302--------------: countGor: ", countGor)
	}()

	return current
}

func closeInput(input Bi, dones bool) {
	if dones {
		for range input {
		}
	}
	close(input)
}

func closeCurrent(current In, dones bool, wg *sync.WaitGroup) {
	wg.Add(1)
	go func() {
		defer wg.Done()
		if dones {
			for range current {
			}
		}
	}()

}

/*
func closePrev(doneS bool, wg *sync.WaitGroup, prev Out) {
	if doneS {
		wg.Add(1)
		go func() {
			defer wg.Done()
			//countPrev := 0
			for range prev {
				//countPrev = countPrev + 1
				//log.Println("------------1021-------------:", "i = ", i, ", countPrev =", countPrev)
			}
		}()
	}
}
*/
