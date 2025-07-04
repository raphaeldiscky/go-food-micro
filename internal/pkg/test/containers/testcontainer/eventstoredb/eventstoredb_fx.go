// Package eventstoredb provides a eventstoredb fx.
package eventstoredb

import (
	"context"
	"testing"

	"github.com/raphaeldiscky/go-food-micro/internal/pkg/eventstoredb/config"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/logger"
)

// EventstoreDBContainerOptionsDecorator is a decorator for the eventstoredb container options.
var EventstoreDBContainerOptionsDecorator = func(t *testing.T, ctx context.Context) interface{} {
	t.Helper()

	return func(c *config.EventStoreDbOptions, logger logger.Logger) (*config.EventStoreDbOptions, error) {
		newOption, err := NewEventStoreDBTestContainers(logger).PopulateContainerOptions(ctx, t)
		if err != nil {
			return nil, err
		}
		newOption.Subscription = c.Subscription

		return newOption, nil
	}
}

// ReplaceEventStoreContainerOptions is a function that replaces the eventstoredb container options.
var ReplaceEventStoreContainerOptions = func(t *testing.T, options *config.EventStoreDbOptions, ctx context.Context, logger logger.Logger) error {
	t.Helper()

	newOption, err := NewEventStoreDBTestContainers(logger).PopulateContainerOptions(ctx, t)
	if err != nil {
		return err
	}

	options.HttpPort = newOption.HttpPort
	options.TcpPort = newOption.TcpPort
	options.Host = newOption.Host

	return nil
}
