package hw06pipelineexecution

import "log"

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	// Place your code here.
	currentOut := in
	for i, stage := range stages {
		log.Printf("Stage %v processing", i)
		nexOut := make(Bi)
		go func(s Stage, in In, out Bi) {
			defer close(out)
			for val := range in {
				select {
				case <-done:
					return
				case out <- val:
				}
			}
		}(stage, currentOut, nexOut)

		currentOut = nexOut

	}
	return currentOut
}
