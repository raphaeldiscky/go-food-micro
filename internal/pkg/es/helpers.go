// Package es provides a helpers for the event sourcing.
package es

import (
	"go.uber.org/fx"

	"github.com/raphaeldiscky/go-food-micro/internal/pkg/es/contracts/projection"
)

// AsProjection is a helper function that annotates a handler with the projection.IProjection interface.
func AsProjection(handler interface{}) interface{} {
	return fx.Annotate(
		handler,
		fx.As(new(projection.IProjection)),
		fx.ResultTags(`group:"projections"`),
	)
}
