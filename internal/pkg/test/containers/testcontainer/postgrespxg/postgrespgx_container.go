// Package postgrespxg provides a postgrespgx container.
package postgrespxg

import (
	"context"
	"fmt"
	"testing"
	"time"

	"emperror.dev/errors"
	"github.com/docker/go-connections/nat"
	"github.com/testcontainers/testcontainers-go/wait"

	testcontainers "github.com/testcontainers/testcontainers-go"

	"github.com/raphaeldiscky/go-food-micro/internal/pkg/logger"
	postgres "github.com/raphaeldiscky/go-food-micro/internal/pkg/postgrespgx"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/test/containers/contracts"
)

// postgresPgxTestContainers represents a postgrespgx test containers.
// https://github.com/testcontainers/testcontainers-go/issues/1359
// https://github.com/testcontainers/testcontainers-go/issues/1249
type postgresPgxTestContainers struct {
	container      testcontainers.Container
	defaultOptions *contracts.PostgresContainerOptions
	logger         logger.Logger
}

// NewPostgresPgxContainers creates a new postgrespgx test containers.
func NewPostgresPgxContainers(l logger.Logger) contracts.PostgresPgxContainer {
	return &postgresPgxTestContainers{
		defaultOptions: &contracts.PostgresContainerOptions{
			Database:  "test_db",
			Port:      "5432/tcp",
			Host:      "localhost",
			UserName:  "testcontainers",
			Password:  "testcontainers",
			Tag:       "latest",
			ImageName: "postgres",
			Name:      "postgresql-testcontainer",
		},
		logger: l,
	}
}

// PopulateContainerOptions populates the container options.
func (g *postgresPgxTestContainers) PopulateContainerOptions(
	ctx context.Context,
	t *testing.T,
	options ...*contracts.PostgresContainerOptions,
) (*postgres.PostgresPgxOptions, error) {
	t.Helper()

	// https://github.com/testcontainers/testcontainers-go
	// https://dev.to/remast/go-integration-tests-using-testcontainers-9o5
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

	// get a free random host hostPort
	hostPort, err := dbContainer.MappedPort(ctx, nat.Port(g.defaultOptions.Port))
	if err != nil {
		return nil, err
	}
	g.defaultOptions.HostPort = hostPort.Int()

	host, err := dbContainer.Host(ctx)
	if err != nil {
		return nil, err
	}

	g.container = dbContainer

	gormOptions := &postgres.PostgresPgxOptions{
		Port:     g.defaultOptions.HostPort,
		Host:     host,
		Password: g.defaultOptions.Password,
		DBName:   g.defaultOptions.Database,
		SSLMode:  false,
		User:     g.defaultOptions.UserName,
	}

	return gormOptions, nil
}

// Start starts the container.
func (g *postgresPgxTestContainers) Start(
	ctx context.Context,
	t *testing.T,
	options ...*contracts.PostgresContainerOptions,
) (*postgres.Pgx, error) {
	t.Helper()

	postgresPgxOptions, err := g.PopulateContainerOptions(ctx, t, options...)
	if err != nil {
		return nil, err
	}

	db, err := postgres.NewPgx(postgresPgxOptions)
	if err != nil {
		return nil, err
	}

	return db, nil
}

// Cleanup cleans up the container.
func (g *postgresPgxTestContainers) Cleanup(ctx context.Context) error {
	if err := g.container.Terminate(ctx); err != nil {
		return errors.WrapIf(err, "failed to terminate container: %s")
	}

	return nil
}

// updateOptions updates the default options with provided options.
func (g *postgresPgxTestContainers) updateOptions(option *contracts.PostgresContainerOptions) {
	if option.ImageName != "" {
		g.defaultOptions.ImageName = option.ImageName
	}
	if option.Host != "" {
		g.defaultOptions.Host = option.Host
	}
	if option.Port != "" {
		g.defaultOptions.Port = option.Port
	}
	if option.UserName != "" {
		g.defaultOptions.UserName = option.UserName
	}
	if option.Password != "" {
		g.defaultOptions.Password = option.Password
	}
	if option.Tag != "" {
		g.defaultOptions.Tag = option.Tag
	}
}

// getRunOptions gets the run options.
func (g *postgresPgxTestContainers) getRunOptions(
	opts ...*contracts.PostgresContainerOptions,
) testcontainers.ContainerRequest {
	if len(opts) > 0 {
		g.updateOptions(opts[0])
	}

	containerReq := testcontainers.ContainerRequest{
		Image:        fmt.Sprintf("%s:%s", g.defaultOptions.ImageName, g.defaultOptions.Tag),
		ExposedPorts: []string{g.defaultOptions.Port},
		WaitingFor: wait.ForListeningPort(nat.Port(g.defaultOptions.Port)).
			WithPollInterval(2 * time.Second),
		Hostname:   g.defaultOptions.Host,
		SkipReaper: true,
		Env: map[string]string{
			"POSTGRES_DB":       g.defaultOptions.Database,
			"POSTGRES_PASSWORD": g.defaultOptions.Password,
			"POSTGRES_USER":     g.defaultOptions.UserName,
		},
	}

	return containerReq
}
