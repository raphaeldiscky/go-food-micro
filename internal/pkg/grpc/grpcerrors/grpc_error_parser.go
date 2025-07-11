// Package grpcerrors provides custom grpc errors.
package grpcerrors

import (
	"context"
	"database/sql"

	"emperror.dev/errors"
	"github.com/go-playground/validator"
	"google.golang.org/grpc/codes"

	"github.com/raphaeldiscky/go-food-micro/internal/pkg/constants"
	customErrors "github.com/raphaeldiscky/go-food-micro/internal/pkg/http/httperrors/customerrors"
	errorUtils "github.com/raphaeldiscky/go-food-micro/internal/pkg/utils/errorutils"
)

// https://github.com/grpc/grpc/blob/master/doc/http-grpc-status-mapping.md
// https://github.com/grpc/grpc/blob/master/doc/statuscodes.md

// ParseError is a function that parses an error.
func ParseError(err error) GrpcErr {
	customErr := customErrors.GetCustomError(err)
	var validatorErr validator.ValidationErrors
	stackTrace := errorUtils.ErrorsWithStack(err)

	if err != nil && customErr != nil {
		switch {
		case customErrors.IsDomainError(err, customErr.Status()):
			return NewDomainGrpcError(
				//nolint:gosec // G115: integer overflow conversion int -> uint32
				codes.Code(customErr.Status()),
				customErr.Error(),
				stackTrace,
			)
		case customErrors.IsApplicationError(err, customErr.Status()):
			return NewApplicationGrpcError(
				//nolint:gosec // G115: integer overflow conversion int -> uint32
				codes.Code(customErr.Status()),
				customErr.Error(),
				stackTrace,
			)
		case customErrors.IsAPIError(err, customErr.Status()):
			return NewAPIGrpcError(
				//nolint:gosec // G115: integer overflow conversion int -> uint32
				codes.Code(customErr.Status()),
				customErr.Error(),
				stackTrace,
			)
		case customErrors.IsBadRequestError(err):
			return NewBadRequestGrpcError(customErr.Error(), stackTrace)
		case customErrors.IsNotFoundError(err):
			return NewNotFoundErrorGrpcError(customErr.Error(), stackTrace)
		case customErrors.IsValidationError(err):
			return NewValidationGrpcError(customErr.Error(), stackTrace)
		case customErrors.IsUnAuthorizedError(err):
			return NewUnAuthorizedErrorGrpcError(customErr.Error(), stackTrace)
		case customErrors.IsForbiddenError(err):
			return NewForbiddenGrpcError(customErr.Error(), stackTrace)
		case customErrors.IsConflictError(err):
			return NewConflictGrpcError(customErr.Error(), stackTrace)
		case customErrors.IsInternalServerError(err):
			return NewInternalServerGrpcError(customErr.Error(), stackTrace)
		case customErrors.IsCustomError(err):
			return NewGrpcError(
				//nolint:gosec // G115: integer overflow conversion int -> uint32
				codes.Code(customErr.Status()),
				//nolint:gosec // G115: integer overflow conversion int -> uint32
				codes.Code(customErr.Status()).String(),
				customErr.Error(),
				stackTrace,
			)
		case customErrors.IsUnMarshalingError(err):
			return NewInternalServerGrpcError(customErr.Error(), stackTrace)
		case customErrors.IsMarshalingError(err):
			return NewInternalServerGrpcError(customErr.Error(), stackTrace)
		default:
			return NewInternalServerGrpcError(err.Error(), stackTrace)
		}
	} else if err != nil && customErr == nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return NewNotFoundErrorGrpcError(err.Error(), stackTrace)
		case errors.Is(err, context.DeadlineExceeded):
			return NewGrpcError(
				codes.DeadlineExceeded,
				constants.ErrRequestTimeoutTitle,
				err.Error(),
				stackTrace,
			)
		case errors.As(err, &validatorErr):
			return NewValidationGrpcError(validatorErr.Error(), stackTrace)
		}
	}

	return nil
}
