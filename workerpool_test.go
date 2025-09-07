package workerpool

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestSubmitExecutesTasks(t *testing.T) {
	wp := NewWorkerPool(3)
	var counter int64
	const n = 50

	for i := 0; i < n; i++ {
		wp.Submit(func() {
			atomic.AddInt64(&counter, 1)
		})
	}

	wp.Stop()
	if got := atomic.LoadInt64(&counter); got != n {
		t.Fatalf("expected %d tasks executed, got %d", n, got)
	}
}

func TestSubmitWait(t *testing.T) {
	wp := NewWorkerPool(2)
	var v int64

	wp.SubmitWait(func() {
		time.Sleep(50 * time.Millisecond)
		atomic.StoreInt64(&v, 42)
	})

	if v != 42 {
		t.Fatalf("SubmitWait did not wait task completion, v = %d", v)
	}
	wp.Stop()
}

func TestStopDoesNotAcceptNew(t *testing.T) {
	wp := NewWorkerPool(1)
	var counter int64

	wp.Submit(func() { atomic.AddInt64(&counter, 1) })
	wp.Stop()
	wp.Submit(func() { atomic.AddInt64(&counter, 1) })

	if got := atomic.LoadInt64(&counter); got != 1 {
		t.Fatalf("expected only 1 task executed, got %d", got)
	}
}

func TestStopWaitExecutesQueued(t *testing.T) {
	wp := NewWorkerPool(1)
	var counter int64

	wp.Submit(func() {
		time.Sleep(100 * time.Millisecond)
		atomic.AddInt64(&counter, 1)
	})
	wp.Submit(func() {
		atomic.AddInt64(&counter, 1)
	})

	start := time.Now()
	wp.StopWait()
	elapsed := time.Since(start)

	if got := atomic.LoadInt64(&counter); got != 2 {
		t.Fatalf("expected 2 tasks executed, got %d", got)
	}
	if elapsed < 100*time.Millisecond {
		t.Fatalf("StopWait returned too early: elapsed %v", elapsed)
	}
}

func TestConcurrentSubmitAndStop(t *testing.T) {
	wp := NewWorkerPool(5)
	var counter int64
	const n = 100

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		for i := 0; i < n; i++ {
			wp.Submit(func() {
				atomic.AddInt64(&counter, 1)
			})
		}
	}()

	go func() {
		defer wg.Done()
		time.Sleep(10 * time.Millisecond)
		wp.Stop()
	}()

	wg.Wait()
	time.Sleep(50 * time.Millisecond)

	got := atomic.LoadInt64(&counter)
	if got < 1 || got > n {
		t.Fatalf("counter out of range: %d", got)
	}
}
