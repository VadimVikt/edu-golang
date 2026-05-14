package hw05parallelexecution

import (
	"errors"
	"sync"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	if n <= 0 {
		return errors.New("n should be greater than 0")
	}
	if m <= 0 { // Возможно рефакторинг см Readme
		return ErrErrorsLimitExceeded
	}

	taskCh := make(chan Task, len(tasks))
	for _, t := range tasks {
		taskCh <- t
	}
	close(taskCh)
	var (
		wg       sync.WaitGroup
		errCount int
		errMu    sync.Mutex
		stopCh   = make(chan struct{})
		errOnce  sync.Once
	)
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				select {
				case <-stopCh:
					return
				case task, ok := <-taskCh:
					if !ok {
						return
					}
					if err := task(); err != nil {
						errMu.Lock()
						errCount++
						if errCount >= m {
							errMu.Unlock()
							errOnce.Do(func() {
								close(stopCh)
							})
							return
						}
						errMu.Unlock()
					}
				}
			}
		}()
	}
	wg.Wait()
	if errCount >= m {
		return ErrErrorsLimitExceeded
	}
	return nil
}
