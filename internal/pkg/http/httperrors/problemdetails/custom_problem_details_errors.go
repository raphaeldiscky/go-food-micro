// Package problemdetails provides problem details.
package problemdetails

import (
	"net/http"
	"time"

	"github.com/raphaeldiscky/go-food-micro/internal/pkg/constants"
)

// NewValidationProblemDetail creates a new validation problem detail.
func NewValidationProblemDetail(detail string, stackTrace string) ProblemDetailErr {
	validationError := &problemDetail{
		Title:      constants.ErrBadRequestTitle,
		Detail:     detail,
		Status:     http.StatusBadRequest,
		Type:       getDefaultType(http.StatusBadRequest),
		Timestamp:  time.Now(),
		StackTrace: stackTrace,
	}

	return validationError
}

// NewConflictProblemDetail creates a new conflict problem detail.
func NewConflictProblemDetail(detail string, stackTrace string) ProblemDetailErr {
	return &problemDetail{
		Title:      constants.ErrConflictTitle,
		Detail:     detail,
		Status:     http.StatusConflict,
		Type:       getDefaultType(http.StatusConflict),
		Timestamp:  time.Now(),
		StackTrace: stackTrace,
	}
}

// NewBadRequestProblemDetail creates a new bad request problem detail.
func NewBadRequestProblemDetail(detail string, stackTrace string) ProblemDetailErr {
	return &problemDetail{
		Title:      constants.ErrBadRequestTitle,
		Detail:     detail,
		Status:     http.StatusBadRequest,
		Type:       getDefaultType(http.StatusBadRequest),
		Timestamp:  time.Now(),
		StackTrace: stackTrace,
	}
}

// NewNotFoundErrorProblemDetail creates a new not found error problem detail.
func NewNotFoundErrorProblemDetail(detail string, stackTrace string) ProblemDetailErr {
	return &problemDetail{
		Title:      constants.ErrNotFoundTitle,
		Detail:     detail,
		Status:     http.StatusNotFound,
		Type:       getDefaultType(http.StatusNotFound),
		Timestamp:  time.Now(),
		StackTrace: stackTrace,
	}
}

// NewUnAuthorizedErrorProblemDetail creates a new unauthorized error problem detail.
func NewUnAuthorizedErrorProblemDetail(detail string, stackTrace string) ProblemDetailErr {
	return &problemDetail{
		Title:      constants.ErrUnauthorizedTitle,
		Detail:     detail,
		Status:     http.StatusUnauthorized,
		Type:       getDefaultType(http.StatusUnauthorized),
		Timestamp:  time.Now(),
		StackTrace: stackTrace,
	}
}

// NewForbiddenProblemDetail creates a new forbidden problem detail.
func NewForbiddenProblemDetail(detail string, stackTrace string) ProblemDetailErr {
	return &problemDetail{
		Title:      constants.ErrForbiddenTitle,
		Detail:     detail,
		Status:     http.StatusForbidden,
		Type:       getDefaultType(http.StatusForbidden),
		Timestamp:  time.Now(),
		StackTrace: stackTrace,
	}
}

// NewInternalServerProblemDetail creates a new internal server problem detail.
func NewInternalServerProblemDetail(detail string, stackTrace string) ProblemDetailErr {
	return &problemDetail{
		Title:      constants.ErrInternalServerErrorTitle,
		Detail:     detail,
		Status:     http.StatusInternalServerError,
		Type:       getDefaultType(http.StatusInternalServerError),
		Timestamp:  time.Now(),
		StackTrace: stackTrace,
	}
}

// NewDomainProblemDetail creates a new domain problem detail.
func NewDomainProblemDetail(status int, detail string, stackTrace string) ProblemDetailErr {
	return &problemDetail{
		Title:      constants.ErrDomainTitle,
		Detail:     detail,
		Status:     status,
		Type:       getDefaultType(http.StatusBadRequest),
		Timestamp:  time.Now(),
		StackTrace: stackTrace,
	}
}

// NewApplicationProblemDetail creates a new application problem detail.
func NewApplicationProblemDetail(status int, detail string, stackTrace string) ProblemDetailErr {
	return &problemDetail{
		Title:      constants.ErrApplicationTitle,
		Detail:     detail,
		Status:     status,
		Type:       getDefaultType(status),
		Timestamp:  time.Now(),
		StackTrace: stackTrace,
	}
}

// NewAPIProblemDetail creates a new api problem detail.
func NewAPIProblemDetail(status int, detail string, stackTrace string) ProblemDetailErr {
	return &problemDetail{
		Title:      constants.ErrAPITitle,
		Detail:     detail,
		Status:     status,
		Type:       getDefaultType(status),
		Timestamp:  time.Now(),
		StackTrace: stackTrace,
	}
}
