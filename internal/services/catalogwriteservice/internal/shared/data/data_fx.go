// Package data contains the data module.
package data

import (
	"go.uber.org/fx"

	"github.com/raphaeldiscky/go-food-micro/internal/services/catalogwriteservice/internal/shared/data/dbcontext"
)

// NewModule is a module that contains the data module.
func NewModule() fx.Option {
	return fx.Module(
		"datafx",
		// - order is not important in provide
		// - provide can have parameter and will resolve if registered
		// - execute its func only if it requested
		fx.Provide(
			dbcontext.NewCatalogsDBContext,
		),
	)
}
