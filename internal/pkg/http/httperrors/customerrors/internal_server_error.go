// Package customerrors provides custom errors.
package customerrors

import (
	"net/http"

	"emperror.dev/errors"
)

// NewInternalServerError creates a new internal server error.
func NewInternalServerError(message string) InternalServerError {
	// `NewPlain` doesn't add stack-trace at all
	internalErrMessage := errors.NewPlain("internal server error")
	// `WrapIf` add stack-trace if not added before
	stackErr := errors.WrapIf(internalErrMessage, message)

	internalServerError := &internalServerError{
		CustomError: NewCustomError(stackErr, http.StatusInternalServerError, message),
	}

	return internalServerError
}

// NewInternalServerErrorWrap creates a new internal server error.
func NewInternalServerErrorWrap(err error, message string) InternalServerError {
	if err == nil {
		return NewInternalServerError(message)
	}

	// `WithMessage` doesn't add stack-trace at all
	internalErrMessage := errors.WithMessage(err, "internal server error")
	// `WrapIf` add stack-trace if not added before
	stackErr := errors.WrapIf(internalErrMessage, message)

	internalServerError := &internalServerError{
		CustomError: NewCustomError(stackErr, http.StatusInternalServerError, message),
	}

	return internalServerError
}

// internalServerError is a internal server error.
type internalServerError struct {
	CustomError
}

// InternalServerError is a internal server error.
type InternalServerError interface {
	CustomError
	isInternalServerError()
}

// isInternalServerError checks if the error is a internal server error.
func (i *internalServerError) isInternalServerError() {
}

// IsInternalServerError checks if the error is a internal server error.
func IsInternalServerError(err error) bool {
	var internalServerErr InternalServerError

	// https://github.com/golang/go/blob/master/src/net/error_windows.go#L10C2-L12C3
	// this doesn't work for a nested internal server error, and we should use errors.As for traversing errors in all levels
	var internalServerError InternalServerError
	if errors.As(err, &internalServerError) {
		return true
	}

	// us, ok := errors.Cause(err).(InternalServerError)
	if errors.As(err, &internalServerErr) {
		return true
	}

	return false
}
