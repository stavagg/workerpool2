package workerpool

func (wp *WorkerPool) worker() {
	defer wp.wg.Done()

	for task := range wp.tasks {
		task()
	}
}
