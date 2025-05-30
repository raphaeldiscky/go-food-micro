// Package customerrors provides custom errors.
package customerrors

import (
	"net/http"

	"emperror.dev/errors"
)

// NewApplicationError creates a new application error.
func NewApplicationError(message string) ApplicationError {
	return NewApplicationErrorWithCode(message, http.StatusInternalServerError)
}

// NewApplicationErrorWithCode creates a new application error with a code.
func NewApplicationErrorWithCode(message string, code int) ApplicationError {
	// `NewPlain` doesn't add stack-trace at all
	applicationErrMessage := errors.NewPlain("application error")
	// `WrapIf` add stack-trace if not added before
	stackErr := errors.WrapIf(applicationErrMessage, message)

	applicationError := &applicationError{
		CustomError: NewCustomError(stackErr, code, message),
	}

	return applicationError
}

// NewApplicationErrorWrap creates a new application error wrap.
func NewApplicationErrorWrap(err error, message string) ApplicationError {
	return NewApplicationErrorWrapWithCode(err, http.StatusInternalServerError, message)
}

// NewApplicationErrorWrapWithCode creates a new application error wrap with a code.
func NewApplicationErrorWrapWithCode(
	err error,
	code int,
	message string,
) ApplicationError {
	if err == nil {
		return NewApplicationErrorWithCode(message, code)
	}

	// `WithMessage` doesn't add stack-trace at all
	applicationErrMessage := errors.WithMessage(err, "application error")
	// `WrapIf` add stack-trace if not added before
	stackErr := errors.WrapIf(applicationErrMessage, message)

	applicationError := &applicationError{
		CustomError: NewCustomError(stackErr, code, message),
	}

	return applicationError
}

// applicationError is a struct that represents a application error.
type applicationError struct {
	CustomError
}

// ApplicationError is a contract that represents a application error.
type ApplicationError interface {
	CustomError
	isApplicationError()
}

func (a *applicationError) isApplicationError() {
}

// IsApplicationError checks if the error is a application error.
func IsApplicationError(err error, code int) bool {
	// https://github.com/golang/go/blob/master/src/net/error_windows.go#L10C2-L12C3
	// this doesn't work for a nested application error, and we should use errors.As for traversing errors in all levels
	var applicationError ApplicationError
	if errors.As(err, &applicationError) {
		return true
	}

	// us, ok := errors.Cause(err).(ApplicationError)
	if errors.As(err, &applicationError) {
		return applicationError.Status() == code
	}

	return false
}
