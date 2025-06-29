// Package contracts provides an elasticsearch container contracts.
package contracts

import (
	"context"
	"testing"

	elasticsearch "github.com/raphaeldiscky/go-food-micro/internal/pkg/elasticsearch"
)

// ElasticsearchContainerOptions represents an elasticsearch container options.
type ElasticsearchContainerOptions struct {
	Host      string
	Port      string
	HostPort  int
	ImageName string
	Name      string
	Tag       string
}

// ElasticsearchContainer is an interface that represents an elasticsearch container.
type ElasticsearchContainer interface {
	PopulateContainerOptions(
		ctx context.Context,
		t *testing.T,
		options ...*ElasticsearchContainerOptions,
	) (*elasticsearch.ElasticOptions, error)
	Cleanup(ctx context.Context) error
}
