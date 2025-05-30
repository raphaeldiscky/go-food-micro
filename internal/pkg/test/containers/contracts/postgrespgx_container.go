// Package contracts provides a postgrespgx container contracts.
package contracts

import (
	"context"
	"testing"

	postgres "github.com/raphaeldiscky/go-food-micro/internal/pkg/postgrespgx"
)

// PostgresPgxContainer is a interface that represents a postgrespgx container.
type PostgresPgxContainer interface {
	PopulateContainerOptions(
		ctx context.Context,
		t *testing.T,
		options ...*PostgresContainerOptions,
	) (*postgres.PostgresPgxOptions, error)
	Cleanup(ctx context.Context) error
}
