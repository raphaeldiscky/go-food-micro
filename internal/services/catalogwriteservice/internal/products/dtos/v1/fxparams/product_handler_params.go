// Package fxparams contains the product handler params.
package fxparams

import (
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/core/messaging/producer"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/logger"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/otel/tracing"
	"go.uber.org/fx"

	"github.com/raphaeldiscky/go-food-micro/internal/services/catalogwriteservice/internal/products/contracts"
	"github.com/raphaeldiscky/go-food-micro/internal/services/catalogwriteservice/internal/shared/data/dbcontext"
)

// ProductHandlerParams is a struct that contains the product handler params.
type ProductHandlerParams struct {
	fx.In

	Log               logger.Logger
	CatalogsDBContext *dbcontext.CatalogsGormDBContext
	RabbitmqProducer  producer.Producer
	Tracer            tracing.AppTracer
	ProductRepository contracts.ProductRepository
}
