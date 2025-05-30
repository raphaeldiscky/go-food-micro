// Package grpc provides a grpc server.
package grpc

import (
	"fmt"
	"net"
	"time"

	"emperror.dev/errors"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/reflection"

	grpcMiddleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpcRecovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpcCtxTags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	googleGrpc "google.golang.org/grpc"

	"github.com/raphaeldiscky/go-food-micro/internal/pkg/grpc/config"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/grpc/handlers/otel"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/grpc/interceptors"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/logger"
)

const (
	maxConnectionIdle = 5
	gRPCTimeout       = 15
	maxConnectionAge  = 5
	gRPCTime          = 10
)

// GrpcServer is an interface that represents a grpc server.
type GrpcServer interface {
	RunGrpcServer(configGrpc ...func(grpcServer *googleGrpc.Server)) error
	GracefulShutdown()
	GetCurrentGrpcServer() *googleGrpc.Server
	GrpcServiceBuilder() *GrpcServiceBuilder
}

// grpcServer is a struct that represents a grpc server.
type grpcServer struct {
	server         *googleGrpc.Server
	config         *config.GrpcOptions
	log            logger.Logger
	serviceName    string
	serviceBuilder *GrpcServiceBuilder
}

// NewGrpcServer is a function that creates a new grpc server.
func NewGrpcServer(
	config *config.GrpcOptions,
	logger logger.Logger,
) GrpcServer {
	unaryServerInterceptors := []googleGrpc.UnaryServerInterceptor{
		interceptors.UnaryServerInterceptor(),
		grpcCtxTags.UnaryServerInterceptor(),
		grpcRecovery.UnaryServerInterceptor(),
	}
	streamServerInterceptors := []googleGrpc.StreamServerInterceptor{
		interceptors.StreamServerInterceptor(),
	}

	s := googleGrpc.NewServer(
		// https://github.com/open-telemetry/opentelemetry-go-contrib/issues/2840
		// https://github.com/open-telemetry/opentelemetry-go-contrib/pull/3002
		// https://github.com/open-telemetry/opentelemetry-go-contrib/blob/main/instrumentation/google.golang.org/grpc/otelgrpc/doc.go
		// https://github.com/open-telemetry/opentelemetry-go-contrib/blob/main/instrumentation/google.golang.org/grpc/otelgrpc/example/server/main.go#L143C3-L143C50
		googleGrpc.StatsHandler(otelgrpc.NewServerHandler()),
		googleGrpc.StatsHandler(otel.NewServerHandler()),

		googleGrpc.KeepaliveParams(keepalive.ServerParameters{
			MaxConnectionIdle: maxConnectionIdle * time.Minute,
			Timeout:           gRPCTimeout * time.Second,
			MaxConnectionAge:  maxConnectionAge * time.Minute,
			Time:              gRPCTime * time.Minute,
		}),
		// https://github.com/open-telemetry/opentelemetry-go-contrib/tree/00b796d0cdc204fa5d864ec690b2ee9656bb5cfc/instrumentation/google.golang.org/grpc/otelgrpc
		// github.com/grpc-ecosystem/go-grpc-middleware
		googleGrpc.StreamInterceptor(grpcMiddleware.ChainStreamServer(
			streamServerInterceptors...,
		)),
		googleGrpc.UnaryInterceptor(grpcMiddleware.ChainUnaryServer(
			unaryServerInterceptors...,
		)),
	)
	healthServer := health.NewServer()
	grpc_health_v1.RegisterHealthServer(s, healthServer)
	healthServer.SetServingStatus(
		config.Name,
		grpc_health_v1.HealthCheckResponse_SERVING,
	)

	return &grpcServer{
		server:         s,
		config:         config,
		log:            logger,
		serviceName:    config.Name,
		serviceBuilder: NewGrpcServiceBuilder(s),
	}
}

// RunGrpcServer is a function that runs the grpc server.
func (s *grpcServer) RunGrpcServer(
	configGrpc ...func(grpcServer *googleGrpc.Server),
) error {
	l, err := net.Listen("tcp", s.config.Port)
	if err != nil {
		return errors.WrapIf(err, "net.Listen")
	}

	if len(configGrpc) > 0 {
		grpcFunc := configGrpc[0]
		if grpcFunc != nil {
			grpcFunc(s.server)
		}
	}

	if s.config.Development {
		reflection.Register(s.server)
	}

	s.log.Infof(
		"[grpcServer.RunGrpcServer] Writer gRPC server is listening on port: %s",
		s.config.Port,
	)

	err = s.server.Serve(l)
	if err != nil {
		s.log.Error(
			fmt.Sprintf(
				"[grpcServer_RunGrpcServer.Serve] grpc server serve error: %+v",
				err,
			),
		)
	}

	return err
}

// GrpcServiceBuilder is a function that returns the grpc service builder.
func (s *grpcServer) GrpcServiceBuilder() *GrpcServiceBuilder {
	return s.serviceBuilder
}

// GetCurrentGrpcServer is a function that returns the current grpc server.
func (s *grpcServer) GetCurrentGrpcServer() *googleGrpc.Server {
	return s.server
}

// GracefulShutdown is a function that gracefully shuts down the grpc server.
func (s *grpcServer) GracefulShutdown() {
	s.server.Stop()
	s.server.GracefulStop()
}
