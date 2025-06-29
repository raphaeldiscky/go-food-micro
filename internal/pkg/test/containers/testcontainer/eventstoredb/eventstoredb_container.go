// Package eventstoredb provides a eventstoredb container.
package eventstoredb

import (
	"context"
	"fmt"
	"testing"
	"time"

	"emperror.dev/errors"
	"github.com/docker/go-connections/nat"
	"github.com/testcontainers/testcontainers-go/wait"

	kdb "github.com/kurrent-io/KurrentDB-Client-Go/kurrentdb"
	testcontainers "github.com/testcontainers/testcontainers-go"

	"github.com/raphaeldiscky/go-food-micro/internal/pkg/eventstoredb"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/eventstoredb/config"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/logger"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/test/containers/contracts"
)

// eventstoredbTestContainers represents a eventstoredb test containers.
type eventstoredbTestContainers struct {
	container      testcontainers.Container
	defaultOptions *contracts.EventstoreDBContainerOptions
	logger         logger.Logger
}

// NewEventStoreDBTestContainers creates a new eventstoredb test containers.
func NewEventStoreDBTestContainers(l logger.Logger) contracts.EventstoreDBContainer {
	return &eventstoredbTestContainers{
		defaultOptions: &contracts.EventstoreDBContainerOptions{
			Ports:   []string{"2113/tcp", "1113/tcp"},
			Host:    "localhost",
			TcpPort: 1113,
			// HTTP is the primary protocol for EventStoreDB. It is used in gRPC communication and HTTP APIs (management, gossip and diagnostics).
			HttpPort:  2113,
			Tag:       "latest",
			ImageName: "kurrentplatform/kurrentdb",
			Name:      "kurrentdb-testcontainers",
		},
		logger: l,
	}
}

// PopulateContainerOptions populates the container options.
func (g *eventstoredbTestContainers) PopulateContainerOptions(
	ctx context.Context,
	t *testing.T,
	options ...*contracts.EventstoreDBContainerOptions,
) (*config.EventStoreDbOptions, error) {
	t.Helper()

	containerReq := g.getRunOptions(options...)

	// @TODO: Using Parallel Container
	dbContainer, err := testcontainers.GenericContainer(
		ctx,
		testcontainers.GenericContainerRequest{
			ContainerRequest: containerReq,
			Started:          true,
		})
	if err != nil {
		return nil, err
	}

	// Clean up the container after the test is complete
	t.Cleanup(func() {
		if err := dbContainer.Terminate(ctx); err != nil {
			t.Fatalf("failed to terminate container: %s", err)
		}
	})

	// get a free random host port for http and grpc port for eventstoredb
	httpPort, err := dbContainer.MappedPort(ctx, nat.Port(g.defaultOptions.Ports[0]))
	if err != nil {
		return nil, err
	}
	g.defaultOptions.HttpPort = httpPort.Int()
	g.logger.Infof("eventstoredb http and grpc port is: %d", httpPort.Int())

	// get a free random host port for tcp port eventstoredb
	tcpPort, err := dbContainer.MappedPort(ctx, nat.Port(g.defaultOptions.Ports[1]))
	if err != nil {
		return nil, err
	}
	g.defaultOptions.TcpPort = tcpPort.Int()

	host, err := dbContainer.Host(ctx)
	if err != nil {
		return nil, err
	}

	g.container = dbContainer

	option := &config.EventStoreDbOptions{
		Host:     host,
		TcpPort:  g.defaultOptions.TcpPort,
		HttpPort: g.defaultOptions.HttpPort,
	}

	return option, nil
}

// Start starts the container.
func (g *eventstoredbTestContainers) Start(
	ctx context.Context,
	t *testing.T,
	options ...*contracts.EventstoreDBContainerOptions,
) (*kdb.Client, error) {
	t.Helper()
	eventstoredbOptions, err := g.PopulateContainerOptions(ctx, t, options...)
	if err != nil {
		return nil, err
	}

	return eventstoredb.NewEventStoreDB(eventstoredbOptions)
}

// Cleanup cleans up the container.
func (g *eventstoredbTestContainers) Cleanup(ctx context.Context) error {
	if err := g.container.Terminate(ctx); err != nil {
		return errors.WrapIf(err, "failed to terminate container: %s")
	}

	return nil
}

// getRunOptions gets the run options.
func (g *eventstoredbTestContainers) getRunOptions(
	opts ...*contracts.EventstoreDBContainerOptions,
) testcontainers.ContainerRequest {
	if len(opts) > 0 {
		g.updateOptions(opts[0])
	}

	return testcontainers.ContainerRequest{
		Image:        fmt.Sprintf("%s:%s", g.defaultOptions.ImageName, g.defaultOptions.Tag),
		ExposedPorts: g.defaultOptions.Ports,
		WaitingFor: wait.ForListeningPort(nat.Port(g.defaultOptions.Ports[0])).
			WithPollInterval(2 * time.Second),
		Hostname: g.defaultOptions.Host,
		// we use `KURRENTDB_IN_MEM` for use eventstoredb in-memory mode in tests
		Env: map[string]string{
			"KURRENTDB_START_STANDARD_PROJECTIONS": "false",
			"KURRENTDB_INSECURE":                   "true",
			"KURRENTDB_ENABLE_ATOM_PUB_OVER_HTTP":  "true",
			"KURRENTDB_MEM_DB":                     "true",
		},
	}
}

func (g *eventstoredbTestContainers) updateOptions(option *contracts.EventstoreDBContainerOptions) {
	if option.ImageName != "" {
		g.defaultOptions.ImageName = option.ImageName
	}
	if option.Host != "" {
		g.defaultOptions.Host = option.Host
	}
	if len(option.Ports) > 0 {
		g.defaultOptions.Ports = option.Ports
	}
	if option.Tag != "" {
		g.defaultOptions.Tag = option.Tag
	}
}
