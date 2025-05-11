package contracts

import (
	"context"
	"testing"

	gormPostgres "github.com/raphaeldiscky/go-food-micro/internal/pkg/postgresgorm"
)

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

type GormContainer interface {
	PopulateContainerOptions(
		ctx context.Context,
		t *testing.T,
		options ...*PostgresContainerOptions,
	) (*gormPostgres.GormOptions, error)
	Cleanup(ctx context.Context) error
}
