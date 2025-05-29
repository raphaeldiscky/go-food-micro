// Package rabbitmq provides a rabbitmq container.
package rabbitmq

import (
	"context"
	"fmt"
	"testing"
	"time"

	"emperror.dev/errors"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/go-connections/nat"
	"github.com/testcontainers/testcontainers-go/wait"

	rabbithole "github.com/michaelklishin/rabbit-hole"
	amqp091 "github.com/rabbitmq/amqp091-go"
	testcontainers "github.com/testcontainers/testcontainers-go"

	"github.com/raphaeldiscky/go-food-micro/internal/pkg/logger"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/rabbitmq/config"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/test/containers/contracts"
)

// rabbitmqTestContainers represents a rabbitmq test containers.
// https://github.com/testcontainers/testcontainers-go/issues/1359
// https://github.com/testcontainers/testcontainers-go/issues/1249
type rabbitmqTestContainers struct {
	container      testcontainers.Container
	defaultOptions *contracts.RabbitMQContainerOptions
	logger         logger.Logger
}

// NewRabbitMQTestContainers creates a new rabbitmq test containers.
func NewRabbitMQTestContainers(l logger.Logger) contracts.RabbitMQContainer {
	return &rabbitmqTestContainers{
		defaultOptions: &contracts.RabbitMQContainerOptions{
			Ports:       []string{"5672/tcp", "15672/tcp"},
			Host:        "localhost",
			VirtualHost: "/",
			UserName:    "guest",
			Password:    "guest",
			HttpPort:    15672,
			HostPort:    5672,
			Tag:         "management",
			ImageName:   "rabbitmq",
			Name:        "rabbitmq-testcontainers",
		},
		logger: l,
	}
}

// PopulateContainerOptions populates the container options.
func (g *rabbitmqTestContainers) PopulateContainerOptions(
	ctx context.Context,
	t *testing.T,
	options ...*contracts.RabbitMQContainerOptions,
) (*config.RabbitmqHostOptions, error) {
	maxRetries := 5
	for i := 0; i < maxRetries; i++ {
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
			if terr := dbContainer.Terminate(ctx); terr != nil {
				t.Fatalf("failed to terminate container: %s", err)
			}
			time.Sleep(time.Second * 1)
		})

		// get a free random host port for rabbitmq `Tcp Port`
		hostPort, err := dbContainer.MappedPort(
			ctx,
			nat.Port(g.defaultOptions.Ports[0]),
		)
		if err != nil {
			return nil, err
		}
		g.defaultOptions.HostPort = hostPort.Int()
		g.logger.Infof("rabbitmq host port is: %d", hostPort.Int())

		// get a free random host port for rabbitmq UI `Http Port`
		uiHTTPPort, err := dbContainer.MappedPort(
			ctx,
			nat.Port(g.defaultOptions.Ports[1]),
		)
		if err != nil {
			return nil, err
		}
		g.defaultOptions.HttpPort = uiHTTPPort.Int()
		g.logger.Infof("rabbitmq ui port is: %d", uiHTTPPort.Int())

		host, err := dbContainer.Host(ctx)
		if err != nil {
			return nil, err
		}

		// Try to connect with a timeout
		connectCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
		isConnectable := IsConnectableWithContext(connectCtx, g.logger, g.defaultOptions)
		cancel()

		if isConnectable {
			g.container = dbContainer
			option := &config.RabbitmqHostOptions{
				UserName:    g.defaultOptions.UserName,
				Password:    g.defaultOptions.Password,
				HostName:    host,
				VirtualHost: g.defaultOptions.VirtualHost,
				Port:        g.defaultOptions.HostPort,
				HttpPort:    g.defaultOptions.HttpPort,
			}

			return option, nil
		}

		// Clean up the failed container
		if err := dbContainer.Terminate(ctx); err != nil {
			g.logger.Errorf("failed to terminate container: %v", err)
		}

		if i < maxRetries-1 {
			g.logger.Infof("Retrying RabbitMQ connection (attempt %d/%d)...", i+1, maxRetries)
			time.Sleep(time.Second * 2)
		}
	}

	return nil, errors.New("failed to connect to RabbitMQ after maximum retries")
}

// Cleanup cleans up the container.
func (g *rabbitmqTestContainers) Cleanup(ctx context.Context) error {
	if err := g.container.Terminate(ctx); err != nil {
		return errors.WrapIf(err, "failed to terminate container: %s")
	}

	return nil
}

// getRunOptions gets the run options.
func (g *rabbitmqTestContainers) getRunOptions(
	opts ...*contracts.RabbitMQContainerOptions,
) testcontainers.ContainerRequest {
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
	containerReq := testcontainers.ContainerRequest{
		Image: fmt.Sprintf(
			"%s:%s",
			g.defaultOptions.ImageName,
			g.defaultOptions.Tag,
		),
		ExposedPorts: g.defaultOptions.Ports,
		WaitingFor: wait.ForAll(
			wait.ForListeningPort(nat.Port(g.defaultOptions.Ports[0])),
			wait.ForLog("Server startup complete"),
			wait.ForLog("started TCP listener"),
		),
		HostConfigModifier: func(hostConfig *container.HostConfig) {
			hostConfig.AutoRemove = true
		},
		Hostname: g.defaultOptions.Host,
		Env: map[string]string{
			"RABBITMQ_DEFAULT_USER": g.defaultOptions.UserName,
			"RABBITMQ_DEFAULT_PASS": g.defaultOptions.Password,
		},
	}

	return containerReq
}

// IsConnectableWithContext checks if the container is connectable.
func IsConnectableWithContext(
	ctx context.Context,
	logger logger.Logger,
	options *contracts.RabbitMQContainerOptions,
) bool {
	// Create a channel to handle the connection attempt with timeout
	connChan := make(chan bool, 1)

	go func() {
		conn, err := amqp091.Dial(options.AmqpEndPoint())
		if err != nil {
			logError(
				logger,
				options.UserName,
				options.Password,
				options.Host,
				options.HostPort,
			)
			connChan <- false

			return
		}
		defer conn.Close()

		if conn.IsClosed() {
			logError(
				logger,
				options.UserName,
				options.Password,
				options.Host,
				options.HostPort,
			)
			connChan <- false

			return
		}

		// Test HTTP connection
		rmqc, err := rabbithole.NewClient(
			options.HTTPEndPoint(),
			options.UserName,
			options.Password,
		)
		if err != nil {
			logger.Errorf(
				"Error creating RabbitMQ HTTP client: %v",
				err,
			)
			connChan <- false

			return
		}

		_, err = rmqc.ListExchanges()
		if err != nil {
			logger.Errorf(
				"Error in creating rabbitmq connection with http host: %s",
				options.HTTPEndPoint(),
			)
			connChan <- false

			return
		}

		logger.Infof(
			"Successfully connected to RabbitMQ on host: amqp://%s:%s@%s:%d",
			options.UserName,
			options.Password,
			options.Host,
			options.HostPort,
		)
		connChan <- true
	}()

	select {
	case result := <-connChan:
		return result
	case <-ctx.Done():
		logger.Errorf("Connection attempt timed out after %v", 10*time.Second)

		return false
	}
}

// logError logs an error.
func logError(
	logger logger.Logger,
	userName string,
	password string,
	host string,
	hostPort int,
) {
	// we should not use `t.Error` or `t.Errorf` for logging errors because it will `fail` our test at the end and, we just should use logs without error like log.Error (not log.Fatal)
	logger.Errorf(
		"Error in creating rabbitmq connection with amqp host: amqp://%s:%s@%s:%d",
		userName,
		password,
		host,
		hostPort,
	)
}
