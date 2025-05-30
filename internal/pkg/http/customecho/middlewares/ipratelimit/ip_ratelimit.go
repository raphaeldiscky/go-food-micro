// Package ipratelimit provides a echo http server ipratelimit middleware.
package ipratelimit

import (
	"log"
	"net/http"
	"strconv"

	"github.com/ulule/limiter/v3/drivers/store/memory"

	echo "github.com/labstack/echo/v4"
	limiter "github.com/ulule/limiter/v3"
)

// ref: https://github.com/ulule/limiter-examples/blob/master/echo/main.go

// IPRateLimit is a function that returns a echo middleware function.
func IPRateLimit(opts ...Option) echo.MiddlewareFunc {
	config := defualtConfig

	for _, opt := range opts {
		opt.apply(&config)
	}

	rate := limiter.Rate{
		Period: config.period,
		Limit:  config.limit,
	}

	var (
		ipRateLimiter *limiter.Limiter
		store         limiter.Store
	)

	store = memory.NewStore()
	ipRateLimiter = limiter.New(store, rate)

	// Return middleware handler.
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) (err error) {
			ip := c.RealIP()
			limiterCtx, err := ipRateLimiter.Get(c.Request().Context(), ip)
			if err != nil {
				log.Printf(
					"IPRateLimit - ipRateLimiter.Get - err: %v, %s on %s",
					err,
					ip,
					c.Request().URL,
				)

				return c.JSON(http.StatusInternalServerError, echo.Map{
					"success": false,
					"message": err,
				})
			}

			h := c.Response().Header()
			h.Set("X-RateLimit-Limit", strconv.FormatInt(limiterCtx.Limit, 10))
			h.Set(
				"X-RateLimit-Remaining",
				strconv.FormatInt(limiterCtx.Remaining, 10),
			)
			h.Set("X-RateLimit-Reset", strconv.FormatInt(limiterCtx.Reset, 10))

			if limiterCtx.Reached {
				log.Printf(
					"Too Many Requests from %s on %s",
					ip,
					c.Request().URL,
				)

				return c.JSON(http.StatusTooManyRequests, echo.Map{
					"success": false,
					"message": "Too Many Requests on " + c.Request().URL.String(),
				})
			}

			// log.Printf("%s request continue", c.RealIP()).
			return next(c)
		}
	}
}
