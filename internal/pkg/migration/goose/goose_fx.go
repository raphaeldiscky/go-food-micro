package goose

import (
	"go.uber.org/fx"

	"github.com/raphaeldiscky/go-food-micro/internal/pkg/migration"
)

var (
	// Module provided to fxlog
	// https://uber-go.github.io/fx/modules.html
	Module = fx.Module( //nolint:gochecknoglobals
		"goosefx",
		mongoProviders,
	)

	mongoProviders = fx.Provide( //nolint:gochecknoglobals
		migration.ProvideConfig,
		NewGoosePostgres,
	)
)
