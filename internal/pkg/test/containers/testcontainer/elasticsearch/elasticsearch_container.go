// Package elasticsearch provides an elasticsearch container.
package elasticsearch

import (
	"context"
	"fmt"
	"testing"
	"time"

	"emperror.dev/errors"
	"github.com/docker/go-connections/nat"
	"github.com/testcontainers/testcontainers-go/wait"

	elasticsearch "github.com/elastic/go-elasticsearch/v8"
	testcontainers "github.com/testcontainers/testcontainers-go"

	elasticsearchPkg "github.com/raphaeldiscky/go-food-micro/internal/pkg/elasticsearch"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/logger"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/test/containers/contracts"
)

// elasticsearchTestContainers represents an elasticsearch test containers.
type elasticsearchTestContainers struct {
	container      testcontainers.Container
	defaultOptions *contracts.ElasticsearchContainerOptions
	logger         logger.Logger
}

// NewElasticsearchTestContainers creates a new elasticsearch test containers.
func NewElasticsearchTestContainers(l logger.Logger) contracts.ElasticsearchContainer {
	return &elasticsearchTestContainers{
		defaultOptions: &contracts.ElasticsearchContainerOptions{
			Port:      "9200/tcp",
			Host:      "localhost",
			Tag:       "8.10.0",
			ImageName: "docker.elastic.co/elasticsearch/elasticsearch",
			Name:      "elasticsearch-testcontainer",
		},
		logger: l,
	}
}

// PopulateContainerOptions populates the container options.
func (e *elasticsearchTestContainers) PopulateContainerOptions(
	ctx context.Context,
	t *testing.T,
	options ...*contracts.ElasticsearchContainerOptions,
) (*elasticsearchPkg.ElasticOptions, error) {
	t.Helper()

	containerReq := e.getRunOptions(options...)

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

	// get a free random host port for elasticsearch
	hostPort, err := dbContainer.MappedPort(
		ctx,
		nat.Port(e.defaultOptions.Port),
	)
	if err != nil {
		return nil, err
	}
	e.defaultOptions.HostPort = hostPort.Int()

	host, err := dbContainer.Host(ctx)
	if err != nil {
		return nil, err
	}

	isConnectable := e.isConnectable(ctx, host, e.defaultOptions.HostPort)
	if !isConnectable {
		return e.PopulateContainerOptions(context.Background(), t, options...)
	}

	e.container = dbContainer

	option := &elasticsearchPkg.ElasticOptions{
		URL: fmt.Sprintf("http://%s:%d", host, e.defaultOptions.HostPort),
	}

	return option, nil
}

// Cleanup cleans up the container.
func (e *elasticsearchTestContainers) Cleanup(ctx context.Context) error {
	if err := e.container.Terminate(ctx); err != nil {
		return errors.WrapIf(err, "failed to terminate container: %s")
	}

	return nil
}

// getRunOptions gets the run options.
func (e *elasticsearchTestContainers) getRunOptions(
	opts ...*contracts.ElasticsearchContainerOptions,
) testcontainers.ContainerRequest {
	if len(opts) > 0 {
		e.updateOptions(opts[0])
	}

	return testcontainers.ContainerRequest{
		Image: fmt.Sprintf(
			"%s:%s",
			e.defaultOptions.ImageName,
			e.defaultOptions.Tag,
		),
		ExposedPorts: []string{e.defaultOptions.Port},
		WaitingFor: wait.ForHTTP("/").
			WithPort(nat.Port(e.defaultOptions.Port)).
			WithPollInterval(2 * time.Second).
			WithStartupTimeout(120 * time.Second),
		Hostname: e.defaultOptions.Host,
		Env: map[string]string{
			"discovery.type":         "single-node",
			"xpack.security.enabled": "false",
			"ES_JAVA_OPTS":           "-Xms512m -Xmx512m",
		},
	}
}

func (e *elasticsearchTestContainers) updateOptions(
	option *contracts.ElasticsearchContainerOptions,
) {
	if option.ImageName != "" {
		e.defaultOptions.ImageName = option.ImageName
	}
	if option.Host != "" {
		e.defaultOptions.Host = option.Host
	}
	if option.Port != "" {
		e.defaultOptions.Port = option.Port
	}
	if option.Tag != "" {
		e.defaultOptions.Tag = option.Tag
	}
}

// isConnectable checks if the container is connectable.
func (e *elasticsearchTestContainers) isConnectable(
	_ context.Context,
	host string,
	port int,
) bool {
	url := fmt.Sprintf("http://%s:%d", host, port)

	client, err := elasticsearch.NewClient(elasticsearch.Config{
		Addresses: []string{url},
	})
	if err != nil {
		e.logger.Errorf(
			"Error creating elasticsearch client with %s: %v",
			url,
			err,
		)

		return false
	}

	_, err = client.Info()
	if err != nil {
		e.logger.Errorf(
			"Error connecting to elasticsearch with %s: %v",
			url,
			err,
		)

		return false
	}

	e.logger.Infof(
		"Opened elasticsearch connection on: %s",
		url,
	)

	return true
}
