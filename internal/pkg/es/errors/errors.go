package errors

import (
	customErrors "github.com/raphaeldiscky/go-food-micro/internal/pkg/http/httperrors/customerrors"

	"emperror.dev/errors"
)

var (
	EventAlreadyExistsError = customErrors.NewConflictError(
		"domain_events event already exists in event registry",
	)
	InvalidEventTypeError = errors.New("invalid event type")
)
