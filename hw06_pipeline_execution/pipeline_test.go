package hw06pipelineexecution

import (
	"fmt"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

const (
	sleepPerStage = time.Millisecond * 100
	fault         = sleepPerStage / 2
)

func TestPipeline(t *testing.T) {
	//t.Parallel()
	// Stage generator
	//g := func(name string, f func(v interface{}) interface{}) Stage {
	//	return func(in In) Out {
	//		out := make(Bi)
	//		go func() {
	//			defer close(out)
	//			for v := range in {
	//				fmt.Println("Что то неопнятное ", v)
	//				time.Sleep(sleepPerStage)
	//				out <- f(v)
	//			}
	//		}()
	//		return out
	//	}
	//}
	g := func(name string, f func(v interface{}) interface{}) Stage {
		return func(in In) Out {
			out := make(Bi)
			go func() {
				defer close(out)
				for v := range in {
					fmt.Printf("Stage %s: received %T = %v\n", name, v, v)
					time.Sleep(sleepPerStage)
					result := f(v)
					fmt.Printf("Stage %s: sending %T = %v\n", name, result, result)
					out <- result
				}
			}()
			return out
		}
	}
	//g := func(name string, f func(v interface{}) interface{}) Stage {
	//	return func(in In) Out {
	//		out := make(Bi)
	//		go func() {
	//			defer close(out)
	//			for v := range in {
	//				fmt.Printf("Stage '%s': received %T = %v\n", name, v, v)
	//				time.Sleep(sleepPerStage)
	//				result := f(v)
	//				fmt.Printf("Stage '%s': sending %T = %v\n", name, result, result)
	//				out <- result
	//			}
	//		}()
	//		return out
	//	}
	//}
	stages := []Stage{
		g("Dummy", func(v interface{}) interface{} { return v }),
		g("Multiplier (* 2)", func(v interface{}) interface{} { return v.(int) * 2 }),
		g("Adder (+ 100)", func(v interface{}) interface{} { return v.(int) + 100 }),
		g("Stringifier", func(v interface{}) interface{} { return strconv.Itoa(v.(int)) }),
	}

	t.Run("simple case", func(t *testing.T) {
		t.Parallel()
		in := make(Bi)
		data := []int{1, 2, 3, 4, 5}

		go func() {
			for _, v := range data {
				in <- v
			}
			close(in)
		}()

		result := make([]string, 0, 10)
		start := time.Now()
		for s := range ExecutePipeline(in, nil, stages...) {
			// Безопасное преобразование с проверкой типа
			str, ok := s.(string)
			if !ok {
				t.Fatalf("Expected string, got %T: %v", s, s)
			}
			fmt.Println("что попало - ", s)
			result = append(result, str)
		}

		elapsed := time.Since(start)

		require.Equal(t, []string{"102", "104", "106", "108", "110"}, result)
		require.Less(t,
			int64(elapsed),
			// ~0.8s for processing 5 values in 4 stages (100ms every) concurrently
			int64(sleepPerStage)*int64(len(stages)+len(data)-1)+int64(fault))
	})

	t.Run("done case", func(t *testing.T) {
		t.Parallel()
		in := make(Bi)
		done := make(Bi)
		data := []int{1, 2, 3, 4, 5}

		// Abort after 200ms
		abortDur := sleepPerStage * 2
		go func() {
			<-time.After(abortDur)
			close(done)
		}()

		go func() {
			for _, v := range data {
				in <- v
			}
			close(in)
		}()

		result := make([]string, 0, 10)
		start := time.Now()
		for s := range ExecutePipeline(in, done, stages...) {
			result = append(result, s.(string))
		}
		elapsed := time.Since(start)

		require.Len(t, result, 0)
		require.Less(t, int64(elapsed), int64(abortDur)+int64(fault))
	})
}

func TestAllStageStop(t *testing.T) {
	t.Parallel()
	wg := sync.WaitGroup{}
	// Stage generator
	g := func(_ string, f func(v interface{}) interface{}) Stage {
		return func(in In) Out {
			out := make(Bi)
			wg.Add(1)
			go func() {
				defer wg.Done()
				defer close(out)
				for v := range in {
					time.Sleep(sleepPerStage)
					out <- f(v)
				}
			}()
			return out
		}
	}

	stages := []Stage{
		g("Dummy", func(v interface{}) interface{} { return v }),
		g("Multiplier (* 2)", func(v interface{}) interface{} { return v.(int) * 2 }),
		g("Adder (+ 100)", func(v interface{}) interface{} { return v.(int) + 100 }),
		g("Stringifier", func(v interface{}) interface{} { return strconv.Itoa(v.(int)) }),
	}

	t.Run("done case", func(t *testing.T) {
		t.Parallel()
		in := make(Bi)
		done := make(Bi)
		data := []int{1, 2, 3, 4, 5}

		// Abort after 200ms
		abortDur := sleepPerStage * 2
		go func() {
			<-time.After(abortDur)
			close(done)
		}()

		go func() {
			for _, v := range data {
				in <- v
			}
			close(in)
		}()

		result := make([]string, 0, 10)
		for s := range ExecutePipeline(in, done, stages...) {
			result = append(result, s.(string))
		}
		wg.Wait()

		require.Len(t, result, 0)

	})
}

func TestStages(t *testing.T) {
	// Stage generator
	g := func(name string, f func(v interface{}) interface{}) Stage {
		return func(in In) Out {
			out := make(Bi)
			go func() {
				defer close(out)
				for v := range in {
					time.Sleep(sleepPerStage)
					out <- f(v)
				}
			}()
			return out
		}
	}
	//stages := []Stage{
	//g("Dummy", func(v interface{}) interface{} {
	//	switch par := v.(type) {
	//	case int:
	//		return par
	//	case string:
	//		return par
	//	default:
	//		return v
	//	}
	//}),
	//g("Multiplier (* 2)", func(v interface{}) interface{} {
	//	return v.(int) * 2
	//}),
	//g("Adder (+ 100)", func(v interface{}) interface{} {
	//	time.Sleep(time.Second)
	//	fmt.Println("Печатаем - ", v.(int)+100)
	//	return v.(int) + 100
	//}),
	//g("Stringifier", func(v interface{}) interface{} { return strconv.Itoa(v.(int)) }),
	//}
	stages := []Stage{
		g("Dummy", func(v interface{}) interface{} { return v }),
		g("Multiplier (* 2)", func(v interface{}) interface{} { return v.(int) * 2 }),
		g("Adder (+ 100)", func(v interface{}) interface{} { return v.(int) + 100 }),
		g("Stringifier", func(v interface{}) interface{} { return strconv.Itoa(v.(int)) }),
	}

	t.Run("simple case", func(t *testing.T) {
		t.Parallel()
		in := make(Bi)
		data := []int{1, 2, 3, 4, 5}

		go func() {
			for _, v := range data {
				in <- v
			}
			close(in)
		}()

		result := make([]string, 0, 10)
		start := time.Now()
		for s := range ExecutePipeline(in, nil, stages...) {
			//switch v := s.(type) {
			//case string:
			//	fmt.Printf("Format %T\n", v)
			//	result = append(result, v)
			//case int:
			//	fmt.Printf("Format %T\n", v)
			//	result = append(result, fmt.Sprintf("%d", v))
			//default:
			//	t.Fatalf("Unexpected type %T", v)
			//}
			result = append(result, s.(string))
		}
		fmt.Println(result)
		elapsed := time.Since(start)

		require.Equal(t, []string{"1", "2", "3", "4", "5"}, result)
		require.Less(t,
			int64(elapsed),
			// ~0.8s for processing 5 values in 4 stages (100ms every) concurrently
			int64(sleepPerStage)*int64(len(stages)+len(data)-1)+int64(fault))
	})

}
