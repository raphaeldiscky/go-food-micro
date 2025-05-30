// Package errors provides a read stream error.
package errors

import (
	"emperror.dev/errors"

	customErrors "github.com/raphaeldiscky/go-food-micro/internal/pkg/http/httperrors/customerrors"
)

// readStreamError is a struct that represents a read stream error.
type readStreamError struct {
	customErrors.InternalServerError
}

// ReadStreamError is a interface that represents a read stream error.
type ReadStreamError interface {
	customErrors.InternalServerError
	IsReadStreamError() bool
}

// NewReadStreamError creates a new read stream error.
func NewReadStreamError(err error) error {
	internal := customErrors.NewInternalServerErrorWrap(err, "unable to read events from stream")
	customErr := customErrors.GetCustomError(internal)

	internalServerErr, ok := customErr.(customErrors.InternalServerError)
	if !ok {
		return errors.Wrap(
			err,
			"failed to convert error to InternalServerError",
		)
	}

	br := &readStreamError{
		InternalServerError: internalServerErr,
	}

	return errors.WithStackIf(br)
}

func (err *readStreamError) IsReadStreamError() bool {
	return true
}

// IsReadStreamError checks if the error is a read stream error.
func IsReadStreamError(err error) bool {
	var rs ReadStreamError
	if errors.As(err, &rs) {
		return rs.IsReadStreamError()
	}

	return false
}
