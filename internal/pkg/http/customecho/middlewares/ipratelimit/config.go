// Package ipratelimit provides a echo http server ipratelimit middleware.
package ipratelimit

import (
	"time"
)

// config is a struct that represents a config.
type config struct {
	period time.Duration
	limit  int64
}

// defualtConfig is a struct that represents a default config.
var defualtConfig = config{
	period: 1 * time.Hour,
	limit:  1000,
}

// Option is a function that applies a config.
type Option interface {
	apply(*config)
}

// optionFunc is a function that represents a option func.
type optionFunc func(*config)

// apply is a function that applies the option.
func (o optionFunc) apply(c *config) {
	o(c)
}

// WithPeriod is a function that sets the period.
func WithPeriod(d time.Duration) Option {
	return optionFunc(func(cfg *config) {
		if cfg.period != 0 {
			cfg.period = d
		}
	})
}

// WithLimit is a function that sets the limit.
func WithLimit(v int64) Option {
	return optionFunc(func(cfg *config) {
		if cfg.limit != 0 {
			cfg.limit = v
		}
	})
}
