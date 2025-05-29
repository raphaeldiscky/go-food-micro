// Package errors provides a delete stream error.
package errors

import (
	"fmt"

	"emperror.dev/errors"

	customErrors "github.com/raphaeldiscky/go-food-micro/internal/pkg/http/httperrors/customerrors"
)

// deleteStreamError is a struct that represents a delete stream error.
type deleteStreamError struct {
	customErrors.InternalServerError
}

// DeleteStreamError is a interface that represents a delete stream error.
type DeleteStreamError interface {
	customErrors.InternalServerError
	IsDeleteStreamError() bool
}

// NewDeleteStreamError creates a new delete stream error.
func NewDeleteStreamError(err error, streamID string) error {
	internal := customErrors.NewInternalServerErrorWrap(
		err,
		fmt.Sprintf("unable to delete stream %s", streamID),
	)
	customErr := customErrors.GetCustomError(internal)

	br := &deleteStreamError{
		InternalServerError: customErr.(customErrors.InternalServerError),
	}

	return errors.WithStackIf(br)
}

// IsDeleteStreamError checks if the error is a delete stream error.
func (err *deleteStreamError) IsDeleteStreamError() bool {
	return true
}

// IsDeleteStreamError checks if the error is a delete stream error.
func IsDeleteStreamError(err error) bool {
	var ds DeleteStreamError
	if errors.As(err, &ds) {
		return ds.IsDeleteStreamError()
	}

	return false
}
