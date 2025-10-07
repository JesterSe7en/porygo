package workerpool

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"
)

func TestWorkerPool(t *testing.T) {
	t.Run("Test New", func(t *testing.T) {
		const numWorkers = 10
		wp := New(numWorkers, 10)
		if wp == nil {
			t.Fatal("expected worker pool to not be nil")
		}

		// verify the integrity of the worker TestWorkerPool
		if cap(wp.jobs) != numWorkers {
			t.Errorf("expected jobs channel capacity to be 10, but got %d", cap(wp.jobs))
		}
		if cap(wp.results) != 10 {
			t.Errorf("expected workers slice length to be 10, but got %d", cap(wp.results))
		}
	})

	t.Run("Test Submit and Run", func(t *testing.T) {
		const numWorkers = 2
		wp := New(numWorkers, 10)
		ctx := context.Background()

		var wg sync.WaitGroup

		wg.Add(1)
		go func() {
			defer wg.Done()
			wp.Run(ctx, numWorkers)
		}()

		time.Sleep(10 * time.Millisecond)

		job := func() Result {
			return Result{Value: "test-value"}
		}

		if err := wp.Submit(ctx, job); err != nil {
			t.Fatalf("Failed to submit job: %v", err)
		}

		result := <-wp.Results()
		if result.Value != "test-value" {
			t.Errorf("Expected result value %q, but got %q", "test-value", result.Value)
		}

		wp.Close()
		wg.Wait()
	})

	t.Run("Test context cancellation", func(t *testing.T) {
		const numWorkers = 2
		wp := New(numWorkers, 10)

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		done := make(chan struct{})

		go func() {
			wp.Run(ctx, numWorkers)
			close(done)
		}()

		jobStarted := make(chan struct{})
		blockingJob := func() Result {
			close(jobStarted)
			<-ctx.Done()
			return Result{Value: "done"}
		}

		if err := wp.Submit(ctx, blockingJob); err != nil {
			t.Fatalf("failed to submit initial job: %v", err)
		}

		<-jobStarted

		cancel()

		select {
		case <-done:
		case <-time.After(time.Second):
			t.Fatal("worker pool did not exit after context cancellation")
		}

		err := wp.Submit(ctx, func() Result { return Result{Value: "test"} })
		if err == nil {
			t.Fatal("expected error when submitting job to canceled context")
		}

		wp.Close()
	})

	t.Run("Test job error", func(t *testing.T) {
		const numWorkers = 2
		wp := New(numWorkers, 10)
		ctx := context.Background()

		var wg sync.WaitGroup

		wg.Add(1)
		go func() {
			defer wg.Done()
			wp.Run(ctx, 1)
		}()

		job := func() Result {
			return Result{Err: errors.New("test-error")}
		}

		if err := wp.Submit(ctx, job); err != nil {
			t.Fatalf("Failed to submit job: %v", err)
		}

		result := <-wp.Results()
		if result.Err == nil {
			t.Fatal("Expected error in result")
		}

		if result.Err.Error() != "test-error" {
			t.Errorf("Expected error message %s, but got %s", "test-error", result.Err.Error())
		}

		wp.Close()
		wg.Wait()
	})
}
