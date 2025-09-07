package workerpool

func (wp *WorkerPool) Submit(task func()) {
	wp.mu.Lock()
	defer wp.mu.Unlock()

	if wp.stopped {
		return
	}
	wp.tasks <- task
}

func (wp *WorkerPool) SubmitWait(task func()) {
	done := make(chan struct{})
	wp.Submit(func() {
		task()
		close(done)
	})
	<-done
}
