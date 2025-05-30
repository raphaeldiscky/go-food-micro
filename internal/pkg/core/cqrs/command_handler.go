// Package cqrs provides a module for the cqrs.
package cqrs

import (
	"context"

	mediatr "github.com/mehdihadeli/go-mediatr"
)

// CommandHandler is a command handler.
type CommandHandler[TCommand Command, TResponse any] interface {
	Handle(ctx context.Context, command TCommand) (TResponse, error)
}

// CommandHandlerVoid is a command handler void.
type CommandHandlerVoid[TCommand Command] interface {
	Handle(ctx context.Context, command TCommand) (*mediatr.Unit, error)
}
