// Package goose provides a migration runner.
package goose

import (
	"go.uber.org/fx"

	"github.com/raphaeldiscky/go-food-micro/internal/pkg/migration"
)

// Module is a module for the goose.
var (
	// Module provided to fxlog
	// https://uber-go.github.io/fx/modules.html
	Module = fx.Module(
		"goosefx",
		mongoProviders,
	)

	// mongoProviders is a module for the goose.
	mongoProviders = fx.Provide(
		migration.ProvideConfig,
		NewGoosePostgres,
	)
)
