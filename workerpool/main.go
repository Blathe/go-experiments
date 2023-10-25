package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type job struct {
	Mux             sync.RWMutex
	Index           int
	StartMessage    string
	Duration        int //in seconds
	CompleteMessage string
	Complete        bool
}

type worker struct {
	Index      int
	CurrentJob int
	HasJob     bool
	WorkerPool *workerpool
}

type workerpool struct {
	Jobs             []*job
	Workers          []*worker
	AvailableWorkers []*worker
}

func (wp *workerpool) HaveWorkersFindJobs(jp *[]*job) {
	for _, worker := range wp.AvailableWorkers {
		if worker.HasJob == false {
			worker.CheckForMoreJobs(jp)
		}
	}
}

func (w *worker) StartJob(j *job, pool *[]*job) {
	j.Mux.Lock()
	j.Complete = true
	j.Mux.Unlock()

	w.HasJob = true
	w.CurrentJob = j.Index

	fmt.Printf("Worker %v starting job %v\n", w.Index, j.Index)

	sleepTime := time.Duration(j.Duration) * time.Second
	time.Sleep(sleepTime)
	*pool = append((*pool)[:j.Index], (*pool)[j.Index:]...)
	fmt.Printf("Worker %v finished job %v\n", w.Index, j.Index)
	fmt.Printf("Removing job %v from pool...\n", j.Index)
	w.HasJob = false

	w.CheckForMoreJobs(pool)
}

func (w *worker) CheckForMoreJobs(pool *[]*job) {
	for _, j := range *pool {
		if j.Complete == false {
			w.StartJob(j, pool)
			return
		}
	}

	fmt.Printf("Worker %v can't find a new job...\n", w.Index)
	time.Sleep(15 * time.Second)
	w.CheckForMoreJobs(pool)
}

func main() {
	var jobPool []*job
	min := 1
	max := 2

	workerPool := workerpool{
		Workers:          nil,
		AvailableWorkers: nil,
		Jobs:             jobPool,
	}

	for i := 0; i < 5000; i++ {
		newJob := job{
			Index:           i,
			StartMessage:    fmt.Sprintf("Job %v started...", i),
			Duration:        rand.Intn(max-min) + min,
			CompleteMessage: fmt.Sprintf("Job %v complete!", i),
			Complete:        false,
		}
		jobPool = append(jobPool, &newJob)
	}

	for i := 0; i < 7; i++ {
		newWorker := worker{
			Index:      i,
			CurrentJob: i,
			HasJob:     false,
			WorkerPool: &workerPool,
		}
		workerPool.Workers = append(workerPool.Workers, &newWorker)
		workerPool.AvailableWorkers = append(workerPool.AvailableWorkers, &newWorker)

		fmt.Println("Spawning worker...")
		go newWorker.StartJob(jobPool[i], &jobPool)
	}

	time.Sleep(3 * time.Second)
	AddMoreJobs(&jobPool)
	time.Sleep(8 * time.Second)
	AddMoreJobs(&jobPool)
	fmt.Scanln()
}

func AddMoreJobs(pool *[]*job) {

	fmt.Println("More jobs being added to the pool... -------------------")

	min := 1
	max := 2

	for i := 0; i < 20; i++ {
		newJob := job{
			Index:           i,
			StartMessage:    fmt.Sprintf("Job %v started...", i),
			Duration:        rand.Intn(max-min) + min,
			CompleteMessage: fmt.Sprintf("Job %v complete!", i),
			Complete:        false,
		}
		*pool = append(*pool, &newJob)
	}
}
