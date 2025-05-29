// Package postgrespxg provides a postgrespgx fx.
package postgrespxg

import (
	"context"
	"testing"

	"github.com/raphaeldiscky/go-food-micro/internal/pkg/logger"
	postgres "github.com/raphaeldiscky/go-food-micro/internal/pkg/postgrespgx"
)

// PostgresPgxContainerOptionsDecorator is a decorator for the postgrespgx container options.
var PostgresPgxContainerOptionsDecorator = func(t *testing.T, ctx context.Context) interface{} {
	t.Helper()

	return func(_ *postgres.PostgresPgxOptions, logger logger.Logger) (*postgres.PostgresPgxOptions, error) {
		return NewPostgresPgxContainers(logger).PopulateContainerOptions(ctx, t)
	}
}
