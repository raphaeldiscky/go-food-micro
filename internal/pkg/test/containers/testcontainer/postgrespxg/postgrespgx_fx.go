package postgrespxg

import (
	"context"
	"testing"

	"github.com/raphaeldiscky/go-food-micro/internal/pkg/logger"
	postgres "github.com/raphaeldiscky/go-food-micro/internal/pkg/postgrespgx"
)

var PostgresPgxContainerOptionsDecorator = func(t *testing.T, ctx context.Context) interface{} {
	return func(c *postgres.PostgresPgxOptions, logger logger.Logger) (*postgres.PostgresPgxOptions, error) {
		return NewPostgresPgxContainers(logger).PopulateContainerOptions(ctx, t)
	}
}
