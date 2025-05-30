// Package customerrors provides custom errors.
package customerrors

import (
	"net/http"

	"emperror.dev/errors"
)

// NewMarshalingError creates a new marshaling error.
func NewMarshalingError(message string) MarshalingError {
	// `NewPlain` doesn't add stack-trace at all
	marshalingErrMessage := errors.NewPlain("marshaling error")
	// `WrapIf` add stack-trace if not added before
	stackErr := errors.WrapIf(marshalingErrMessage, message)

	marshalingError := &marshalingError{
		CustomError: NewCustomError(stackErr, http.StatusInternalServerError, message),
	}

	return marshalingError
}

// NewMarshalingErrorWrap creates a new marshaling error.
func NewMarshalingErrorWrap(err error, message string) MarshalingError {
	if err == nil {
		return NewMarshalingError(message)
	}

	// `WithMessage` doesn't add stack-trace at all
	marshalingErrMessage := errors.WithMessage(err, "marshaling error")
	// `WrapIf` add stack-trace if not added before
	stackErr := errors.WrapIf(marshalingErrMessage, message)

	marshalingError := &marshalingError{
		CustomError: NewCustomError(stackErr, http.StatusInternalServerError, message),
	}

	return marshalingError
}

// marshalingError is a marshaling error.
type marshalingError struct {
	CustomError
}

// MarshalingError is a marshaling error.
type MarshalingError interface {
	InternalServerError
	isMarshalingError()
}

// isMarshalingError checks if the error is a marshaling error.
func (m *marshalingError) isMarshalingError() {
}

// isInternalServerError checks if the error is a internal server error.
func (m *marshalingError) isInternalServerError() {
}

// IsMarshalingError checks if the error is a marshaling error.
func IsMarshalingError(err error) bool {
	var marshalingErr MarshalingError

	// https://github.com/golang/go/blob/master/src/net/error_windows.go#L10C2-L12C3
	// this doesn't work for a nested marshaling error, and we should use errors.As for traversing errors in all levels
	var marshalingError MarshalingError
	if errors.As(err, &marshalingError) {
		return true
	}

	// us, ok := errors.Cause(err).(MarshalingError)
	if errors.As(err, &marshalingErr) {
		return true
	}

	return false
}
