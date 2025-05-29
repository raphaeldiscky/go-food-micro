// Package redis provides a redis fx.
package redis

import (
	"context"
	"testing"

	"github.com/raphaeldiscky/go-food-micro/internal/pkg/logger"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/redis"
)

// RedisContainerOptionsDecorator is a decorator for the redis container options.
var RedisContainerOptionsDecorator = func(t *testing.T, ctx context.Context) interface{} {
	return func(c *redis.RedisOptions, logger logger.Logger) (*redis.RedisOptions, error) {
		return NewRedisTestContainers(logger).PopulateContainerOptions(ctx, t)
	}
}
