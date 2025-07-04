// Package handlers provides a echo http server handlers.
package handlers

import (
	"emperror.dev/errors"

	echo "github.com/labstack/echo/v4"

	problemDetails "github.com/raphaeldiscky/go-food-micro/internal/pkg/http/httperrors/problemdetails"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/logger"
)

// ProblemDetailErrorHandlerFunc is a function that handles the problem detail error.
func ProblemDetailErrorHandlerFunc(
	err error,
	c echo.Context,
	logger logger.Logger,
) {
	var problem problemDetails.ProblemDetailErr

	// if error was not problem detail we will convert the error to a problem detail.
	if ok := errors.As(err, &problem); !ok {
		problem = problemDetails.ParseError(err)
	}

	if !c.Response().Committed && problem != nil {
		// `WriteTo` will set `Response status code` to our problem details status
		if _, err := problemDetails.WriteTo(problem, c.Response()); err != nil {
			logger.Error(err)
		}
	}
}
