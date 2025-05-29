// Package customerrors provides custom errors.
package customerrors

import (
	"net/http"

	"emperror.dev/errors"
)

// NewConflictError creates a new conflict error.
func NewConflictError(message string) ConflictError {
	// `NewPlain` doesn't add stack-trace at all
	conflictErrMessage := errors.NewPlain("conflict error")
	// `WrapIf` add stack-trace if not added before
	stackErr := errors.WrapIf(conflictErrMessage, message)

	conflictError := &conflictError{
		CustomError: NewCustomError(stackErr, http.StatusConflict, message),
	}

	return conflictError
}

// NewConflictErrorWrap creates a new conflict error wrap.
func NewConflictErrorWrap(err error, message string) ConflictError {
	if err == nil {
		return NewConflictError(message)
	}

	// `WithMessage` doesn't add stack-trace at all
	conflictErrMessage := errors.WithMessage(err, "conflict error")
	// `WrapIf` add stack-trace if not added before
	stackErr := errors.WrapIf(conflictErrMessage, message)

	conflictError := &conflictError{
		CustomError: NewCustomError(stackErr, http.StatusConflict, message),
	}

	return conflictError
}

// conflictError is a struct that represents a conflict error.
type conflictError struct {
	CustomError
}

// ConflictError is a contract that represents a conflict error.
type ConflictError interface {
	CustomError
	isConflictError()
}

func (c *conflictError) isConflictError() {
}

// IsConflictError checks if the error is a conflict error.
func IsConflictError(err error) bool {
	// https://github.com/golang/go/blob/master/src/net/error_windows.go#L10C2-L12C3
	// this doesn't work for a nested notfound error, and we should use errors.As for traversing errors in all levels
	var conflictError ConflictError
	if errors.As(err, &conflictError) {
		return true
	}

	// us, ok := errors.Cause(err).(ConflictError)
	if errors.As(err, &conflictError) {
		return true
	}

	return false
}
