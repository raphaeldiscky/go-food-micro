// Package mongo provides the mongo docker test container options decorator.
package mongo

import (
	"context"
	"testing"

	"github.com/raphaeldiscky/go-food-micro/internal/pkg/logger"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/mongodb"
)

// MongoDockerTestContainerOptionsDecorator is a decorator for the mongo docker test container options.
var MongoDockerTestContainerOptionsDecorator = func(t *testing.T, ctx context.Context) interface{} {
	return func(_ *mongodb.MongoDbOptions, _ logger.Logger) (*mongodb.MongoDbOptions, error) {
		return NewMongoDockerTest().PopulateContainerOptions(ctx, t)
	}
}
