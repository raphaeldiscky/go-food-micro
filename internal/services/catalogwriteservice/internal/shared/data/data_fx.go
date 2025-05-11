package data

import (
	"github.com/raphaeldiscky/go-food-micro/internal/services/catalogwriteservice/internal/shared/data/dbcontext"

	"go.uber.org/fx"
)

// https://uber-go.github.io/fx/modules.html
var Module = fx.Module(
	"datafx",
	// - order is not important in provide
	// - provide can have parameter and will resolve if registered
	// - execute its func only if it requested
	fx.Provide(
		dbcontext.NewCatalogsDBContext,
	),
)
