// Package errors provides a truncate stream error.
package errors

import (
	"fmt"

	"emperror.dev/errors"

	customErrors "github.com/raphaeldiscky/go-food-micro/internal/pkg/http/httperrors/customerrors"
)

// truncateStreamError is a struct that represents a truncate stream error.
type truncateStreamError struct {
	customErrors.InternalServerError
}

// TruncateStreamError is a interface that represents a truncate stream error.
type TruncateStreamError interface {
	customErrors.InternalServerError
	IsTruncateStreamError() bool
}

// NewTruncateStreamError creates a new truncate stream error.
func NewTruncateStreamError(err error, streamID string) error {
	internal := customErrors.NewInternalServerErrorWrap(
		err,
		fmt.Sprintf("unable to truncate stream %s", streamID),
	)
	customErr := customErrors.GetCustomError(internal)

	internalServerErr, ok := customErr.(customErrors.InternalServerError)
	if !ok {
		return errors.Wrap(err, "failed to convert error to InternalServerError")
	}

	br := &truncateStreamError{
		InternalServerError: internalServerErr,
	}

	return errors.WithStackIf(br)
}

// IsTruncateStreamError checks if the error is a truncate stream error.
func (err *truncateStreamError) IsTruncateStreamError() bool {
	return true
}

// IsTruncateStreamError checks if the error is a truncate stream error.
func IsTruncateStreamError(err error) bool {
	var rs TruncateStreamError
	if errors.As(err, &rs) {
		return rs.IsTruncateStreamError()
	}

	return false
}
