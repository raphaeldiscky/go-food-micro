// Package problemdetail provides a problem detail middleware.
package problemdetail

import (
	"github.com/labstack/echo/v4/middleware"

	echo "github.com/labstack/echo/v4"

	problemDetails "github.com/raphaeldiscky/go-food-micro/internal/pkg/http/httperrors/problemdetails"
)

// ProblemDetail is the problem detail middleware.
func ProblemDetail(opts ...Option) echo.MiddlewareFunc {
	cfg := config{}
	for _, opt := range opts {
		opt.apply(&cfg)
	}

	if cfg.Skipper == nil {
		cfg.Skipper = middleware.DefaultSkipper
	}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if cfg.Skipper(c) {
				return next(c)
			}

			err := next(c)

			prbError := problemDetails.ParseError(err)

			if cfg.ProblemParser != nil {
				prbError = cfg.ProblemParser(prbError)
			}

			if prbError != nil {
				// handle echo error in this middleware and raise echo errorhandler func and our custom error handler
				// when we call c.Error more than once, `c.Response().Committed` becomes true and response doesn't write to client again in our error handler
				// Error will update response status with occurred error object status code
				c.Error(prbError)
			}

			return prbError
		}
	}
}
