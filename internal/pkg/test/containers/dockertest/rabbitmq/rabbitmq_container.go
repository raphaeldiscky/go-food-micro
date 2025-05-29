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

type rabbitmqDockerTest struct {
	resource       *dockertest.Resource
	defaultOptions *contracts.RabbitMQContainerOptions
	logger         logger.Logger
}

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
			Name:        "rabbitmq-dockertest",
		},
		logger: logger,
	}
}

func (g *rabbitmqDockerTest) PopulateContainerOptions(
	ctx context.Context,
	t *testing.T,
	options ...*contracts.RabbitMQContainerOptions,
) (*config.RabbitmqHostOptions, error) {
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
		log.Fatalf("Could not start resource (RabbitMQ Container): %s", err)
	}

	resource.Expire(
		120,
	) // Tell docker to hard kill the container in 120 seconds exponential backoff-retry, because the application_exceptions in the container might not be ready to accept connections yet

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

	t.Cleanup(func() { _ = resource.Close() })

	// isConnectable := isConnectable(g.logger, g.defaultOptions)
	// if !isConnectable {
	//	return g.PopulateContainerOptions(context.Background(), t, options...)
	//}

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

func (g *rabbitmqDockerTest) Cleanup(ctx context.Context) error {
	return nil
}

func (g *rabbitmqDockerTest) getRunOptions(
	opts ...*contracts.RabbitMQContainerOptions,
) *dockertest.RunOptions {
	if len(opts) > 0 && opts[0] != nil {
		option := opts[0]
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

	runOptions := &dockertest.RunOptions{
		Repository: g.defaultOptions.ImageName,
		Tag:        g.defaultOptions.Tag,
		Env: []string{
			"RABBITMQ_DEFAULT_USER=" + g.defaultOptions.UserName,
			"RABBITMQ_DEFAULT_PASS=" + g.defaultOptions.Password,
		},
		Hostname:     g.defaultOptions.Host,
		ExposedPorts: g.defaultOptions.Ports,
	}

	return runOptions
}
