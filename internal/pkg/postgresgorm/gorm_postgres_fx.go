// Package postgresgorm provides a set of functions for the postgres gorm.
package postgresgorm

import (
	"fmt"

	"go.uber.org/fx"

	"github.com/raphaeldiscky/go-food-micro/internal/pkg/health/contracts"
)

// Module is the module for the postgres gorm.
// https://uber-go.github.io/fx/modules.html
var Module = fx.Module(
	"gormpostgresfx",
	fx.Provide(
		provideConfig,
		NewGorm,
		NewSQLDB,

		fx.Annotate(
			NewGormHealthChecker,
			fx.As(new(contracts.Health)),
			fx.ResultTags(fmt.Sprintf(`group:"%s"`, "healths")),
		),
	),
)
