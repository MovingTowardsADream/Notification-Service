package outbox

type WorkerRun interface {
	Run() error
	Stop()
}

type Worker struct {
	workers map[string]WorkerRun
}

func NewWorker(workers map[string]WorkerRun) *Worker {
	return &Worker{
		workers: workers,
	}
}

func (w *Worker) WorkerRun() error {
	for _, worker := range w.workers {
		go func(w WorkerRun) {
			_ = w.Run()
		}(worker)
	}
	return nil
}

func (w *Worker) Shutdown() error {
	for _, worker := range w.workers {
		worker.Stop()
	}
	return nil
}
