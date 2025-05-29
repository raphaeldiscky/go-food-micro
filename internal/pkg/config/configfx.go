// Package config provides a module for the config.
package config

import (
	"go.uber.org/fx"

	"github.com/raphaeldiscky/go-food-micro/internal/pkg/config/environment"
)

// Module is a module that provides a environment.Environment.
// https://uber-go.github.io/fx/modules.html
var Module = fx.Module(
	"configfx",
	fx.Provide(func() environment.Environment {
		return environment.ConfigAppEnv()
	}),
)

// ModuleFunc is a function that returns a fx.Option.
var ModuleFunc = func(e environment.Environment) fx.Option {
	return fx.Module(
		"configfx",
		fx.Provide(func() environment.Environment {
			return environment.ConfigAppEnv(e)
		}),
	)
}
