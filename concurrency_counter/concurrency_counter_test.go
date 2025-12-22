package concurrency_counter

import (
	"fmt"
	"runtime"
	"sync"
	"sync/atomic"
	"testing"
)

// Benchmarks comparing common concurrency counter patterns.
// Variants:
// - sync.Mutex (write)
// - sync.RWMutex (write + read-heavy)
// - atomic.AddInt64
// - channels (various buffer sizes)
// - worker-pool (fixed workers consuming jobs)

func BenchmarkMutexParallel(b *testing.B) {
	var mu sync.Mutex
	var counter int64

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			mu.Lock()
			counter++
			mu.Unlock()
		}
	})
}

func BenchmarkRWMutexWriteParallel(b *testing.B) {
	var mu sync.RWMutex
	var counter int64

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			mu.Lock()
			counter++
			mu.Unlock()
		}
	})
}

func BenchmarkRWMutexReadParallel(b *testing.B) {
	var mu sync.RWMutex
	var counter int64

	// populate counter
	counter = 42

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			mu.RLock()
			_ = counter
			mu.RUnlock()
		}
	})
}

func BenchmarkAtomicParallel(b *testing.B) {
	var counter int64

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			atomic.AddInt64(&counter, 1)
		}
	})
}

func BenchmarkChannelBufferedSizes(b *testing.B) {
	sizes := []int{0, 1, 16, 1024}

	for _, sz := range sizes {
		b.Run(fmt.Sprintf("buf=%d", sz), func(b *testing.B) {
			ch := make(chan struct{}, sz)
			var wg sync.WaitGroup
			var counter int64

			wg.Add(1)
			go func() {
				for range ch {
					atomic.AddInt64(&counter, 1)
				}
				wg.Done()
			}()

			b.RunParallel(func(pb *testing.PB) {
				for pb.Next() {
					ch <- struct{}{}
				}
			})

			close(ch)
			wg.Wait()
		})
	}
}

func BenchmarkWorkerPool(b *testing.B) {
	workers := runtime.GOMAXPROCS(0)

	b.Run(fmt.Sprintf("workers=%d", workers), func(b *testing.B) {
		jobs := make(chan struct{}, 1024)
		var wg sync.WaitGroup
		var counter int64

		wg.Add(workers)
		for i := 0; i < workers; i++ {
			go func() {
				for range jobs {
					atomic.AddInt64(&counter, 1)
				}
				wg.Done()
			}()
		}

		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				jobs <- struct{}{}
			}
		})

		close(jobs)
		wg.Wait()
	})
}
