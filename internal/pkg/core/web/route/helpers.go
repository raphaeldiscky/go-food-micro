package route

import (
	"fmt"

	"go.uber.org/fx"
)

// AsRoute annotates the given constructor to state that
// when we register multiple handlers with output type `echo.HandlerFunc` we get exception `type already provided`, so we should use tags annotation
// it provides a route to the "routes" group.
func AsRoute(handler interface{}, routeGroupName string) interface{} {
	return fx.Annotate(
		handler,
		fx.As(new(Endpoint)),
		fx.ResultTags(fmt.Sprintf(`group:"%s"`, routeGroupName)),
	)
}
