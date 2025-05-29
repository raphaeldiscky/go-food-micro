// Package customerrors provides custom errors.
package customerrors

import (
	"net/http"

	"emperror.dev/errors"
)

// NewNotFoundError creates a new not found error.
func NewNotFoundError(message string) NotFoundError {
	// `NewPlain` doesn't add stack-trace at all
	notFoundErrMessage := errors.NewPlain("not found error")
	// `WrapIf` add stack-trace if not added before
	stackErr := errors.WrapIf(notFoundErrMessage, message)

	notFoundError := &notFoundError{
		CustomError: NewCustomError(stackErr, http.StatusBadRequest, message),
	}

	return notFoundError
}

// NewNotFoundErrorWrap creates a new not found error.
func NewNotFoundErrorWrap(err error, message string) NotFoundError {
	if err == nil {
		return NewNotFoundError(message)
	}

	// `WithMessage` doesn't add stack-trace at all
	notFoundErrMessage := errors.WithMessage(err, "not found error")
	// `WrapIf` add stack-trace if not added before
	stackErr := errors.WrapIf(notFoundErrMessage, message)

	notFoundError := &notFoundError{
		CustomError: NewCustomError(stackErr, http.StatusNotFound, message),
	}

	return notFoundError
}

// notFoundError is a not found error.
type notFoundError struct {
	CustomError
}

// NotFoundError is a not found error.
type NotFoundError interface {
	CustomError
	isNotFoundError()
}

// isNotFoundError checks if the error is a not found error.
func (n *notFoundError) isNotFoundError() {
}

// IsNotFoundError checks if the error is a not found error.
func IsNotFoundError(err error) bool {
	// https://github.com/golang/go/blob/master/src/net/error_windows.go#L10C2-L12C3
	// this doesn't work for a nested notfound error, and we should use errors.As for traversing errors in all levels
	var notFoundError NotFoundError
	if errors.As(err, &notFoundError) {
		return true
	}

	// us, ok := errors.Cause(err).(NotFoundError)
	if errors.As(err, &notFoundError) {
		return true
	}

	return false
}
