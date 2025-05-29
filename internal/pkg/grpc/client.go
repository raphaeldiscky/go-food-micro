// Package grpc provides a grpc client.
package grpc

import (
	"fmt"
	"time"

	"emperror.dev/errors"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/raphaeldiscky/go-food-micro/internal/pkg/grpc/config"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/grpc/handlers/otel"
)

// grpcClient is a struct that represents a grpc client.
type grpcClient struct {
	conn *grpc.ClientConn
}

// GrpcClient is an interface that represents a grpc client.
type GrpcClient interface {
	GetGrpcConnection() *grpc.ClientConn
	Close() error
	// WaitForAvailableConnection waiting for grpc endpoint becomes ready in the given timeout
	WaitForAvailableConnection() error
}

// NewGrpcClient is a function that creates a new grpc client.
func NewGrpcClient(config *config.GrpcOptions) (GrpcClient, error) {
	// Grpc Client to call Grpc Server
	// https://sahansera.dev/building-grpc-client-go/
	// https://github.com/open-telemetry/opentelemetry-go-contrib/blob/df16f32df86b40077c9c90d06f33c4cdb6dd5afa/instrumentation/google.golang.org/grpc/otelgrpc/example_interceptor_test.go
	conn, err := grpc.Dial(fmt.Sprintf("%s%s", config.Host, config.Port),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		// https://github.com/open-telemetry/opentelemetry-go-contrib/blob/main/instrumentation/google.golang.org/grpc/otelgrpc/example/client/main.go#L47C3-L47C52
		// https://github.com/open-telemetry/opentelemetry-go-contrib/blob/main/instrumentation/google.golang.org/grpc/otelgrpc/doc.go
		grpc.WithStatsHandler(otelgrpc.NewClientHandler()),
		grpc.WithStatsHandler(otel.NewClientHandler()),
	)
	if err != nil {
		return nil, err
	}

	return &grpcClient{conn: conn}, err
}

// GetGrpcConnection is a function that returns the grpc connection.
func (g *grpcClient) GetGrpcConnection() *grpc.ClientConn {
	return g.conn
}

// Close is a function that closes the grpc connection.
func (g *grpcClient) Close() error {
	return g.conn.Close()
}

// WaitForAvailableConnection is a function that waits for the grpc connection to be available.
func (g *grpcClient) WaitForAvailableConnection() error {
	timeout := time.Second * 20

	err := waitUntilConditionMet(func() bool {
		return g.conn.GetState() == connectivity.Ready
	}, timeout)

	state := g.conn.GetState()
	fmt.Printf("grpc state is:%s\n", state)

	return err
}

// waitUntilConditionMet is a function that waits until a condition is met.
func waitUntilConditionMet(
	conditionToMet func() bool,
	timeout ...time.Duration,
) error {
	timeOutTime := 20 * time.Second
	if len(timeout) >= 0 && timeout != nil {
		timeOutTime = timeout[0]
	}

	startTime := time.Now()
	timeOutExpired := false
	meet := conditionToMet()
	for !meet {
		if timeOutExpired {
			return errors.New("grpc connection could not be established in the given timeout")
		}
		time.Sleep(time.Second * 2)
		meet = conditionToMet()
		timeOutExpired = time.Since(startTime) > timeOutTime
	}

	return nil
}
