package gorm

import (
	"context"
	"testing"

	gormPostgres "github.com/raphaeldiscky/go-food-micro/internal/pkg/postgresgorm"
)

var GormDockerTestConatnerOptionsDecorator = func(t *testing.T, ctx context.Context) interface{} {
	return func(c *gormPostgres.GormOptions) (*gormPostgres.GormOptions, error) {
		return NewGormDockerTest().PopulateContainerOptions(ctx, t)
	}
}
