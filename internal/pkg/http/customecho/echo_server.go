// Package customecho provides a custom echo server.
package customecho

import (
	"context"
	"fmt"
	"strings"

	"github.com/labstack/echo/v4/middleware"
	"go.opentelemetry.io/otel/metric"

	echo "github.com/labstack/echo/v4"

	"github.com/raphaeldiscky/go-food-micro/internal/pkg/constants"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/http/customecho/config"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/http/customecho/contracts"
	handlers "github.com/raphaeldiscky/go-food-micro/internal/pkg/http/customecho/handlers"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/http/customecho/middlewares/ipratelimit"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/http/customecho/middlewares/log"
	otelMetrics "github.com/raphaeldiscky/go-food-micro/internal/pkg/http/customecho/middlewares/otelmetrics"
	oteltracing "github.com/raphaeldiscky/go-food-micro/internal/pkg/http/customecho/middlewares/oteltracing"
	problemdetail "github.com/raphaeldiscky/go-food-micro/internal/pkg/http/customecho/middlewares/problemdetail"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/logger"
)

// echoHTTPServer is a echo http server.
type echoHTTPServer struct {
	echo         *echo.Echo
	config       *config.EchoHTTPOptions
	log          logger.Logger
	meter        metric.Meter
	routeBuilder *contracts.RouteBuilder
}

// NewEchoHTTPServer creates a new echo http server.
func NewEchoHTTPServer(
	config *config.EchoHTTPOptions,
	logger logger.Logger,
	meter metric.Meter,
) contracts.EchoHTTPServer {
	e := echo.New()
	e.HideBanner = true

	return &echoHTTPServer{
		echo:         e,
		config:       config,
		log:          logger,
		meter:        meter,
		routeBuilder: contracts.NewRouteBuilder(e),
	}
}

// RunHTTPServer runs the http server.
func (s *echoHTTPServer) RunHTTPServer(
	configEcho ...func(echo *echo.Echo),
) error {
	s.echo.Server.ReadTimeout = constants.ReadTimeout
	s.echo.Server.WriteTimeout = constants.WriteTimeout
	s.echo.Server.MaxHeaderBytes = constants.MaxHeaderBytes

	if len(configEcho) > 0 {
		ehcoFunc := configEcho[0]
		if ehcoFunc != nil {
			configEcho[0](s.echo)
		}
	}

	// https://echo.labstack.com/guide/http_server/
	return s.echo.Start(s.config.Port)
}

// Logger returns the logger.
func (s *echoHTTPServer) Logger() logger.Logger {
	return s.log
}

// Cfg returns the config.
func (s *echoHTTPServer) Cfg() *config.EchoHTTPOptions {
	return s.config
}

// RouteBuilder returns the route builder.
func (s *echoHTTPServer) RouteBuilder() *contracts.RouteBuilder {
	return s.routeBuilder
}

// ConfigGroup configures the group.
func (s *echoHTTPServer) ConfigGroup(
	groupName string,
	groupFunc func(group *echo.Group),
) {
	groupFunc(s.echo.Group(groupName))
}

// AddMiddlewares adds the middlewares.
func (s *echoHTTPServer) AddMiddlewares(middlewares ...echo.MiddlewareFunc) {
	if len(middlewares) > 0 {
		s.echo.Use(middlewares...)
	}
}

// GracefulShutdown shuts down the http server.
func (s *echoHTTPServer) GracefulShutdown(ctx context.Context) error {
	err := s.echo.Shutdown(ctx)
	if err != nil {
		return err
	}

	return nil
}

// SetupDefaultMiddlewares sets up the default middlewares.
func (s *echoHTTPServer) SetupDefaultMiddlewares() {
	skipper := func(c echo.Context) bool {
		return strings.Contains(c.Request().URL.Path, "swagger") ||
			strings.Contains(c.Request().URL.Path, "metrics") ||
			strings.Contains(c.Request().URL.Path, "health") ||
			strings.Contains(c.Request().URL.Path, "favicon.ico")
	}

	// set error handler
	s.echo.HTTPErrorHandler = func(err error, c echo.Context) {
		// bypass skip endpoints and its error
		if skipper(c) {
			return
		}

		handlers.ProblemDetailErrorHandlerFunc(err, c, s.log)
	}

	// log errors and information
	s.echo.Use(
		log.EchoLogger(
			s.log,
			log.WithSkipper(skipper),
		),
	)
	s.echo.Use(
		oteltracing.HTTPTrace(
			oteltracing.WithSkipper(skipper),
			oteltracing.WithServiceName(s.config.Name),
		),
	)
	s.echo.Use(
		otelMetrics.HTTPMetrics(
			otelMetrics.WithServiceName(s.config.Name),
			otelMetrics.WithSkipper(skipper)),
	)
	s.echo.Use(middleware.BodyLimit(constants.BodyLimit))
	s.echo.Use(ipratelimit.IPRateLimit())
	s.echo.Use(middleware.RequestID())
	s.echo.Use(middleware.GzipWithConfig(middleware.GzipConfig{
		Level:   constants.GzipLevel,
		Skipper: skipper,
	}))
	// should be last middleware
	s.echo.Use(problemdetail.ProblemDetail(problemdetail.WithSkipper(skipper)))
}

// ApplyVersioningFromHeader applies the versioning from the header.
func (s *echoHTTPServer) ApplyVersioningFromHeader() {
	s.echo.Pre(apiVersion)
}

// GetEchoInstance returns the echo instance.
func (s *echoHTTPServer) GetEchoInstance() *echo.Echo {
	return s.echo
}

// apiVersion is the api version.
func apiVersion(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		req := c.Request()
		headers := req.Header

		apiVersion := headers.Get("version")

		req.URL.Path = fmt.Sprintf("/%s%s", apiVersion, req.URL.Path)

		return next(c)
	}
}
