// Package contracts provides a module for the contracts.
package contracts

import (
	"go.uber.org/fx"

	"github.com/raphaeldiscky/go-food-micro/internal/pkg/config/environment"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/logger"
)

// ApplicationBuilder is a builder for the application.
type ApplicationBuilder interface {
	// ProvideModule register modules directly instead and modules should not register with `provided` function
	ProvideModule(module fx.Option)
	// Provide register functions constructors as dependency resolver
	Provide(constructors ...interface{})
	// Decorate decorates the application
	Decorate(constructors ...interface{})
	Build() Application

	GetProvides() []interface{}
	GetDecorates() []interface{}
	Options() []fx.Option
	Logger() logger.Logger
	Environment() environment.Environment
}
