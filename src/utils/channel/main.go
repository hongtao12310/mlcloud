package main

import (
    "github.com/deepinsight/mlcloud/src/utils/log"
    "time"
    "fmt"
)

// implement the concurrency http request

var (
    MaxWorkers = 5
)


// Job represents the job to be run
type Job struct {
    Id      int     `json:"id"`
    Name    string  `json:"name"`
}

// A buffered channel that we can send work requests on.
var JobQueue chan Job = make(chan Job)

var ResultMap map[string]chan string

func init() {
    ResultMap = make(map[string]chan string)
}

// Worker represents the worker that executes the job
type Worker struct {
    WorkerPool  chan chan Job
    JobChannel  chan Job
    quit    	chan bool
}

func NewWorker(workerPool chan chan Job) Worker {
    return Worker {
        WorkerPool: workerPool,
        JobChannel: make(chan Job),
        quit:       make(chan bool)}
}

// Start method starts the run loop for the worker, listening for a quit channel in
// case we need to stop it
func (w Worker) Start() {
    go func() {
        for {
            // register the current worker into the worker queue.
            w.WorkerPool <- w.JobChannel
            log.Debug("assigned job channel to worker pool")

            select {
            case job := <-w.JobChannel:
                log.Debugf("receive job from channel: %#v\n sleep 5 seconds", job)
                time.Sleep(5 * time.Second)
                ResultMap[job.Name] <- fmt.Sprintf("finish job: %s", job.Name)

            case <-w.quit:
            // we have received a signal to stop
                log.Debugf("receive quit signal from worker")
                return
            }
        }
    }()
}

// Stop signals the worker to stop listening for work requests.
func (w Worker) Stop() {
    go func() {
        w.quit <- true
    }()
}

func feedJob(id int) {
    name := fmt.Sprintf("job-%d", id)

    job := Job {
        Id: id,
        Name: name,
    }

    // push the job name to result map
    ResultMap[job.Name] = make(chan string)

    JobQueue <- job

}


type Dispatcher struct {
    // A pool of workers channels that are registered with the dispatcher
    WorkerPool chan chan Job
    Workers []Worker
}

func NewDispatcher(maxWorkers int) *Dispatcher {
    pool := make(chan chan Job, maxWorkers)
    return &Dispatcher{
        WorkerPool: pool,
    }
}

func (d *Dispatcher) Run() {
    // starting n number of workers
    for i := 0; i < MaxWorkers; i++ {
        worker := NewWorker(d.WorkerPool)
        d.Workers = append(d.Workers, worker)
        worker.Start()
    }

    go d.dispatch()
}

func (d *Dispatcher) Stop() {
    for _, worker := range d.Workers {
        worker.Stop()
    }
}


func (d *Dispatcher) dispatch() {
    for {
        select {
        case job := <-JobQueue:
        // a job request has been received
            go func(job Job) {

                log.Debug("ready to get a job channel from worker poll")

                // try to obtain a worker job channel that is available.
                // this will block until a worker is idle
                jobChannel := <- d.WorkerPool

                log.Debug("got a job channel from worker poll")
                // dispatch the job to the worker job channel
                jobChannel <- job
            } (job)
        }
    }
}



func main() {
    dispatcher := NewDispatcher(MaxWorkers)
    dispatcher.Run()

    //for id := 0; id < MaxWorkers+1; id++ {
    //    feedJob(id)
    //}

    feedJob(1)

    for name, c := range ResultMap {
        select {
        case result := <- c:
            log.Debugf("jobName: %s, result: %s", name, result)
        }
    }


    // stop worker

    //wg := &sync.WaitGroup{}
    //
    //wg.Add(2)
    //queue := make(chan string)
    //
    //quit := make(chan bool)
    //// producer
    //go func(wgc *sync.WaitGroup) {
    //    queue <- "hello world"
    //
    //    wgc.Done()
    //}(wg)
    //
    //// consumer
    //go func(wgc *sync.WaitGroup) {
    //    for  {
    //        select {
    //        case msg := <- queue:
    //            log.Debug("received message: ", msg)
    //        case <- quit:
    //            log.Debug("consumer quit....")
    //            break
    //
    //        }
    //    }
    //
    //    wgc.Done()
    //
    //}(wg)
    //
    //
    //time.Sleep(1 * time.Second)
    //
    //quit <- true
    //
    //wg.Wait()
    //time.Sleep(5 * time.Second)
    log.Debug("Main Thread End")
}


