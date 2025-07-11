// Package fxapp provides a module for the fxapp.
package fxapp

import (
	"go.uber.org/fx"

	"github.com/raphaeldiscky/go-food-micro/internal/pkg/config/environment"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/fxapp/contracts"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/logger"
	loggerConfig "github.com/raphaeldiscky/go-food-micro/internal/pkg/logger/config"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/logger/logrous"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/logger/models"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/logger/zap"
)

// applicationBuilder is a builder for the application.
type applicationBuilder struct {
	provides    []interface{}
	decorates   []interface{}
	options     []fx.Option
	logger      logger.Logger
	environment environment.Environment
}

// NewApplicationBuilder creates a new application builder.
func NewApplicationBuilder(environments ...environment.Environment) contracts.ApplicationBuilder {
	env := environment.ConfigAppEnv(environments...)

	var logger logger.Logger
	logoption, err := loggerConfig.ProvideLogConfig(env)
	if err != nil || logoption == nil {
		logger = zap.NewZapLogger(logoption, env)
	} else if logoption.LogType == models.Logrus {
		logger = logrous.NewLogrusLogger(logoption, env)
	} else {
		logger = zap.NewZapLogger(logoption, env)
	}

	return &applicationBuilder{logger: logger, environment: env}
}

// ProvideModule provides a module.
func (a *applicationBuilder) ProvideModule(module fx.Option) {
	a.options = append(a.options, module)
}

// Provide provides a constructor.
func (a *applicationBuilder) Provide(constructors ...interface{}) {
	a.provides = append(a.provides, constructors...)
}

// Decorate decorates a constructor.
func (a *applicationBuilder) Decorate(constructors ...interface{}) {
	a.decorates = append(a.decorates, constructors...)
}

// Build builds the application.
func (a *applicationBuilder) Build() contracts.Application {
	app := NewApplication(a.provides, a.decorates, a.options, a.logger, a.environment)

	return app
}

// GetProvides gets the provides.
func (a *applicationBuilder) GetProvides() []interface{} {
	return a.provides
}

// GetDecorates gets the decorates.
func (a *applicationBuilder) GetDecorates() []interface{} {
	return a.decorates
}

// Options gets the options.
func (a *applicationBuilder) Options() []fx.Option {
	return a.options
}

// Logger gets the logger.
func (a *applicationBuilder) Logger() logger.Logger {
	return a.logger
}

// Environment gets the environment.
func (a *applicationBuilder) Environment() environment.Environment {
	return a.environment
}
