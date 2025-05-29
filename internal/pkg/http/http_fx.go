package http

import (
	"go.uber.org/fx"

	"github.com/raphaeldiscky/go-food-micro/internal/pkg/http/client"
	customEcho "github.com/raphaeldiscky/go-food-micro/internal/pkg/http/customecho"
)

// Module provided to fxlog
// https://uber-go.github.io/fx/modules.html
var Module = fx.Module("httpfx",
	// - order is not important in provide
	// - provide can have parameter and will resolve if registered
	// - execute its func only if it requested
	client.Module,
	customEcho.Module,
)
