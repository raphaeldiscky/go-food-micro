// Package gomigrate provides a migration runner.
package gomigrate

import (
	"go.uber.org/fx"

	"github.com/raphaeldiscky/go-food-micro/internal/pkg/migration"
)

// Module is a module for the gomigrate.
var (
	// Module provided to fxlog
	// https://uber-go.github.io/fx/modules.html
	Module = fx.Module( //nolint:gochecknoglobals
		"gomigratefx",
		mongoProviders,
	)

	// mongoProviders is a module for the gomigrate.
	mongoProviders = fx.Provide( //nolint:gochecknoglobals
		migration.ProvideConfig,
		NewGoMigratorPostgres,
	)
)
