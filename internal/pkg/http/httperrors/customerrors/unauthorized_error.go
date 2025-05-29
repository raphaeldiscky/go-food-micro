// Package customerrors provides custom errors.
package customerrors

import (
	"net/http"

	"emperror.dev/errors"
)

// NewUnAuthorizedError creates a new unauthorized error.
func NewUnAuthorizedError(message string) UnauthorizedError {
	// `NewPlain` doesn't add stack-trace at all
	unAuthorizedErrMessage := errors.NewPlain("unauthorized error")
	// `WrapIf` add stack-trace if not added before
	stackErr := errors.WrapIf(unAuthorizedErrMessage, message)

	unAuthorizedError := &unauthorizedError{
		CustomError: NewCustomError(stackErr, http.StatusUnauthorized, message),
	}

	return unAuthorizedError
}

// NewUnAuthorizedErrorWrap creates a new unauthorized error.
func NewUnAuthorizedErrorWrap(err error, message string) UnauthorizedError {
	if err == nil {
		return NewUnAuthorizedError(message)
	}

	// `WithMessage` doesn't add stack-trace at all
	unAuthorizedErrMessage := errors.WithMessage(err, "unauthorized error")
	// `WrapIf` add stack-trace if not added before
	stackErr := errors.WrapIf(unAuthorizedErrMessage, message)

	unAuthorizedError := &unauthorizedError{
		CustomError: NewCustomError(stackErr, http.StatusUnauthorized, message),
	}

	return unAuthorizedError
}

// unauthorizedError is a unauthorized error.
type unauthorizedError struct {
	CustomError
}

// UnauthorizedError is a unauthorized error.
type UnauthorizedError interface {
	CustomError
	isUnAuthorizedError()
}

// isUnAuthorizedError checks if the error is a unauthorized error.
func (u *unauthorizedError) isUnAuthorizedError() {
}

// IsUnAuthorizedError checks if the error is a unauthorized error.
func IsUnAuthorizedError(err error) bool {
	// https://github.com/golang/go/blob/master/src/net/error_windows.go#L10C2-L12C3
	// this doesn't work for a nested unauthorized error, and we should use errors.As for traversing errors in all levels
	var unauthorizedError UnauthorizedError
	if errors.As(err, &unauthorizedError) {
		return true
	}

	// us, ok := errors.Cause(err).(UnauthorizedError)
	if errors.As(err, &unauthorizedError) {
		return true
	}

	return false
}
