// Package rabbitmq provides a set of functions for the rabbitmq package.
package rabbitmq

import (
	"context"

	"emperror.dev/errors"

	"github.com/raphaeldiscky/go-food-micro/internal/pkg/health/contracts"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/rabbitmq/types"
)

// gormHealthChecker is a struct that contains the rabbitmq connection.
type gormHealthChecker struct {
	connection types.IConnection
}

// NewRabbitMQHealthChecker creates a new rabbitmq health checker.
func NewRabbitMQHealthChecker(connection types.IConnection) contracts.Health {
	return &gormHealthChecker{connection}
}

// CheckHealth checks the health of the rabbitmq connection.
func (g gormHealthChecker) CheckHealth(_ context.Context) error {
	if g.connection.IsConnected() {
		return nil
	}

	return errors.New("rabbitmq is not available")
}

// GetHealthName returns the name of the rabbitmq health checker.
func (g gormHealthChecker) GetHealthName() string {
	return "rabbitmq"
}
