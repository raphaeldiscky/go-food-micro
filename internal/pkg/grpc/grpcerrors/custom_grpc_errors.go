// Package grpcerrors provides custom grpc errors.
package grpcerrors

import (
	"time"

	"google.golang.org/grpc/codes"

	"github.com/raphaeldiscky/go-food-micro/internal/pkg/constants"
)

// NewValidationGrpcError is a function that creates a new validation grpc error.
func NewValidationGrpcError(detail string, stackTrace string) GrpcErr {
	validationError := &grpcErr{
		Title:      constants.ErrBadRequestTitle,
		Detail:     detail,
		Status:     codes.InvalidArgument,
		Timestamp:  time.Now(),
		StackTrace: stackTrace,
	}

	return validationError
}

// NewConflictGrpcError is a function that creates a new conflict grpc error.
func NewConflictGrpcError(detail string, stackTrace string) GrpcErr {
	return &grpcErr{
		Title:      constants.ErrConflictTitle,
		Detail:     detail,
		Status:     codes.AlreadyExists,
		Timestamp:  time.Now(),
		StackTrace: stackTrace,
	}
}

// NewBadRequestGrpcError is a function that creates a new bad request grpc error.
func NewBadRequestGrpcError(detail string, stackTrace string) GrpcErr {
	return &grpcErr{
		Title:      constants.ErrBadRequestTitle,
		Detail:     detail,
		Status:     codes.InvalidArgument,
		Timestamp:  time.Now(),
		StackTrace: stackTrace,
	}
}

// NewNotFoundErrorGrpcError is a function that creates a new not found grpc error.
func NewNotFoundErrorGrpcError(detail string, stackTrace string) GrpcErr {
	return &grpcErr{
		Title:      constants.ErrNotFoundTitle,
		Detail:     detail,
		Status:     codes.NotFound,
		Timestamp:  time.Now(),
		StackTrace: stackTrace,
	}
}

// NewUnAuthorizedErrorGrpcError is a function that creates a new unauthorized grpc error.
func NewUnAuthorizedErrorGrpcError(detail string, stackTrace string) GrpcErr {
	return &grpcErr{
		Title:      constants.ErrUnauthorizedTitle,
		Detail:     detail,
		Status:     codes.Unauthenticated,
		Timestamp:  time.Now(),
		StackTrace: stackTrace,
	}
}

// NewForbiddenGrpcError is a function that creates a new forbidden grpc error.
func NewForbiddenGrpcError(detail string, stackTrace string) GrpcErr {
	return &grpcErr{
		Title:      constants.ErrForbiddenTitle,
		Detail:     detail,
		Status:     codes.PermissionDenied,
		Timestamp:  time.Now(),
		StackTrace: stackTrace,
	}
}

// NewInternalServerGrpcError is a function that creates a new internal server grpc error.
func NewInternalServerGrpcError(detail string, stackTrace string) GrpcErr {
	return &grpcErr{
		Title:      constants.ErrInternalServerErrorTitle,
		Detail:     detail,
		Status:     codes.Internal,
		Timestamp:  time.Now(),
		StackTrace: stackTrace,
	}
}

// NewDomainGrpcError is a function that creates a new domain grpc error.
func NewDomainGrpcError(status codes.Code, detail string, stackTrace string) GrpcErr {
	return &grpcErr{
		Title:      constants.ErrDomainTitle,
		Detail:     detail,
		Status:     status,
		Timestamp:  time.Now(),
		StackTrace: stackTrace,
	}
}

// NewApplicationGrpcError is a function that creates a new application grpc error.
func NewApplicationGrpcError(status codes.Code, detail string, stackTrace string) GrpcErr {
	return &grpcErr{
		Title:      constants.ErrApplicationTitle,
		Detail:     detail,
		Status:     status,
		Timestamp:  time.Now(),
		StackTrace: stackTrace,
	}
}

// NewApiGrpcError is a function that creates a new api grpc error.
func NewApiGrpcError(status codes.Code, detail string, stackTrace string) GrpcErr {
	return &grpcErr{
		Title:      constants.ErrAPITitle,
		Detail:     detail,
		Status:     status,
		Timestamp:  time.Now(),
		StackTrace: stackTrace,
	}
}
