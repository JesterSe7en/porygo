// Copyright (c) 2025 Alexander Chan
// SPDX-License-Identifier: MIT

package workerpool

import (
	"context"
	"sync"
)

type WorkerPool struct {
	jobs    chan Job
	results chan Result
	wg      sync.WaitGroup
}

func New(workerCount int, bufferSize int) *WorkerPool {
	return &WorkerPool{
		jobs:    make(chan Job, bufferSize),
		results: make(chan Result, bufferSize),
	}
}

// Submit adds a job to the pool, but will not block indefinitely.
// It returns an error if the context is canceled before the job can be submitted.
func (wp *WorkerPool) Submit(ctx context.Context, job Job) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case wp.jobs <- job:
		return nil
	}
}

// Run starts the workers. It now accepts a context to enable graceful shutdown.
func (wp *WorkerPool) Run(ctx context.Context, workerCount int) {
	wp.wg.Add(workerCount)
	for range workerCount {
		go func() {
			defer wp.wg.Done()
			for {
				// essentially checking if either a cancel was requested
				// or the job was closed either due to finishing or being canceled
				select {
				case <-ctx.Done():
					return
				case job, ok := <-wp.jobs:
					if !ok {
						return
					}
					wp.results <- job()
				}
			}
		}()
	}
}

// Close waits for all jobs to be processed and then closes the results channel.
// It should be called after all jobs have been submitted.
func (wp *WorkerPool) Close() {
	close(wp.jobs)
	wp.wg.Wait()
	close(wp.results)
}

func (wp *WorkerPool) Results() <-chan Result {
	return wp.results
}
