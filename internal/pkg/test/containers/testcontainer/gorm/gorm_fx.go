// Package gorm provides a gorm fx.
package gorm

import (
	"context"
	"testing"

	"github.com/raphaeldiscky/go-food-micro/internal/pkg/logger"
	gormPostgres "github.com/raphaeldiscky/go-food-micro/internal/pkg/postgresgorm"
)

// GormContainerOptionsDecorator is a decorator for the gorm container options.
var GormContainerOptionsDecorator = func(t *testing.T, ctx context.Context) interface{} {
	return func(_ *gormPostgres.GormOptions, logger logger.Logger) (*gormPostgres.GormOptions, error) {
		return NewGormTestContainers(logger).PopulateContainerOptions(ctx, t)
	}
}
