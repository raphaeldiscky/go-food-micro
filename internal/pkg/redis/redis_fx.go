// Package redis provides a set of functions for the redis package.
package redis

import (
	"context"
	"fmt"

	"go.uber.org/fx"

	redis "github.com/redis/go-redis/v9"

	"github.com/raphaeldiscky/go-food-micro/internal/pkg/health/contracts"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/logger"
)

// Module provided to fxlog.
var (
	// Module provided to fxlog
	// https://uber-go.github.io/fx/modules.html
	Module = fx.Module(
		"redisfx",
		redisProviders,
		redisInvokes,
	)

	redisProviders = fx.Options(fx.Provide(
		NewRedisClient,
		func(client *redis.Client) redis.UniversalClient {
			return client
		},
		//// will create new instance of redis client instead of reusing current instance of `redis.Client`
		// fx.Annotate(
		//	NewRedisClient,
		//	fx.As(new(redis.UniversalClient)),
		// ),
		fx.Annotate(
			NewRedisHealthChecker,
			fx.As(new(contracts.Health)),
			fx.ResultTags(fmt.Sprintf(`group:"%s"`, "healths")),
		),
		provideConfig))

	redisInvokes = fx.Options(
		fx.Invoke(registerHooks),
	)
)

// registerHooks registers the hooks for the redis client.
func registerHooks(
	lc fx.Lifecycle,
	client redis.UniversalClient,
	logger logger.Logger,
) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			return client.Ping(ctx).Err()
		},
		OnStop: func(_ context.Context) error {
			if err := client.Close(); err != nil {
				logger.Errorf("error in closing redis: %v", err)
			} else {
				logger.Info("redis closed gracefully")
			}

			return nil
		},
	})
}
