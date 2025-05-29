// Package contracts provides a echo http server contracts.
package contracts

import echo "github.com/labstack/echo/v4"

// RouteBuilder is a struct that represents a route builder.
type RouteBuilder struct {
	echo *echo.Echo
}

// NewRouteBuilder is a function that creates a new route builder.
func NewRouteBuilder(echo *echo.Echo) *RouteBuilder {
	return &RouteBuilder{echo: echo}
}

// RegisterRoutes is a function that registers the routes.
func (r *RouteBuilder) RegisterRoutes(builder func(e *echo.Echo)) *RouteBuilder {
	builder(r.echo)

	return r
}

// RegisterGroupFunc is a function that registers the group func.
func (r *RouteBuilder) RegisterGroupFunc(
	groupName string,
	builder func(g *echo.Group),
) *RouteBuilder {
	builder(r.echo.Group(groupName))

	return r
}

// RegisterGroup is a function that registers the group.
func (r *RouteBuilder) RegisterGroup(groupName string) *RouteBuilder {
	r.echo.Group(groupName)

	return r
}

// Build is a function that builds the echo.
func (r *RouteBuilder) Build() *echo.Echo {
	return r.echo
}
