// Package customerrors provides custom errors.
package customerrors

import (
	"net/http"

	"emperror.dev/errors"
)

// NewForbiddenError creates a new forbidden error.
func NewForbiddenError(message string) ForbiddenError {
	// `NewPlain` doesn't add stack-trace at all
	forbiddenErrMessage := errors.NewPlain("forbidden error")
	// `WrapIf` add stack-trace if not added before
	stackErr := errors.WrapIf(forbiddenErrMessage, message)

	forbiddenError := &forbiddenError{
		CustomError: NewCustomError(stackErr, http.StatusForbidden, message),
	}

	return forbiddenError
}

// NewForbiddenErrorWrap creates a new forbidden error.
func NewForbiddenErrorWrap(err error, message string) ForbiddenError {
	if err == nil {
		return NewForbiddenError(message)
	}

	// `WithMessage` doesn't add stack-trace at all
	forbiddenErrMessage := errors.WithMessage(err, "forbidden error")
	// `WrapIf` add stack-trace if not added before
	stackErr := errors.WrapIf(forbiddenErrMessage, message)

	forbiddenError := &forbiddenError{
		CustomError: NewCustomError(stackErr, http.StatusForbidden, message),
	}

	return forbiddenError
}

// forbiddenError is a forbidden error.
type forbiddenError struct {
	CustomError
}

// ForbiddenError is a forbidden error.
type ForbiddenError interface {
	CustomError
	isForbiddenError()
}

// isForbiddenError checks if the error is a forbidden error.
func (f *forbiddenError) isForbiddenError() {
}

// IsForbiddenError checks if the error is a forbidden error.
func IsForbiddenError(err error) bool {
	// https://github.com/golang/go/blob/master/src/net/error_windows.go#L10C2-L12C3
	// this doesn't work for a nested forbidden error, and we should use errors.As for traversing errors in all levels
	var forbiddenError ForbiddenError
	if errors.As(err, &forbiddenError) {
		return true
	}

	// us, ok := errors.Cause(err).(ForbiddenError)
	if errors.As(err, &forbiddenError) {
		return true
	}

	return false
}
