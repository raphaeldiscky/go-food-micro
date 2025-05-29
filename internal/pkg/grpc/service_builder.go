// Package grpc provides a grpc service builder.
package grpc

import (
	"google.golang.org/grpc"
)

// GrpcServiceBuilder is a struct that represents a grpc service builder.
type GrpcServiceBuilder struct {
	server *grpc.Server
}

// NewGrpcServiceBuilder is a function that creates a new grpc service builder.
func NewGrpcServiceBuilder(server *grpc.Server) *GrpcServiceBuilder {
	return &GrpcServiceBuilder{server: server}
}

// RegisterRoutes is a function that registers the routes.
func (r *GrpcServiceBuilder) RegisterRoutes(builder func(s *grpc.Server)) *GrpcServiceBuilder {
	builder(r.server)

	return r
}

// Build is a function that builds the grpc server.
func (r *GrpcServiceBuilder) Build() *grpc.Server {
	return r.server
}
