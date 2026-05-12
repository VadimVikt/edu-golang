package hw05parallelexecution

import (
	"errors"
	"fmt"
	"sync"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	//fmt.Printf("Количество горутин до создания: %d\n", runtime.NumGoroutine())
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
					fmt.Println("Превышено кол во ошибков")
					return
				case task, ok := <-taskCh:
					//fmt.Println("Задача прочитана из канала")
					if !ok {
						return
					}
					if err := task(); err != nil {
						errMu.Lock()
						errCount++
						fmt.Println("Количество ошибков - ", errCount)
						if errCount >= m {
							errMu.Unlock()
							errOnce.Do(func() {
								fmt.Println("Закрыли stopCh")
								close(stopCh)
							})
							return
						}
						errMu.Unlock()
					}
				}

			}
		}()
		//fmt.Printf("Количество горутин после создания: %d\n", runtime.NumGoroutine())
	}
	wg.Wait()
	//fmt.Printf("Количество горутин после завершения: %d\n", runtime.NumGoroutine())
	if errCount >= m {
		return ErrErrorsLimitExceeded
	}
	return nil
}
