package contracts

import (
	"context"

	"github.com/raphaeldiscky/go-food-micro/internal/pkg/http/customecho/config"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/logger"

	"github.com/labstack/echo/v4"
)

type EchoHttpServer interface {
	RunHttpServer(configEcho ...func(echo *echo.Echo)) error
	GracefulShutdown(ctx context.Context) error
	ApplyVersioningFromHeader()
	GetEchoInstance() *echo.Echo
	Logger() logger.Logger
	Cfg() *config.EchoHttpOptions
	SetupDefaultMiddlewares()
	RouteBuilder() *RouteBuilder
	AddMiddlewares(middlewares ...echo.MiddlewareFunc)
	ConfigGroup(groupName string, groupFunc func(group *echo.Group))
}
