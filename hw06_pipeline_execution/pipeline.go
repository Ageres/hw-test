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

//================================================================================

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	if len(stages) == 0 {
		return in
	}

	stageLen := len(stages)

	stageChans := make([]Bi, stageLen+1)
	for i := range stageChans {
		stageChans[i] = make(Bi)
	}

	wg := sync.WaitGroup{}

	wg.Add(1)
	go func() {
		defer wg.Done()
		defer close(stageChans[0])
		j := 0
		for v := range in {
			log.Printf("------------100--------------: j = %v, v = %v", j, v)
			stageChans[0] <- v
			log.Printf("------------101--------------: j = %v, v = %v", j, v)
			j = j + 1

		}
	}()

	for i, stage := range stages {
		log.Printf("------------200--------------: i = %v", i)
		wg.Add(1)

		var out = stage(stageChans[i])
		go func() {
			defer log.Printf("------------399--------------: i = %v, end", i)
			defer wg.Done()
			defer close(stageChans[i+1])
			log.Printf("------------301--------------: i = %v, start", i)

			for {
				log.Printf("------------302--------------: i = %v", i)
				select {
				case <-done:
					for range out {
					}
					log.Println("------------311--------------:", "i = ", i, ", done")
					return
				case o, ok := <-out:
					log.Printf("------------303--------------: i = %v", i)
					if ok {
						log.Printf("------------304--------------: ok, i = %v", i)
						stageChans[i+1] <- o
					} else {
						log.Printf("------------305--------------: close, i = %v", i)
						return
					}
				}
			}
		}()

	}

	go func() {
		log.Println("------------901-------------- wg.Wait start")
		wg.Wait()
		log.Println("------------902-------------- wg.Wait end")
	}()

	return stageChans[stageLen]

}

//================================================================================

func ExecutePipelineOld1(in In, done In, stages ...Stage) Out {
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
