// Package mongo provides a mongo container.
package mongo

import (
	"context"
	"fmt"
	"testing"
	"time"

	"emperror.dev/errors"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/go-connections/nat"
	"github.com/testcontainers/testcontainers-go/wait"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	testcontainers "github.com/testcontainers/testcontainers-go"

	"github.com/raphaeldiscky/go-food-micro/internal/pkg/logger"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/mongodb"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/test/containers/contracts"
)

const (
	connectTimeout  = 60 * time.Second
	maxConnIdleTime = 3 * time.Minute
	minPoolSize     = 20
	maxPoolSize     = 300
)

// mongoTestContainers represents a mongo test containers.
type mongoTestContainers struct {
	container      testcontainers.Container
	defaultOptions *contracts.MongoContainerOptions
	logger         logger.Logger
}

// NewMongoTestContainers creates a new mongo test containers.
func NewMongoTestContainers(l logger.Logger) contracts.MongoContainer {
	return &mongoTestContainers{
		defaultOptions: &contracts.MongoContainerOptions{
			Database:  "test_db",
			Port:      "27017/tcp",
			Host:      "localhost",
			UserName:  "testcontainers",
			Password:  "testcontainers",
			Tag:       "latest",
			ImageName: "mongo",
			Name:      "mongo-testcontainer",
		},
		logger: l,
	}
}

// PopulateContainerOptions populates the container options.
func (g *mongoTestContainers) PopulateContainerOptions(
	ctx context.Context,
	t *testing.T,
	options ...*contracts.MongoContainerOptions,
) (*mongodb.MongoDbOptions, error) {
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

	option := &mongodb.MongoDbOptions{
		User:     g.defaultOptions.UserName,
		Password: g.defaultOptions.Password,
		UseAuth:  false,
		Host:     host,
		Port:     g.defaultOptions.HostPort,
		Database: g.defaultOptions.Database,
	}

	return option, nil
}

// Cleanup cleans up the container.
func (g *mongoTestContainers) Cleanup(ctx context.Context) error {
	if err := g.container.Terminate(ctx); err != nil {
		return errors.WrapIf(err, "failed to terminate container: %s")
	}

	return nil
}

// updateOptions updates the default options with provided options.
func (g *mongoTestContainers) updateOptions(option *contracts.MongoContainerOptions) {
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
func (g *mongoTestContainers) getRunOptions(
	opts ...*contracts.MongoContainerOptions,
) testcontainers.ContainerRequest {
	if len(opts) > 0 {
		g.updateOptions(opts[0])
	}

	containerReq := testcontainers.ContainerRequest{
		Image: fmt.Sprintf(
			"%s:%s",
			g.defaultOptions.ImageName,
			g.defaultOptions.Tag,
		),
		ExposedPorts: []string{g.defaultOptions.Port},
		WaitingFor: wait.ForListeningPort(nat.Port(g.defaultOptions.Port)).
			WithPollInterval(2 * time.Second),
		Hostname: g.defaultOptions.Host,
		HostConfigModifier: func(hostConfig *container.HostConfig) {
			hostConfig.AutoRemove = true
		},
		Env: map[string]string{
			"MONGO_INITDB_ROOT_USERNAME": g.defaultOptions.UserName,
			"MONGO_INITDB_ROOT_PASSWORD": g.defaultOptions.Password,
		},
	}

	return containerReq
}

// isConnectable checks if the container is connectable.
func isConnectable(
	ctx context.Context,
	logger logger.Logger,
	mongoOptions *contracts.MongoContainerOptions,
) bool {
	uriAddress := fmt.Sprintf(
		"mongodb://%s:%s@%s:%d",
		mongoOptions.UserName,
		mongoOptions.Password,
		mongoOptions.Host,
		mongoOptions.HostPort,
	)
	opt := options.Client().ApplyURI(uriAddress).
		SetConnectTimeout(connectTimeout).
		SetMaxConnIdleTime(maxConnIdleTime).
		SetMinPoolSize(minPoolSize).
		SetMaxPoolSize(maxPoolSize)
	opt = opt.SetAuth(
		options.Credential{
			Username: mongoOptions.UserName,
			Password: mongoOptions.Password,
		},
	)

	mongoClient, err := mongo.Connect(ctx, opt)

	defer func() {
		if err := mongoClient.Disconnect(ctx); err != nil {
			logger.Error(
				"error in disconnecting mongodb client: %v",
				err,
			)
		}
	}()

	if err != nil {
		logError(logger, mongoOptions.Host, mongoOptions.HostPort)

		return false
	}

	err = mongoClient.Ping(ctx, nil)
	if err != nil {
		logError(logger, mongoOptions.Host, mongoOptions.HostPort)

		return false
	}
	logger.Infof(
		"Opened mongodb connection on host: %s:%d",
		mongoOptions.Host,
		mongoOptions.HostPort,
	)

	return true
}

// logError logs an error.
func logError(logger logger.Logger, host string, hostPort int) {
	// we should not use `t.Error` or `t.Errorf` for logging errors because it will `fail` our test at the end and, we just should use logs without error like log.Error (not log.Fatal)
	logger.Errorf(
		"Error in creating mongodb connection with %s:%d", host, hostPort,
	)
}
