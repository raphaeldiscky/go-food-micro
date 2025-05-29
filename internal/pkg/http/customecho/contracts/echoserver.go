// Package contracts provides a echo http server contracts.
package contracts

import (
	"context"

	echo "github.com/labstack/echo/v4"

	"github.com/raphaeldiscky/go-food-micro/internal/pkg/http/customecho/config"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/logger"
)

// EchoHttpServer is an interface that represents a echo http server.
type EchoHttpServer interface {
	RunHttpServer(configEcho ...func(echo *echo.Echo)) error
	GracefulShutdown(ctx context.Context) error
	ApplyVersioningFromHeader()
	GetEchoInstance() *echo.Echo
	Logger() logger.Logger
	Cfg() *config.EchoHTTPOptions
	SetupDefaultMiddlewares()
	RouteBuilder() *RouteBuilder
	AddMiddlewares(middlewares ...echo.MiddlewareFunc)
	ConfigGroup(groupName string, groupFunc func(group *echo.Group))
}
