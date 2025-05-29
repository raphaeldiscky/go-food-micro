// Package contracts provides a redis container contracts.
package contracts

import (
	"context"
	"testing"

	redis2 "github.com/raphaeldiscky/go-food-micro/internal/pkg/redis"
)

// RedisContainerOptions represents a redis container options.
type RedisContainerOptions struct {
	Host      string
	Port      string
	HostPort  int
	Database  int
	ImageName string
	Name      string
	Tag       string
	PoolSize  int
}

// RedisContainer is a interface that represents a redis container.
type RedisContainer interface {
	PopulateContainerOptions(
		ctx context.Context,
		t *testing.T,
		options ...*RedisContainerOptions,
	) (*redis2.RedisOptions, error)
	Cleanup(ctx context.Context) error
}
