package main

import (
    "testing"
    "sync"
    "time"
    "github.com/deepinsight/mlcloud/src/utils/log"
    "fmt"
)

func TestWaitGroup(t *testing.T)  {
    var vg sync.WaitGroup

    vg.Add(2)
    // producer
    go func() {
        t.Log("producer sleep 1 second")
        time.Sleep(1 * time.Second)

        vg.Done()
    }()


    // consumer
    go func() {
        t.Log("producer sleep 2 second")
        time.Sleep(2 * time.Second)
        vg.Done()

    }()

    vg.Wait()
}

func TestDeadLock(t *testing.T)  {

    queue := make(chan string)

    quit := make(chan bool)
    // producer
    go func() {
        queue <- "hello world"
    }()

    // consumer
    go func() {
        for  {
            select {
            case msg := <- queue:
                t.Log("received message: ", msg)
            case <- quit:
                log.Debug("consumer quit....")

            }
        }

    }()


    time.Sleep(1 * time.Second)

    quit <- true
}

type Worker struct {
    Name string
    Queue chan string
    Quit chan bool
    Avaliable bool
}

type WorkerPool struct {
    workers []Worker
}

func (self *WorkerPool) GetAvaWorker () Worker  {
    for _, worker := range self.workers {
        if worker.Avaliable {
            return worker
        }
    }
}

func (self *Worker) start ()  {
    go func() {
        self.Avaliable = false

        for  {
            select {
            case msg := <- self.Queue:
                log.Debugf("receive message: %s from queue ", msg)
                self.Avaliable = true
            case <- self.Quit:
                log.Debug("work End....")
                break
            }
        }

    }()
}

func (self *Worker) stop ()  {
    self.Quit <- true
}


const (
    MaxWorkers = 5

)


func TestWoker(t *testing.T)  {
    workerPool := WorkerPool{

    }

    for id := 0; id < 5; id ++ {
        name := fmt.Sprintf("worker:%d", id)
        worker := Worker{
            Name: name,
            Queue: make(chan string),
            Quit: make(chan bool),
            Avaliable: false,
        }
        workerPool.workers = append(workerPool.workers, worker)
    }


}