package core

import (
	"go.uber.org/fx"

	"github.com/raphaeldiscky/go-food-micro/internal/pkg/core/serializer/json"
)

// Module provided to fxlog
// https://uber-go.github.io/fx/modules.html
var Module = fx.Module(
	"corefx",
	fx.Provide(
		json.NewDefaultJsonSerializer,
		json.NewDefaultEventJsonSerializer,
		json.NewDefaultMessageJsonSerializer,
		json.NewDefaultMetadataJsonSerializer,
	),
)
