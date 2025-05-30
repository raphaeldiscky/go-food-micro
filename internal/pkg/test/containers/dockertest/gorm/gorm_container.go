// Package gorm provides a gorm container.
package gorm

import (
	"context"
	"log"
	"strconv"
	"testing"

	"github.com/ory/dockertest/v3/docker"
	"github.com/phayes/freeport"

	dockertest "github.com/ory/dockertest/v3"

	gormPostgres "github.com/raphaeldiscky/go-food-micro/internal/pkg/postgresgorm"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/test/containers/contracts"
)

// gormDockerTest represents a gorm docker test.
type gormDockerTest struct {
	resource       *dockertest.Resource
	defaultOptions *contracts.PostgresContainerOptions
	logger         *log.Logger
}

// NewGormDockerTest creates a new gorm docker test.
func NewGormDockerTest() contracts.GormContainer {
	return &gormDockerTest{
		defaultOptions: &contracts.PostgresContainerOptions{
			Database:  "test_db",
			Port:      "5432",
			Host:      "localhost",
			UserName:  "dockertest",
			Password:  "dockertest",
			Tag:       "latest",
			ImageName: "postgres",
			Name:      "postgresql-dockertest",
		},
		logger: log.Default(),
	}
}

// PopulateContainerOptions populates the container options.
func (g *gormDockerTest) PopulateContainerOptions(
	ctx context.Context,
	t *testing.T,
	options ...*contracts.PostgresContainerOptions,
) (*gormPostgres.GormOptions, error) {
	t.Helper()

	// https://github.com/ory/dockertest/blob/v3/examples/PostgreSQL.md
	// https://github.com/bozd4g/fb.testcontainers
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	runOption := g.getRunOptions(options...)
	// pulls an image, creates a container based on it and runs it
	resource, err := pool.RunWithOptions(
		runOption,
		func(config *docker.HostConfig) {
			// set AutoRemove to true so that stopped container goes away by itself
			config.AutoRemove = true
			config.RestartPolicy = docker.RestartPolicy{Name: "no"}
		})
	if err != nil {
		log.Fatalf(
			"Could not start resource (Postgresql Test Container): %s",
			err,
		)
	}

	if err := resource.Expire(120); err != nil {
		// Log the error but don't fail the test, as this is just a cleanup timeout
		g.logger.Printf("Error setting container expiration: %v", err)
	} // Tell docker to hard kill the container in 120 seconds exponential backoff-retry, because the application_exceptions in the container might not be ready to accept connections yet

	g.resource = resource
	port, err := strconv.Atoi(resource.GetPort("5432/tcp"))
	if err != nil {
		return nil, err
	}
	g.defaultOptions.HostPort = port

	t.Cleanup(func() {
		if err := resource.Close(); err != nil {
			log.Fatalf("Error closing gorm container: %v", err)
		}
	})

	var postgresoptions *gormPostgres.GormOptions

	if err = pool.Retry(func() error {
		postgresoptions = &gormPostgres.GormOptions{
			Port:     g.defaultOptions.HostPort,
			Host:     g.defaultOptions.Host,
			Password: g.defaultOptions.Password,
			DBName:   g.defaultOptions.Database,
			SSLMode:  false,
			User:     g.defaultOptions.UserName,
		}

		return nil
	}); err != nil {
		log.Fatalf("Could not connect to docker: %s", err)

		return nil, err
	}

	return postgresoptions, nil
}

// Cleanup cleans up the container.
func (g *gormDockerTest) Cleanup(ctx context.Context) error {
	return g.resource.Close()
}

// updateOptions updates the default options with provided options.
func (g *gormDockerTest) updateOptions(option *contracts.PostgresContainerOptions) {
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
func (g *gormDockerTest) getRunOptions(
	opts ...*contracts.PostgresContainerOptions,
) *dockertest.RunOptions {
	if len(opts) > 0 {
		g.updateOptions(opts[0])
	}

	hostFreePort, err := freeport.GetFreePort()
	if err != nil {
		log.Fatal(err)
	}
	g.defaultOptions.HostPort = hostFreePort

	return &dockertest.RunOptions{
		Repository: g.defaultOptions.ImageName,
		Tag:        g.defaultOptions.Tag,
		Env: []string{
			"POSTGRES_USER=" + g.defaultOptions.UserName,
			"POSTGRES_PASSWORD=" + g.defaultOptions.Password,
			"POSTGRES_DB=" + g.defaultOptions.Database,
			"listen_addresses = '*'",
		},
		Hostname:     g.defaultOptions.Host,
		ExposedPorts: []string{g.defaultOptions.Port},
		PortBindings: map[docker.Port][]docker.PortBinding{
			docker.Port(g.defaultOptions.Port): {
				{
					HostIP:   "0.0.0.0",
					HostPort: strconv.Itoa(g.defaultOptions.HostPort),
				},
			},
		},
	}
}
