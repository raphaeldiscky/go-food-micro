// Package errors provides a stream not found error.
package errors

import (
	"fmt"

	"emperror.dev/errors"

	customErrors "github.com/raphaeldiscky/go-food-micro/internal/pkg/http/httperrors/customerrors"
)

// streamNotFoundError is a struct that represents a stream not found error.
type streamNotFoundError struct {
	customErrors.NotFoundError
}

// StreamNotFoundError is a interface that represents a stream not found error.
type StreamNotFoundError interface {
	customErrors.NotFoundError
	IsStreamNotFoundError() bool
}

// NewStreamNotFoundError creates a new stream not found error.
func NewStreamNotFoundError(err error, streamID string) error {
	notFound := customErrors.NewNotFoundErrorWrap(
		err,
		fmt.Sprintf("stream with streamId %s not found", streamID),
	)
	customErr := customErrors.GetCustomError(notFound)
	br := &streamNotFoundError{
		NotFoundError: customErr.(customErrors.NotFoundError),
	}

	return errors.WithStackIf(br)
}

// IsStreamNotFoundError checks if the error is a stream not found error.
func (err *streamNotFoundError) IsStreamNotFoundError() bool {
	return true
}

// IsStreamNotFoundError checks if the error is a stream not found error.
func IsStreamNotFoundError(err error) bool {
	var rs StreamNotFoundError
	if errors.As(err, &rs) {
		return rs.IsStreamNotFoundError()
	}

	return false
}
