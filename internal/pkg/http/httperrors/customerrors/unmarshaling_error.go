// Package customerrors provides custom errors.
package customerrors

import (
	"net/http"

	"emperror.dev/errors"
)

// NewUnMarshalingError creates a new unMarshaling error.
func NewUnMarshalingError(message string) UnMarshalingError {
	// `NewPlain` doesn't add stack-trace at all
	unMarshalingErrMessage := errors.NewPlain("unMarshaling error")
	// `WrapIf` add stack-trace if not added before
	stackErr := errors.WrapIf(unMarshalingErrMessage, message)

	unMarshalingError := &unMarshalingError{
		CustomError: NewCustomError(stackErr, http.StatusInternalServerError, message),
	}

	return unMarshalingError
}

// NewUnMarshalingErrorWrap creates a new unMarshaling error.
func NewUnMarshalingErrorWrap(err error, message string) UnMarshalingError {
	if err == nil {
		return NewUnMarshalingError(message)
	}

	// `WithMessage` doesn't add stack-trace at all
	unMarshalingErrMessage := errors.WithMessage(err, "unMarshaling error")
	// `WrapIf` add stack-trace if not added before
	stackErr := errors.WrapIf(unMarshalingErrMessage, message)

	unMarshalingError := &unMarshalingError{
		CustomError: NewCustomError(stackErr, http.StatusInternalServerError, message),
	}

	return unMarshalingError
}

// unMarshalingError is a unMarshaling error.
type unMarshalingError struct {
	CustomError
}

// UnMarshalingError is a unMarshaling error.
type UnMarshalingError interface {
	InternalServerError
	isUnMarshalingError()
}

// isUnMarshalingError checks if the error is a unMarshaling error.
func (u *unMarshalingError) isUnMarshalingError() {
}

// isInternalServerError checks if the error is a internal server error.
func (u *unMarshalingError) isInternalServerError() {
}

// IsUnMarshalingError checks if the error is a unMarshaling error.
func IsUnMarshalingError(err error) bool {
	// https://github.com/golang/go/blob/master/src/net/error_windows.go#L10C2-L12C3
	// this doesn't work for a nested unMarshaling error, and we should use errors.As for traversing errors in all levels
	var unMarshalingError UnMarshalingError
	if errors.As(err, &unMarshalingError) {
		return true
	}

	// us, ok := errors.Cause(err).(UnMarshalingError)
	if errors.As(err, &unMarshalingError) {
		return true
	}

	return false
}
