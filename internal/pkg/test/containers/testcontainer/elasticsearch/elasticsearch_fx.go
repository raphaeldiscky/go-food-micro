// Package elasticsearch provides an elasticsearch fx.
package elasticsearch

import (
	"context"
	"testing"

	"github.com/raphaeldiscky/go-food-micro/internal/pkg/elasticsearch"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/logger"
)

// ElasticsearchContainerOptionsDecorator is a decorator for the elasticsearch container options.
var ElasticsearchContainerOptionsDecorator = func(t *testing.T, ctx context.Context) interface{} {
	t.Helper()

	return func(c *elasticsearch.ElasticOptions, logger logger.Logger) (*elasticsearch.ElasticOptions, error) {
		return NewElasticsearchTestContainers(logger).PopulateContainerOptions(ctx, t)
	}
}
