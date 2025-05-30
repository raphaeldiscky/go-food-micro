// Package contracts provides a gorm container contracts.
package contracts

import (
	"context"
	"testing"

	gormPostgres "github.com/raphaeldiscky/go-food-micro/internal/pkg/postgresgorm"
)

// PostgresContainerOptions represents a postgres container options.
type PostgresContainerOptions struct {
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

// GormContainer is a interface that represents a gorm container.
type GormContainer interface {
	PopulateContainerOptions(
		ctx context.Context,
		t *testing.T,
		options ...*PostgresContainerOptions,
	) (*gormPostgres.GormOptions, error)
	Cleanup(ctx context.Context) error
}
