// Package mongo provides a mongo fx.
package mongo

import (
	"context"
	"testing"

	"github.com/raphaeldiscky/go-food-micro/internal/pkg/logger"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/mongodb"
)

// MongoContainerOptionsDecorator is a decorator for the mongo container options.
var MongoContainerOptionsDecorator = func(t *testing.T, ctx context.Context) interface{} {
	t.Helper()

	return func(c *mongodb.MongoDbOptions, logger logger.Logger) (*mongodb.MongoDbOptions, error) {
		return NewMongoTestContainers(logger).PopulateContainerOptions(ctx, t)
	}
}
