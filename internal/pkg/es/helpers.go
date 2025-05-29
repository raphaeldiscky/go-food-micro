package es

import (
	"go.uber.org/fx"

	"github.com/raphaeldiscky/go-food-micro/internal/pkg/es/contracts/projection"
)

func AsProjection(handler interface{}) interface{} {
	return fx.Annotate(
		handler,
		fx.As(new(projection.IProjection)),
		fx.ResultTags(`group:"projections"`),
	)
}
