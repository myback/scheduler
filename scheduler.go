package scheduler

import (
	"context"
	"os"
	"sync"
	"time"

	"github.com/spaolacci/murmur3"
)

var wg = sync.WaitGroup{}

type Scheduler struct {
	mu   sync.Mutex
	jobs map[uint32]*job
}

func (sch *Scheduler) add(ctx context.Context, key uint32, j *job) error {
	sch.mu.Lock()
	defer sch.mu.Unlock()

	if _, ok := sch.jobs[key]; ok {
		return ErrorJobAlreadyExists
	}

	j.start(ctx)
	sch.jobs[key] = j

	return nil
}

func (sch *Scheduler) cancel(key uint32) error {
	sch.mu.Lock()
	defer sch.mu.Unlock()

	if _, ok := sch.jobs[key]; !ok {
		return ErrorJobNotFound
	}

	sch.jobs[key].cancel()
	delete(sch.jobs, key)

	return nil
}

func (sch *Scheduler) update(key uint32, intervalSec time.Duration) error {
	sch.mu.Lock()
	defer sch.mu.Unlock()

	if _, ok := sch.jobs[key]; !ok {
		return ErrorJobNotFound
	}

	sch.jobs[key].updateInterval(intervalSec)

	return nil
}

func (sch *Scheduler) AddJob(ctx context.Context, namespacedName, jobName string, intervalSec int, j func()) error {
	return sch.add(ctx, murmur3.Sum32([]byte(namespacedName+"/"+jobName)), &job{
		namespacedName: namespacedName,
		ticker:         time.NewTicker(time.Duration(intervalSec) * time.Second),
		jobFunc:        j,
	})
}

func (sch *Scheduler) CancelJob(namespacedName string) error {
	if namespacedName == "" {
		return ErrorJobNotFound
	}

	id := murmur3.Sum32([]byte(namespacedName))

	var err error
	for _id, j := range sch.jobs {
		if j.namespacedName == namespacedName || _id == id {
			if e := sch.cancel(_id); e != nil {
				err = e
			}
		}
	}

	return err
}

func (sch *Scheduler) GetNumberJobs() int {
	sch.mu.Lock()
	defer sch.mu.Unlock()

	return len(sch.jobs)
}

func (sch *Scheduler) GracefulShutdown() {
	sch.mu.Lock()
	defer sch.mu.Unlock()

	for id, j := range sch.jobs {
		j.cancel()
		delete(sch.jobs, id)
	}
}

func (sch *Scheduler) Wait() {
	wg.Wait()
}

func (sch *Scheduler) WaitWithTimeout(timeout int) {
	done := make(chan bool)

	go func() {
		t := time.NewTimer(time.Duration(timeout) * time.Second)
		defer t.Stop()

		for {
			select {
			case <-t.C:
				os.Exit(201)
			case <-done:
				return
			}
		}
	}()

	wg.Wait()
	close(done)
}

func (sch *Scheduler) UpdateJobInterval(namespacedName string, intervalSec int) error {
	return sch.update(murmur3.Sum32([]byte(namespacedName)), time.Duration(intervalSec)*time.Second)
}

func New() *Scheduler {
	return &Scheduler{
		mu:   sync.Mutex{},
		jobs: map[uint32]*job{},
	}
}
