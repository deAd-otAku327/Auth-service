package cryptor

type workerPool struct {
	tasks chan func()
}

func NewWorkerPool(maxWorkers int) *workerPool {
	pool := &workerPool{
		tasks: make(chan func(), maxWorkers),
	}

	for i := 0; i < maxWorkers; i++ {
		go pool.worker()
	}
	return pool
}

func (p *workerPool) worker() {
	for task := range p.tasks {
		task()
	}
}

func (p *workerPool) Add(task func()) {
	p.tasks <- task
}
