// Package errors provides a set of errors for the es package.
package errors

import (
	"emperror.dev/errors"

	customErrors "github.com/raphaeldiscky/go-food-micro/internal/pkg/http/httperrors/customerrors"
)

var (
	// EventAlreadyExistsError is a error that represents a event already exists.
	EventAlreadyExistsError = customErrors.NewConflictError(
		"domain_events event already exists in event registry",
	)
	// InvalidEventTypeError is a error that represents a invalid event type.
	InvalidEventTypeError = errors.New("invalid event type")
)
