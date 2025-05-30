// Package interceptors provides a error interceptor.
package interceptors

import (
	"context"

	"emperror.dev/errors"
	"google.golang.org/grpc"

	"github.com/raphaeldiscky/go-food-micro/internal/pkg/grpc/grpcerrors"
)

// UnaryServerInterceptor is a function that returns a problem-detail error to client.
func UnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		_ *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		resp, err := handler(ctx, req)

		var grpcErr grpcerrors.GrpcErr

		// if error was not `grpcErr` we will convert the error to a `grpcErr`
		if ok := errors.As(err, &grpcErr); !ok {
			grpcErr = grpcerrors.ParseError(err)
		}

		if grpcErr != nil {
			return nil, grpcErr.ToGrpcResponseErr()
		}

		return resp, err
	}
}

// StreamServerInterceptor is a function that returns a problem-detail error to client.
func StreamServerInterceptor() grpc.StreamServerInterceptor {
	return func(
		srv interface{},
		ss grpc.ServerStream,
		_ *grpc.StreamServerInfo,
		handler grpc.StreamHandler,
	) error {
		err := handler(srv, ss)

		var grpcErr grpcerrors.GrpcErr

		// if error was not `grpcErr` we will convert the error to a `grpcErr`
		if ok := errors.As(err, &grpcErr); !ok {
			grpcErr = grpcerrors.ParseError(err)
		}

		if grpcErr != nil {
			return grpcErr.ToGrpcResponseErr()
		}

		return err
	}
}
