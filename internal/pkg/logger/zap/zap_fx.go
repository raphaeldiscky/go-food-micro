// Package zap provides a logger for the application.
package zap

import (
	"go.uber.org/fx"

	"github.com/raphaeldiscky/go-food-micro/internal/pkg/logger"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/logger/config"
)

// https://uber-go.github.io/fx/modules.html

// Module is a module for the zap logger.
var Module = fx.Module("zapfx",

	// - order is not important in provide
	// - provide can have parameter and will resolve if registered
	// - execute its func only if it requested
	fx.Provide(
		config.ProvideLogConfig,
		NewZapLogger,
		fx.Annotate(
			NewZapLogger,
			fx.As(new(logger.Logger))),
	),
)

// ModuleFunc is a function that returns a module for the zap logger.
var ModuleFunc = func(l logger.Logger) fx.Option {
	return fx.Module(
		"zapfx",

		fx.Provide(config.ProvideLogConfig),
		fx.Supply(fx.Annotate(l, fx.As(new(logger.Logger)))),
		fx.Supply(fx.Annotate(l, fx.As(new(ZapLogger)))),
	)
}
