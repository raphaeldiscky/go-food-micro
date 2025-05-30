// Package logrous provides a logger for the application.
package logrous

import (
	"go.uber.org/fx"

	"github.com/raphaeldiscky/go-food-micro/internal/pkg/logger"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/logger/config"
)

// https://uber-go.github.io/fx/modules.html

// Module is a module for the logrous logger.
var Module = fx.Module("logrousfx",
	// - order is not important in provide
	// - provide can have parameter and will resolve if registered
	// - execute its func only if it requested
	fx.Provide(
		fx.Annotate(
			NewLogrusLogger,
			fx.As(new(logger.Logger)),
		),
		config.ProvideLogConfig,
	))

// ModuleFunc is a function that returns a module for the logrous logger.
var ModuleFunc = func(l logger.Logger) fx.Option {
	return fx.Module("logrousfx",

		fx.Provide(config.ProvideLogConfig),
		fx.Supply(fx.Annotate(l, fx.As(new(logger.Logger)))),
	)
}
