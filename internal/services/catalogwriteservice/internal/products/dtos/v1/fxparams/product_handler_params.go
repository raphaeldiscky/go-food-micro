package fxparams

import (
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/core/messaging/producer"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/logger"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/otel/tracing"
	"github.com/raphaeldiscky/go-food-micro/internal/services/catalogwriteservice/internal/shared/data/dbcontext"

	"go.uber.org/fx"
)

// ProductHandlerParams is a struct that contains the product handler params.
type ProductHandlerParams struct {
	fx.In

	Log               logger.Logger
	CatalogsDBContext *dbcontext.CatalogsGormDBContext
	RabbitmqProducer  producer.Producer
	Tracer            tracing.AppTracer
}
