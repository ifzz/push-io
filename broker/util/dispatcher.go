package util

type Dispatcher struct {
	// A pool of workers channels that are registered with the dispatcher
	workerPool chan chan Job
	jobQueue chan Job
}

func NewDispatcher(jobQueue chan Job) *Dispatcher {
	workerPool := make(chan chan Job, config.MaxWorkers)

	return &Dispatcher{
		workerPool: workerPool,
		jobQueue: jobQueue}
}

func (d *Dispatcher) Run() {
	// starting n number of workers
	for i := 0; i < config.MaxWorkers; i++ {
		worker := NewWorker(d.workerPool)
		worker.Start()
	}

	go d.dispatch()
}

func (d *Dispatcher) dispatch() {
	for {
		select {
		// a job request has been received
		case job := <-d.jobQueue:
		// try to obtain a worker job channel that is available
		// this will block until a worker is idle
			go func(job Job) {
				jobChannel := <-d.workerPool

				// dispatch the job to the worker job channel
				jobChannel <- job
			}(job)
		}
	}
}
