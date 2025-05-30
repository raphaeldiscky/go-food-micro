// Package pipelines provides a config for the pipelines.
package pipelines

import (
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/logger"
	defaultLogger "github.com/raphaeldiscky/go-food-micro/internal/pkg/logger/defaultlogger"
)

// config is a config for the pipelines.
type config struct {
	logger      logger.Logger
	serviceName string
}

// defaultConfig is a default config for the pipelines.
var defaultConfig = &config{
	serviceName: "app",
	logger:      defaultLogger.GetLogger(),
}

// Option specifies instrumentation configuration options.
type Option interface {
	apply(*config)
}

// optionFunc is a function that applies an option to a config.
type optionFunc func(*config)

// apply applies an option to a config.
func (o optionFunc) apply(c *config) {
	o(c)
}

// WithServiceName sets the service name for the pipelines.
func WithServiceName(v string) Option {
	return optionFunc(func(cfg *config) {
		if cfg.serviceName != "" {
			cfg.serviceName = v
		}
	})
}

// WithLogger sets the logger for the pipelines.
func WithLogger(l logger.Logger) Option {
	return optionFunc(func(cfg *config) {
		if cfg.logger != nil {
			cfg.logger = l
		}
	})
}
