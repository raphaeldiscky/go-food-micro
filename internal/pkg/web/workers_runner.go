// Package web provides a workers runner.
package web

import "context"

// WorkersRunner is a struct that represents a workers runner.
type WorkersRunner struct {
	workers []Worker
	errChan chan error
}

// NewWorkersRunner is a function that creates a new workers runner.
func NewWorkersRunner(workers []Worker) *WorkersRunner {
	return &WorkersRunner{workers: workers, errChan: make(chan error)}
}

// Start is a function that starts the workers runner.
func (r *WorkersRunner) Start(ctx context.Context) chan error {
	if len(r.workers) == 0 {
		return nil
	}

	for _, w := range r.workers {
		err := w.Start(ctx)
		go func() {
			for {
				select {
				case e := <-err:
					r.errChan <- e

					return
				case <-ctx.Done():
					stopErr := r.Stop(ctx)
					if stopErr != nil {
						r.errChan <- stopErr

						return
					}

					return
				}
			}
		}()
	}

	return r.errChan
}

// Stop is a function that stops the workers runner.
func (r *WorkersRunner) Stop(ctx context.Context) error {
	if len(r.workers) == 0 {
		return nil
	}

	for _, w := range r.workers {
		err := w.Stop(ctx)
		if err != nil {
			return err
		}
	}

	return nil
}
