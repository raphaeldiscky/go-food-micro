// Package customerrors provides custom errors.
package customerrors

import (
	"net/http"

	"emperror.dev/errors"
)

// NewValidationError creates a new validation error.
func NewValidationError(message string) ValidationError {
	// `NewPlain` doesn't add stack-trace at all
	validationErrMessage := errors.NewPlain("validation error")
	// `WrapIf` add stack-trace if not added before
	stackErr := errors.WrapIf(validationErrMessage, message)

	validationError := &validationError{
		CustomError: NewCustomError(stackErr, http.StatusBadRequest, message),
	}

	return validationError
}

// NewValidationErrorWrap creates a new validation error.
func NewValidationErrorWrap(err error, message string) ValidationError {
	if err == nil {
		return NewValidationError(message)
	}

	// `WithMessage` doesn't add stack-trace at all
	validationErrMessage := errors.WithMessage(err, "validation error")
	// `WrapIf` add stack-trace if not added before
	stackErr := errors.WrapIf(validationErrMessage, message)

	validationError := &validationError{
		CustomError: NewCustomError(stackErr, http.StatusBadRequest, message),
	}

	return validationError
}

// validationError is a validation error.
type validationError struct {
	CustomError
}

// ValidationError is a validation error.
type ValidationError interface {
	BadRequestError
	isValidationError()
}

// isValidationError checks if the error is a validation error.
func (v *validationError) isValidationError() {
}

// isBadRequestError checks if the error is a bad request error.
func (v *validationError) isBadRequestError() {
}

// IsValidationError checks if the error is a validation error.
func IsValidationError(err error) bool {
	// https://github.com/golang/go/blob/master/src/net/error_windows.go#L10C2-L12C3
	// this doesn't work for a nested validation error, and we should use errors.As for traversing errors in all levels
	var validationError ValidationError
	if errors.As(err, &validationError) {
		return true
	}

	// us, ok := errors.Cause(err).(ValidationError)
	if errors.As(err, &validationError) {
		return true
	}

	return false
}
