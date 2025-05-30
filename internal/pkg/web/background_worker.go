// Package web provides a background worker.
package web

import (
	"context"
)

// Worker is a interface that represents a worker.
type Worker interface {
	Start(ctx context.Context) chan error
	Stop(ctx context.Context) error
}

type (
	// ExecutionFunc is a function that represents a execution func.
	ExecutionFunc func(ctx context.Context) error
	// StopFunc is a function that represents a stop func.
	StopFunc func(ctx context.Context) error
)

// BackgroundWorker is a struct that represents a background worker.
type BackgroundWorker struct {
	ctx           context.Context
	executionFunc ExecutionFunc
	stopFunc      StopFunc
	cancelFunc    context.CancelFunc
	errChan       chan error
}

// NewBackgroundWorker is a function that creates a new background worker.
func NewBackgroundWorker(executionFunc ExecutionFunc, stopFunc StopFunc) Worker {
	return &BackgroundWorker{
		executionFunc: executionFunc,
		stopFunc:      stopFunc,
		errChan:       make(chan error),
	}
}

// Start is a function that starts the background worker.
func (b BackgroundWorker) Start(ctx context.Context) chan error {
	b.ctx, b.cancelFunc = context.WithCancel(ctx)
	go func() {
		if b.executionFunc == nil {
			return
		}

		err := b.executionFunc(b.ctx)
		if err != nil {
			b.cancelFunc()
			b.errChan <- err
		}
	}()

	return b.errChan
}

// Stop is a function that stops the background worker.
func (b BackgroundWorker) Stop(ctx context.Context) error {
	if b.executionFunc == nil {
		return nil
	}
	if b.stopFunc != nil {
		return b.stopFunc(ctx)
	}
	if b.cancelFunc != nil {
		b.cancelFunc()
	}

	return nil
}
