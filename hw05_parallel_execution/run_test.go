package hw05parallelexecution

import (
	"errors"
	"fmt"
	"math/rand"
	"runtime"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.uber.org/goleak"
)

func TestRun(t *testing.T) {
	defer goleak.VerifyNone(t)

	t.Run("if were errors in first M tasks, than finished not more N+M tasks", func(t *testing.T) {
		tasksCount := 50
		tasks := make([]Task, 0, tasksCount)

		var runTasksCount int32

		for i := 0; i < tasksCount; i++ {
			err := fmt.Errorf("error from task %d", i)
			tasks = append(tasks, func() error {
				time.Sleep(time.Millisecond * time.Duration(rand.Intn(100)))
				atomic.AddInt32(&runTasksCount, 1)
				return err
			})
		}

		workersCount := 10
		maxErrorsCount := 23
		err := Run(tasks, workersCount, maxErrorsCount)

		require.Truef(t, errors.Is(err, ErrErrorsLimitExceeded), "actual err - %v", err)
		require.LessOrEqual(t, runTasksCount, int32(workersCount+maxErrorsCount), "extra tasks were started")
	})

	t.Run("tasks without errors", func(t *testing.T) {
		tasksCount := 50
		tasks := make([]Task, 0, tasksCount)

		var runTasksCount int32
		var sumTime time.Duration

		for i := 0; i < tasksCount; i++ {
			taskSleep := time.Millisecond * time.Duration(rand.Intn(100))
			sumTime += taskSleep

			tasks = append(tasks, func() error {
				time.Sleep(taskSleep)
				atomic.AddInt32(&runTasksCount, 1)
				return nil
			})
		}

		workersCount := 5
		maxErrorsCount := 1

		start := time.Now()
		err := Run(tasks, workersCount, maxErrorsCount)
		elapsedTime := time.Since(start)
		require.NoError(t, err)

		require.Equal(t, int32(tasksCount), runTasksCount, "not all tasks were completed")
		require.LessOrEqual(t, int64(elapsedTime), int64(sumTime/2), "tasks were run sequentially?")
	})

	t.Run("tasks with errors without time.Sleep", func(t *testing.T) {
		tasksCount := 50
		workersCount := 5
		maxErrorsCount := 1

		var (
			runTasksCount     int32
			maxConcurrent     int32
			currentConcurrent int32
			mu                sync.Mutex
		)
		tasks := make([]Task, 0, tasksCount)
		for i := 0; i < tasksCount; i++ {
			done := make(chan struct{})
			close(done)
			tasks = append(tasks, func() error {
				atomic.AddInt32(&currentConcurrent, 1)
				mu.Lock()
				if currentConcurrent > maxConcurrent {
					maxConcurrent = currentConcurrent
				}
				mu.Unlock()
				<-done
				atomic.AddInt32(&currentConcurrent, 1)
				atomic.AddInt32(&runTasksCount, 1)
				return nil
			})
		}
		start := time.Now()
		err := Run(tasks, workersCount, maxErrorsCount)
		elapsedTime := time.Since(start)
		require.NoError(t, err)
		require.Equal(t, int32(tasksCount), runTasksCount, "не все задачи были выполнены")
		require.GreaterOrEqual(t, maxConcurrent, int32(workersCount), "задачи выполнялись не параллельно")
		require.Less(t, elapsedTime, 100*time.Millisecond, "задачи выполнялись слишком долго (последовательно?)")

	})
}

func TestMyTask(t *testing.T) {
	tasks := []Task{
		func() error { fmt.Println("task 1 ---ok"); return nil },
		func() error { fmt.Println("task 2 ---err"); return errors.New("err") },
		func() error { fmt.Println("task 3 ---ok"); return nil },
		func() error { fmt.Println("task 4 ---err"); return errors.New("err") },
		func() error { fmt.Println("task 5 ---err"); return errors.New("err") },
		func() error { fmt.Println("task 6 ---ok"); return nil },
	}
	workersCount := 4
	maxErrorsCount := 4
	err := Run(tasks, workersCount, maxErrorsCount)
	_ = err
	fmt.Printf("Количество горутин после завершения теста: %d\n", runtime.NumGoroutine())
}

func TestNoError(t *testing.T) {
	tasks := []Task{
		func() error { fmt.Println("task 1 ---ok"); return nil },
		func() error { fmt.Println("task 2 ---err"); return errors.New("err") },
	}
	workersCount := 2
	maxErrorsCount := 0
	err := Run(tasks, workersCount, maxErrorsCount)
	require.Equal(t, err, ErrErrorsLimitExceeded)
}
