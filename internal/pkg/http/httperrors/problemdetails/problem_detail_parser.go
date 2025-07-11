// Package problemdetails provides problem details.
package problemdetails

import (
	"context"
	"database/sql"
	"net/http"
	"reflect"

	"emperror.dev/errors"
	"github.com/go-playground/validator"

	"github.com/raphaeldiscky/go-food-micro/internal/pkg/constants"
	customErrors "github.com/raphaeldiscky/go-food-micro/internal/pkg/http/httperrors/customerrors"
	typeMapper "github.com/raphaeldiscky/go-food-micro/internal/pkg/reflection/typemapper"
	errorUtils "github.com/raphaeldiscky/go-food-micro/internal/pkg/utils/errorutils"
)

// ProblemDetailParser is a struct that represents a problem detail parser.
type ProblemDetailParser struct {
	internalErrors map[reflect.Type]func(err error) ProblemDetailErr
}

// ErrorParserFunc is a function that parses an error and returns a problem detail error.
type ErrorParserFunc func(err error) ProblemDetailErr

// NewProblemDetailParser creates a new problem detail parser.
func NewProblemDetailParser(
	builder func(builder *OptionBuilder),
) *ProblemDetailParser {
	optionBuilder := NewOptionBuilder()
	builder(optionBuilder)
	items := optionBuilder.Build()

	return &ProblemDetailParser{internalErrors: items}
}

// ResolveError resolves an error and returns a problem detail error.
func (p *ProblemDetailParser) ResolveError(err error) ProblemDetailErr {
	errType := typeMapper.GetReflectType(err)
	problem := p.internalErrors[errType]
	if problem != nil {
		return problem(err)
	}

	return nil
}

// ParseError parses an error and returns a problem detail error.
func ParseError(err error) ProblemDetailErr {
	stackTrace := errorUtils.ErrorsWithStack(err)
	customErr := customErrors.GetCustomError(err)
	var validatorErr validator.ValidationErrors

	if err != nil && customErr != nil {
		switch {
		case customErrors.IsDomainError(err, customErr.Status()):
			return NewDomainProblemDetail(
				customErr.Status(),
				customErr.Error(),
				stackTrace,
			)
		case customErrors.IsApplicationError(err, customErr.Status()):
			return NewApplicationProblemDetail(
				customErr.Status(),
				customErr.Error(),
				stackTrace,
			)
		case customErrors.IsAPIError(err, customErr.Status()):
			return NewAPIProblemDetail(
				customErr.Status(),
				customErr.Error(),
				stackTrace,
			)
		case customErrors.IsBadRequestError(err):
			return NewBadRequestProblemDetail(customErr.Error(), stackTrace)
		case customErrors.IsNotFoundError(err):
			return NewNotFoundErrorProblemDetail(customErr.Error(), stackTrace)
		case customErrors.IsValidationError(err):
			return NewValidationProblemDetail(customErr.Error(), stackTrace)
		case customErrors.IsUnAuthorizedError(err):
			return NewUnAuthorizedErrorProblemDetail(
				customErr.Error(),
				stackTrace,
			)
		case customErrors.IsForbiddenError(err):
			return NewForbiddenProblemDetail(customErr.Error(), stackTrace)
		case customErrors.IsConflictError(err):
			return NewConflictProblemDetail(customErr.Error(), stackTrace)
		case customErrors.IsInternalServerError(err):
			return NewInternalServerProblemDetail(customErr.Error(), stackTrace)
		case customErrors.IsCustomError(err):
			return NewProblemDetailFromCodeAndDetail(
				customErr.Status(),
				customErr.Error(),
				stackTrace,
			)
		case customErrors.IsUnMarshalingError(err):
			return NewInternalServerProblemDetail(err.Error(), stackTrace)
		case customErrors.IsMarshalingError(err):
			return NewInternalServerProblemDetail(err.Error(), stackTrace)

		default:
			return NewInternalServerProblemDetail(err.Error(), stackTrace)
		}
	} else if err != nil && customErr == nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return NewNotFoundErrorProblemDetail(err.Error(), stackTrace)
		case errors.Is(err, context.DeadlineExceeded):
			return NewProblemDetail(
				http.StatusRequestTimeout,
				constants.ErrRequestTimeoutTitle,
				err.Error(),
				stackTrace,
			)
		case errors.As(err, &validatorErr):
			return NewValidationProblemDetail(validatorErr.Error(), stackTrace)
		default:
			return NewInternalServerProblemDetail(err.Error(), stackTrace)
		}
	}

	return nil
}
