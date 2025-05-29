// Package rabbitmq provides a rabbitmq fx.
package rabbitmq

import (
	"context"
	"testing"

	"github.com/raphaeldiscky/go-food-micro/internal/pkg/logger"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/rabbitmq/config"
)

// RabbitmqContainerOptionsDecorator is a decorator for the rabbitmq container options.
var RabbitmqContainerOptionsDecorator = func(t *testing.T, ctx context.Context) interface{} {
	return func(c *config.RabbitmqOptions, logger logger.Logger) (*config.RabbitmqOptions, error) {
		rabbitmqHostOptions, err := NewRabbitMQTestContainers(
			logger,
		).PopulateContainerOptions(ctx, t)
		c.RabbitmqHostOptions = rabbitmqHostOptions

		return c, err
	}
}
