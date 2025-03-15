package outbox

type WorkerRun interface {
	Run() error
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
		go func() {
			_ = worker.Run()
		}()
	}
	return nil
}
