// Package log provides a echo http server log middleware.
package log

import "github.com/labstack/echo/v4/middleware"

// config defines the config for Logger middleware.
type config struct {
	// Skipper defines a function to skip middleware.
	Skipper middleware.Skipper
}

// Option specifies instrumentation configuration options.
type Option interface {
	apply(*config)
}

// optionFunc is a function that represents a option func.
type optionFunc func(*config)

// apply is a function that applies the option.
func (o optionFunc) apply(c *config) {
	o(c)
}

// WithSkipper specifies a skipper for allowing requests to skip generating spans.
func WithSkipper(skipper middleware.Skipper) Option {
	return optionFunc(func(cfg *config) {
		cfg.Skipper = skipper
	})
}
