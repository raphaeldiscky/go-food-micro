package config

import "go.uber.org/fx"

// https://uber-go.github.io/fx/modules.html

// NewModule creates and returns the app config module.
func NewModule() fx.Option {
	return fx.Module("appconfigfx",
		// - order is not important in provide
		// - provide can have parameter and will resolve if registered
		// - execute its func only if it requested
		fx.Provide(
			NewConfig,
		),
	)
}
