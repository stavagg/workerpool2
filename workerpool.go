package workerpool

import (
	"sync"
)

type WorkerPool struct {
	tasks   chan func()
	workers int
	wg      sync.WaitGroup
	once    sync.Once
	stopped bool
	mu      sync.Mutex
}

func NewWorkerPool(numberOfWorkers int) *WorkerPool {
	wp := &WorkerPool{
		tasks:   make(chan func(), numberOfWorkers),
		workers: numberOfWorkers,
	}

	for i := 0; i < numberOfWorkers; i++ {
		wp.wg.Add(1)
		go wp.worker()
	}

	return wp
}
