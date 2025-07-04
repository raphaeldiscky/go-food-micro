// Package cqrs provides a module for the cqrs.
package cqrs

import mediatr "github.com/mehdihadeli/go-mediatr"

// HandlerRegisterer for registering `RequestHandler` to mediatr registry, if handler implements this interface it will be registered automatically.
// HandlerRegisterer for registering `RequestHandler` to mediatr registry, if handler implements this interface it will be registered automatically.
type HandlerRegisterer interface {
	RegisterHandler() error
}

// RequestHandlerWithRegisterer for registering `RequestHandler` to mediatr registry and handling `RequestHandler`.
type RequestHandlerWithRegisterer[TRequest any, TResponse any] interface {
	HandlerRegisterer
	mediatr.RequestHandler[TRequest, TResponse]
}
