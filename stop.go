package workerpool

func (wp *WorkerPool) Stop() {
	wp.once.Do(func() {
		wp.mu.Lock()
		wp.stopped = true
		wp.mu.Unlock()

		close(wp.tasks)
		wp.wg.Wait()
	})
}

func (wp *WorkerPool) StopWait() {
	wp.Stop()
}
