package errors

import (
	"emperror.dev/errors"

	customErrors "github.com/raphaeldiscky/go-food-micro/internal/pkg/http/httperrors/customerrors"
)

var (
	EventAlreadyExistsError = customErrors.NewConflictError(
		"domain_events event already exists in event registry",
	)
	InvalidEventTypeError = errors.New("invalid event type")
)
