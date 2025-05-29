// Package contracts provides the contracts for the http errors.
package contracts

import (
	"fmt"

	"emperror.dev/errors"
)

// Causer is a contract that represents a causer.
type Causer interface {
	Cause() error
}

// StackTracer is a contract that represents a stack tracer.
type StackTracer interface {
	StackTrace() errors.StackTrace
}

// Wrapper is a contract that represents a wrapper.
type Wrapper interface {
	Unwrap() error
}

// Formatter is a contract that represents a formatter.
type Formatter interface {
	Format(f fmt.State, verb rune)
}

// BaseError is a contract that represents a base error.
type BaseError interface {
	error
	Wrapper
	Causer
	Formatter
}
