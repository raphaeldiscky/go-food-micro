// Package client provides a http client module.
package client

import "go.uber.org/fx"

// Module is a fx.Module that provides the http client module.
// https://uber-go.github.io/fx/modules.html
var Module = fx.Module("clientfx",
	// Module provided to fxlog
	// - order is not important in provide
	// - provide can have parameter and will resolve if registered
	// - execute its func only if it requested
	fx.Provide(NewHTTPClient),
)
