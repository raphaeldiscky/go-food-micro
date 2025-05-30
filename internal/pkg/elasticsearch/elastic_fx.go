// Package elasticsearch provides the elasticsearch module.
package elasticsearch

import (
	"go.uber.org/fx"
)

// Module is the module for the elasticsearch.
// https://uber-go.github.io/fx/modules.html
var Module = fx.Module("elasticfx",
	fx.Provide(provideConfig),
	fx.Provide(NewElasticClient),
)
