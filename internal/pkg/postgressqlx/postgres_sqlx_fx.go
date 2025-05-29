// Package postgressqlx provides a set of functions for the postgressqlx.
package postgressqlx

import (
	"context"

	"go.uber.org/fx"

	"github.com/raphaeldiscky/go-food-micro/internal/pkg/logger"
)

// Module is the module for the postgressqlx.
// https://uber-go.github.io/fx/modules.html
var Module = fx.Module("postgressqlxfx",
	fx.Provide(NewSqlxConn, provideConfig),
	fx.Invoke(registerHooks),
)

// registerHooks registers the hooks for the postgressqlx.
func registerHooks(lc fx.Lifecycle, pgxClient *Sqlx, logger logger.Logger) {
	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			pgxClient.Close()
			logger.Info("Sqlx postgres connection closed gracefully")

			return nil
		},
	})
}
