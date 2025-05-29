// Package mongo provides a mongo container.
package mongo

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"testing"

	"github.com/ory/dockertest/v3/docker"

	dockertest "github.com/ory/dockertest/v3"

	"github.com/raphaeldiscky/go-food-micro/internal/pkg/mongodb"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/test/containers/contracts"
)

// mongoDockerTest represents a mongo docker test.
type mongoDockerTest struct {
	resource       *dockertest.Resource
	defaultOptions *contracts.MongoContainerOptions
	logger         *log.Logger
}

// NewMongoDockerTest creates a new mongo docker test.
func NewMongoDockerTest() contracts.MongoContainer {
	return &mongoDockerTest{
		defaultOptions: &contracts.MongoContainerOptions{
			Database:  "test_db",
			Port:      "27017",
			Host:      "localhost",
			UserName:  "dockertest",
			Password:  "dockertest",
			Tag:       "latest",
			ImageName: "mongo",
			Name:      "mongo-dockertest",
		},
		logger: log.Default(),
	}
}

// PopulateContainerOptions populates the container options.
func (g *mongoDockerTest) PopulateContainerOptions(
	ctx context.Context,
	t *testing.T,
	options ...*contracts.MongoContainerOptions,
) (*mongodb.MongoDbOptions, error) {
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	runOptions := g.getRunOptions(options...)

	// pull mongodb docker image for version 5.0
	resource, err := pool.RunWithOptions(
		runOptions,
		func(config *docker.HostConfig) {
			// set AutoRemove to true so that stopped container goes away by itself
			config.AutoRemove = true
			config.RestartPolicy = docker.RestartPolicy{
				Name: "no",
			}
		},
	)
	if err != nil {
		log.Fatalf("Could not start resource (Mongo Container): %s", err)
	}

	if err := resource.Expire(120); err != nil {
		// Log the error but don't fail the test, as this is just a cleanup timeout
		g.logger.Printf("Error setting container expiration: %v", err)
	} // Tell docker to hard kill the container in 120 seconds exponential backoff-retry, because the application_exceptions in the container might not be ready to accept connections yet

	g.resource = resource
	port, err := strconv.Atoi(
		resource.GetPort(fmt.Sprintf("%s/tcp", g.defaultOptions.Port)),
	)
	if err != nil {
		log.Fatalf("Error converting port to int: %v", err)
	}
	g.defaultOptions.HostPort = port

	t.Cleanup(func() {
		if err := resource.Close(); err != nil {
			log.Fatalf("Error closing mongo container: %v", err)
		}
	})

	go func() {
		for range ctx.Done() {
			if err := resource.Close(); err != nil {
				log.Fatalf("Error closing mongo container: %v", err)
			}
		}
	}()

	mongoOptions := &mongodb.MongoDbOptions{
		User:     g.defaultOptions.UserName,
		Password: g.defaultOptions.Password,
		UseAuth:  false,
		Host:     g.defaultOptions.Host,
		Port:     g.defaultOptions.HostPort,
		Database: g.defaultOptions.Database,
	}

	return mongoOptions, nil
}

// Cleanup cleans up the container.
func (g *mongoDockerTest) Cleanup(_ context.Context) error {
	return g.resource.Close()
}

// getRunOptions gets the run options.
func (g *mongoDockerTest) getRunOptions(
	opts ...*contracts.MongoContainerOptions,
) *dockertest.RunOptions {
	if len(opts) > 0 && opts[0] != nil {
		option := opts[0]
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

	runOptions := &dockertest.RunOptions{
		Repository: g.defaultOptions.ImageName,
		Tag:        g.defaultOptions.Tag,
		Env: []string{
			"MONGO_INITDB_ROOT_USERNAME=" + g.defaultOptions.UserName,
			"MONGO_INITDB_ROOT_PASSWORD=" + g.defaultOptions.Password,
		},
		Hostname:     g.defaultOptions.Host,
		ExposedPorts: []string{g.defaultOptions.Port},
	}

	return runOptions
}
