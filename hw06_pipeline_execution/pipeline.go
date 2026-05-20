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
	current := in
	for _, stage := range stages {
		// Создаем выходной канал для текущей стадии
		out := make(Bi)
		// Запускаем стадию в горутине
		go func(s Stage, in In, out Bi) {
			defer close(out)
			for val := range s(orDone(in, done)) {
				select {
				case out <- val:
				case <-done:
					return
				}
			}
		}(stage, current, out)
		current = out
	}
	return current
}

// orDone оборачивает канал, чтобы он закрывался при done
func orDone(in In, done In) In {
	out := make(Bi)
	go func() {
		defer close(out)
		for {
			select {
			case <-done:
				return
			case val, ok := <-in:
				if !ok {
					return
				}
				select {
				case out <- val:
				case <-done:
					return
				}
			}
		}
	}()
	return out
}
