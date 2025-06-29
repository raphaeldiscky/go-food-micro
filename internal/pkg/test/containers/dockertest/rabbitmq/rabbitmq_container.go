// Package rabbitmq provides a rabbitmq docker test.
package rabbitmq

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"testing"

	"github.com/ory/dockertest/v3/docker"

	dockertest "github.com/ory/dockertest/v3"

	"github.com/raphaeldiscky/go-food-micro/internal/pkg/logger"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/rabbitmq/config"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/test/containers/contracts"
)

// rabbitmqDockerTest is a struct that represents a rabbitmq docker test.
type rabbitmqDockerTest struct {
	resource       *dockertest.Resource
	defaultOptions *contracts.RabbitMQContainerOptions
	logger         logger.Logger
}

// NewRabbitMQDockerTest creates a new rabbitmq docker test.
func NewRabbitMQDockerTest(logger logger.Logger) contracts.RabbitMQContainer {
	return &rabbitmqDockerTest{
		defaultOptions: &contracts.RabbitMQContainerOptions{
			Ports:       []string{"5672", "15672"},
			Host:        "localhost",
			VirtualHost: "",
			UserName:    "dockertest",
			Password:    "dockertest",
			Tag:         "management",
			ImageName:   "rabbitmq",
			Name:        "", // Will be set dynamically to avoid conflicts
		},
		logger: logger,
	}
}

// PopulateContainerOptions populates the container options.
func (g *rabbitmqDockerTest) PopulateContainerOptions(
	ctx context.Context,
	t *testing.T,
	options ...*contracts.RabbitMQContainerOptions,
) (*config.RabbitmqHostOptions, error) {
	t.Helper()

	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	// Generate unique container name to avoid conflicts
	uniqueName := fmt.Sprintf("rabbitmq-dockertest-%s", t.Name())
	g.defaultOptions.Name = uniqueName

	// Clean up any existing containers with the same name
	if err := g.cleanupExistingContainer(pool, uniqueName); err != nil {
		g.logger.Errorf("Error cleaning up existing container: %v", err)
	}

	runOptions := g.getRunOptions(options...)

	// pull rabbitmq docker image
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
		log.Fatalf("Could not start resource (RabbitMQ Container): %s", err)
	}

	if err := resource.Expire(120); err != nil {
		// Log the error but don't fail the test, as this is just a cleanup timeout
		g.logger.Errorf("Error setting container expiration: %v", err)
	}

	g.resource = resource
	hostPort, err := strconv.Atoi(
		resource.GetPort(fmt.Sprintf("%s/tcp", g.defaultOptions.Ports[0])),
	) // 5672
	if err != nil {
		return nil, err
	}
	httpPort, err := strconv.Atoi(
		resource.GetPort(fmt.Sprintf("%s/tcp", g.defaultOptions.Ports[1])),
	) // 15672
	if err != nil {
		return nil, err
	}

	g.defaultOptions.HostPort = hostPort
	g.defaultOptions.HttpPort = httpPort

	t.Cleanup(func() {
		if err := resource.Close(); err != nil {
			g.logger.Errorf("Error closing rabbitmq container: %v", err)
		}
	})

	var rabbitmqoptions *config.RabbitmqHostOptions
	if err = pool.Retry(func() error {
		rabbitmqoptions = &config.RabbitmqHostOptions{
			UserName:    g.defaultOptions.UserName,
			Password:    g.defaultOptions.Password,
			HostName:    g.defaultOptions.Host,
			VirtualHost: g.defaultOptions.VirtualHost,
			Port:        g.defaultOptions.HostPort,
			HttpPort:    g.defaultOptions.HttpPort,
		}

		return nil
	}); err != nil {
		log.Fatalf("Could not connect to docker: %s", err)

		return nil, err
	}

	return rabbitmqoptions, nil
}

// Cleanup cleans up the rabbitmq docker test.
func (g *rabbitmqDockerTest) Cleanup(ctx context.Context) error {
	return nil
}

// getRunOptions gets the run options.
func (g *rabbitmqDockerTest) getRunOptions(
	opts ...*contracts.RabbitMQContainerOptions,
) *dockertest.RunOptions {
	if len(opts) > 0 {
		g.updateOptions(opts[0])
	}

	return &dockertest.RunOptions{
		Repository: g.defaultOptions.ImageName,
		Tag:        g.defaultOptions.Tag,
		Name:       g.defaultOptions.Name,
		Env: []string{
			"RABBITMQ_DEFAULT_USER=" + g.defaultOptions.UserName,
			"RABBITMQ_DEFAULT_PASS=" + g.defaultOptions.Password,
		},
		Hostname:     g.defaultOptions.Host,
		ExposedPorts: []string{"5672", "15672"},
	}
}

func (g *rabbitmqDockerTest) updateOptions(option *contracts.RabbitMQContainerOptions) {
	if option.ImageName != "" {
		g.defaultOptions.ImageName = option.ImageName
	}
	if option.Host != "" {
		g.defaultOptions.Host = option.Host
	}
	if len(option.Ports) > 0 {
		g.defaultOptions.Ports = option.Ports
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

func (g *rabbitmqDockerTest) cleanupExistingContainer(
	pool *dockertest.Pool,
	containerName string,
) error {
	// Try to remove any existing container with the same name
	// This is a best-effort cleanup, we don't fail if container doesn't exist
	if err := pool.RemoveContainerByName(containerName); err != nil {
		// Container doesn't exist or already removed, which is fine
		g.logger.Debugf(
			"Container %s not found for cleanup (this is normal): %v",
			containerName,
			err,
		)
	}

	return nil
}
