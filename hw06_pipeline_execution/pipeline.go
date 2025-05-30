package hw06pipelineexecution

import (
	"strconv"
	"strings"
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

	/*
		for inItem := range in {
			wg.Add(1)
			go func() {
				defer wg.Done()
				var sb strings.Builder
				for _, stage := range stages {
					outStage := stage(inItem)


					if _, isInt := outStageItem; isInt {
						outStageItem = strconv.Itoa(outStageItem)
					}
					sb.WriteString(outStageItem)
				}
			}()
		}
	*/

	for _, stage := range stages {
		wg.Add(1)
		go func() {
			defer wg.Done()
			outStage := stage(in)
			var sb strings.Builder
			for outStageItemIf := range outStage {
				//log.Printf("----201---- r: %v, type: %T", r, r)
				var outStageItemStr string
				if outStageItemInt, isInt := outStageItemIf.(int); isInt {
					outStageItemStr = strconv.Itoa(outStageItemInt)
				} else {
					outStageItemStr = outStageItemIf.(string)
				}
				sb.WriteString(outStageItemStr)
				/*
					ri, ok := outStageItemIf.(int)
					//log.Println("----202---- ok:", ok)
					if ok {
						//log.Printf("----203---- ok: true, r: %v, type_r: %T, ri: %v", r, r, ri)
						rs := strconv.Itoa(ri)
						outCh <- rs
					} else {
						//log.Printf("----204---- ok: false, r: %v, type_r: %T, ri: %v", r, r, ri)
						outCh <- outStageItemIf.(string)
					}
				*/
			}
			outCh <- sb.String()
		}()
	}
	go func() {
		wg.Wait()
		close(outCh)

	}()
	return outCh
}
