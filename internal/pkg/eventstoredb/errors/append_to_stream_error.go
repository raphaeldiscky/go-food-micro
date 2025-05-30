// Package errors provides a append to stream error.
package errors

import (
	"fmt"

	"emperror.dev/errors"

	customErrors "github.com/raphaeldiscky/go-food-micro/internal/pkg/http/httperrors/customerrors"
)

// appendToStreamError is a struct that represents a append to stream error.
type appendToStreamError struct {
	customErrors.BadRequestError
}

// AppendToStreamError is a interface that represents a append to stream error.
type AppendToStreamError interface {
	customErrors.BadRequestError
	IsAppendToStreamError() bool
}

// NewAppendToStreamError creates a new append to stream error.
func NewAppendToStreamError(err error, streamID string) error {
	bad := customErrors.NewBadRequestErrorWrap(
		err,
		fmt.Sprintf("unable to append events to stream %s", streamID),
	)
	customErr := customErrors.GetCustomError(bad)

	badRequestErr, ok := customErr.(customErrors.BadRequestError)
	if !ok {
		return errors.Wrap(
			err,
			fmt.Sprintf("failed to convert error to BadRequestError: %v", customErr),
		)
	}

	br := &appendToStreamError{
		BadRequestError: badRequestErr,
	}

	return errors.WithStackIf(br)
}

// IsAppendToStreamError checks if the error is a append to stream error.
func (err *appendToStreamError) IsAppendToStreamError() bool {
	return true
}

// IsAppendToStreamError checks if the error is a append to stream error.
func IsAppendToStreamError(err error) bool {
	var an AppendToStreamError
	if errors.As(err, &an) {
		return an.IsAppendToStreamError()
	}

	return false
}
