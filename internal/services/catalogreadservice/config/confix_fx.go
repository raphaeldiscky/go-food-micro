package config

import (
	"go.uber.org/fx"
)

// https://uber-go.github.io/fx/modules.html

// NewModule is a fx module for the catalog read service.
func NewModule() fx.Option {
	return fx.Module("app"+
		"configfx",
		// - order is not important in provide
		// - provide can have parameter and will resolve if registered
		// - execute its func only if it requested
		fx.Provide(
			NewConfig,
		),
	)
}
