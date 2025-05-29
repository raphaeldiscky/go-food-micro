// Package contracts provides a mongo container contracts.
package contracts

import (
	"context"
	"testing"

	"github.com/raphaeldiscky/go-food-micro/internal/pkg/mongodb"
)

// MongoContainerOptions represents a mongo container options.
type MongoContainerOptions struct {
	Database  string
	Host      string
	Port      string
	HostPort  int
	UserName  string
	Password  string
	ImageName string
	Name      string
	Tag       string
}

// MongoContainer is a interface that represents a mongo container.
type MongoContainer interface {
	PopulateContainerOptions(
		ctx context.Context,
		t *testing.T,
		options ...*MongoContainerOptions,
	) (*mongodb.MongoDbOptions, error)
	Cleanup(ctx context.Context) error
}
