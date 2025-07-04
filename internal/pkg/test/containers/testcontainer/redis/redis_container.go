// Package redis provides a redis container.
package redis

import (
	"context"
	"fmt"
	"testing"
	"time"

	"emperror.dev/errors"
	"github.com/docker/go-connections/nat"
	"github.com/testcontainers/testcontainers-go/wait"

	redis "github.com/redis/go-redis/v9"
	testcontainers "github.com/testcontainers/testcontainers-go"

	"github.com/raphaeldiscky/go-food-micro/internal/pkg/logger"
	redis2 "github.com/raphaeldiscky/go-food-micro/internal/pkg/redis"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/test/containers/contracts"
)

// redisTestContainers represents a redis test containers.
type redisTestContainers struct {
	container      testcontainers.Container
	defaultOptions *contracts.RedisContainerOptions
	logger         logger.Logger
}

// NewRedisTestContainers creates a new redis test containers.
func NewRedisTestContainers(l logger.Logger) contracts.RedisContainer {
	return &redisTestContainers{
		defaultOptions: &contracts.RedisContainerOptions{
			Port:      "6379/tcp",
			Host:      "localhost",
			Database:  0,
			PoolSize:  300,
			Tag:       "latest",
			ImageName: "redis",
			Name:      "redis-testcontainers",
		},
		logger: l,
	}
}

// PopulateContainerOptions populates the container options.
func (g *redisTestContainers) PopulateContainerOptions(
	ctx context.Context,
	t *testing.T,
	options ...*contracts.RedisContainerOptions,
) (*redis2.RedisOptions, error) {
	t.Helper()
	// https://github.com/testcontainers/testcontainers-go
	// https://dev.to/remast/go-integration-tests-using-testcontainers-9o5
	containerReq := g.getRunOptions(options...)

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
	hostPort, err := dbContainer.MappedPort(
		ctx,
		nat.Port(g.defaultOptions.Port),
	)
	if err != nil {
		return nil, err
	}
	g.defaultOptions.HostPort = hostPort.Int()

	host, err := dbContainer.Host(ctx)
	if err != nil {
		return nil, err
	}

	isConnectable := isConnectable(ctx, g.logger, g.defaultOptions)
	if !isConnectable {
		return g.PopulateContainerOptions(context.Background(), t, options...)
	}

	g.container = dbContainer

	reidsOptions := &redis2.RedisOptions{
		Database: g.defaultOptions.Database,
		Host:     host,
		Port:     g.defaultOptions.HostPort,
		PoolSize: g.defaultOptions.PoolSize,
	}

	return reidsOptions, nil
}

// Cleanup cleans up the container.
func (g *redisTestContainers) Cleanup(ctx context.Context) error {
	if err := g.container.Terminate(ctx); err != nil {
		return errors.WrapIf(err, "failed to terminate container: %s")
	}

	return nil
}

// getRunOptions gets the run options.
func (g *redisTestContainers) getRunOptions(
	opts ...*contracts.RedisContainerOptions,
) testcontainers.ContainerRequest {
	if len(opts) > 0 {
		g.updateOptions(opts[0])
	}

	return testcontainers.ContainerRequest{
		Image: fmt.Sprintf(
			"%s:%s",
			g.defaultOptions.ImageName,
			g.defaultOptions.Tag,
		),
		ExposedPorts: []string{g.defaultOptions.Port},
		WaitingFor: wait.ForListeningPort(nat.Port(g.defaultOptions.Port)).
			WithPollInterval(2 * time.Second),
		Hostname: g.defaultOptions.Host,
		Env:      map[string]string{},
	}
}

func (g *redisTestContainers) updateOptions(option *contracts.RedisContainerOptions) {
	if option.ImageName != "" {
		g.defaultOptions.ImageName = option.ImageName
	}
	if option.Host != "" {
		g.defaultOptions.Host = option.Host
	}
	if option.Port != "" {
		g.defaultOptions.Port = option.Port
	}
	if option.Tag != "" {
		g.defaultOptions.Tag = option.Tag
	}
}

// isConnectable checks if the container is connectable.
func isConnectable(
	ctx context.Context,
	logger logger.Logger,
	options *contracts.RedisContainerOptions,
) bool {
	redisClient := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%d", options.Host, options.HostPort),
	})

	defer func() {
		if err := redisClient.Close(); err != nil {
			logger.Errorf("Error closing redis connection: %v", err)
		}
	}()

	err := redisClient.Ping(ctx).Err()
	if err != nil {
		// we should not use `t.Error` or `t.Errorf` for logging errors because it will `fail` our test at the end and, we just should use logs without error like log.Error (not log.Fatal)
		logger.Errorf(
			"Error in creating redis connection with %s:%d",
			options.Host,
			options.HostPort,
		)

		return false
	}

	logger.Infof(
		"Opened redis connection on host: %s:%d",
		options.Host,
		options.HostPort,
	)

	return true
}
