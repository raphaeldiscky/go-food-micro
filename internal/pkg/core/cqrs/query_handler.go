// Package cqrs provides a module for the cqrs.
package cqrs

import (
	"context"
)

// QueryHandler is a query handler.
type QueryHandler[TQuery Query, TResponse any] interface {
	Handle(ctx context.Context, query TQuery) (TResponse, error)
}
