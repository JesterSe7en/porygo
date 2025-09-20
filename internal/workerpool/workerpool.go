// Copyright (c) 2025 Alexander Chan
// SPDX-License-Identifier: MIT

package workerpool

import "sync"

type WorkerPool struct {
	jobs    chan Job
	results chan Result
	wg      sync.WaitGroup
}

func New(workerCount int, bufferSize int) *WorkerPool {
	return &WorkerPool{
		jobs:    make(chan Job, bufferSize),
		results: make(chan Result, bufferSize),
		// WaitGroup not needed as it is initialized to its zero value - aka ready to use
	}
}

func (wp *WorkerPool) Submit(job Job) {
	wp.jobs <- job
}

func (wp *WorkerPool) Run(workerCount int) {
	for range workerCount {
		// WaitGroup.Go only available in Go >v1.25
		wp.wg.Go(func() {
			for job := range wp.jobs {
				wp.results <- job()
			}
		})
	}
}

func (wp *WorkerPool) Close() {
	close(wp.jobs)
	wp.wg.Wait()
	close(wp.results)
}

func (wp *WorkerPool) Results() <-chan Result {
	return wp.results
}
