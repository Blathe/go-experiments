package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type job struct {
	Mux      sync.RWMutex
	Index    int
	Duration int //in seconds
	Complete bool
}

type worker struct {
	Index      int
	CurrentJob int
	HasJob     bool
	WorkerPool *workerpool
}

type workerpool struct {
	Jobs    []*job
	Workers []*worker
}

type jobpool struct {
	Mux  sync.RWMutex
	Jobs []*job
}

func (wp *workerpool) HaveWorkersFindJobs(jp *jobpool) {
	for _, worker := range wp.Workers {
		if worker.HasJob == false {
			worker.CheckForMoreJobs(jp)
		}
	}
}

func (w *worker) StartJob(j *job, pool *jobpool) {
	j.Mux.Lock()
	j.Complete = true
	j.Mux.Unlock()

	w.HasJob = true
	w.CurrentJob = j.Index

	fmt.Printf("Worker %v starting job %v\n", w.Index, j.Index)

	sleepTime := time.Duration(j.Duration) * time.Second
	time.Sleep(sleepTime)
	pool.Mux.Lock()
	pool.Jobs = append(pool.Jobs[:j.Index], pool.Jobs[j.Index+1:]...)
	pool.Mux.Unlock()
	fmt.Printf("Worker %v finished job %v\n", w.Index, j.Index)
	fmt.Printf("Removing job %v from pool...%v jobs left...\n", j.Index, len(pool.Jobs))
	w.HasJob = false

	w.CheckForMoreJobs(pool)
}

func (w *worker) CheckForMoreJobs(pool *jobpool) {
	for _, j := range pool.Jobs {
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
	var jobs []*job
	min := 1
	max := 2

	workerPool := workerpool{
		Workers: nil,
		Jobs:    jobs,
	}

	jobPool := jobpool{
		Jobs: jobs,
	}

	for i := 0; i < 5000; i++ {
		newJob := job{
			Index:    i,
			Duration: rand.Intn(max-min) + min,
			Complete: false,
		}
		jobPool.Jobs = append(jobPool.Jobs, &newJob)
	}

	for i := 0; i < 7; i++ {
		newWorker := worker{
			Index:      i,
			CurrentJob: i,
			HasJob:     false,
			WorkerPool: &workerPool,
		}
		workerPool.Workers = append(workerPool.Workers, &newWorker)

		fmt.Println("Spawning worker...")
		go newWorker.StartJob(jobPool.Jobs[i], &jobPool)
	}

	time.Sleep(3 * time.Second)
	AddMoreJobs(jobPool.Jobs)
	time.Sleep(8 * time.Second)
	AddMoreJobs(jobPool.Jobs)
	fmt.Scanln()
}

func AddMoreJobs(pool []*job) {

	fmt.Println("More jobs being added to the pool... -------------------")

	min := 1
	max := 2

	for i := 0; i < 20; i++ {
		newJob := job{
			Index:    i,
			Duration: rand.Intn(max-min) + min,
			Complete: false,
		}
		pool = append(pool, &newJob)
	}
}
