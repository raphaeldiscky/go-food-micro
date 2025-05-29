// Package rabbitmq provides a rabbitmq docker test.
package rabbitmq

import (
	"context"
	"testing"

	"github.com/raphaeldiscky/go-food-micro/internal/pkg/logger"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/rabbitmq/config"
)

// RabbitmqDockerTestContainerOptionsDecorator is a decorator that populates the rabbitmq docker test container options.
var RabbitmqDockerTestContainerOptionsDecorator = func(t *testing.T, ctx context.Context) interface{} {
	t.Helper()

	return func(c *config.RabbitmqOptions, logger logger.Logger) (*config.RabbitmqOptions, error) {
		rabbitmqHostOptions, err := NewRabbitMQDockerTest(logger).PopulateContainerOptions(ctx, t)
		c.RabbitmqHostOptions = rabbitmqHostOptions

		return c, err
	}
}
