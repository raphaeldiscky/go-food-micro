// Package problemdetail provides a problem detail middleware.
package problemdetail

import (
	"github.com/labstack/echo/v4/middleware"

	problemDetails "github.com/raphaeldiscky/go-food-micro/internal/pkg/http/httperrors/problemdetails"
)

// config is the config for the problem detail middleware.
type config struct {
	Skipper       middleware.Skipper
	ProblemParser problemDetails.ErrorParserFunc
}

// Option is the option for the problem detail middleware.
type Option interface {
	apply(*config)
}

// optionFunc is the option function for the problem detail middleware.
type optionFunc func(*config)

// apply applies the option to the config.
func (o optionFunc) apply(c *config) {
	o(c)
}

// WithSkipper is the option function for the problem detail middleware.
func WithSkipper(skipper middleware.Skipper) Option {
	return optionFunc(func(cfg *config) {
		cfg.Skipper = skipper
	})
}

// WithErrorParser is the option function for the problem detail middleware.
func WithErrorParser(errorParser problemDetails.ErrorParserFunc) Option {
	return optionFunc(func(cfg *config) {
		cfg.ProblemParser = errorParser
	})
}
