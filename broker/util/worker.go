package util

type Handle func(Action)

// Job represents the job to be run
type Job struct {
	Payload Action
    Do Handle
}

// Worker represents the worker that executes the job
type Worker struct {
	WorkerPool chan chan Job
	JobChannel chan Job
	quit chan bool
}

func NewWorker(workerPool chan chan Job) Worker {
	return Worker{
		WorkerPool: workerPool,
		JobChannel: make(chan Job),
		quit: make(chan bool),
	}
}

// Start method starts the run loop for the worker, listening for a quit channel in case we need to stop it
func (w Worker) Start() {
	go func() {
		for {
			// register the current worker into the worker queue
			w.WorkerPool <- w.JobChannel
			select {
			case job := <-w.JobChannel:
				job.Do(job.Payload)

			case <-w.quit:
			// we have received a signal to stop
				return
			}
		}
	}()
}

// Stop method signals the worker to stop listening for work requests
func (w Worker) Stop() {
	go func() {
		w.quit <- true
	}()
}
