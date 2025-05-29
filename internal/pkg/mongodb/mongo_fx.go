// Package mongodb provides a module for the mongodb.
package mongodb

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/raphaeldiscky/go-food-micro/internal/pkg/health/contracts"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/logger"
)

// Module is a module for the mongodb.
var (
	// Module provided to fxlog
	// https://uber-go.github.io/fx/modules.html
	Module = fx.Module(
		"mongofx",
		mongoProviders,
		mongoInvokes,
	)

	// mongoProviders is a module for the mongodb.
	mongoProviders = fx.Provide(
		provideConfig,
		NewMongoDB,
		fx.Annotate(
			NewMongoHealthChecker,
			fx.As(new(contracts.Health)),
			fx.ResultTags(fmt.Sprintf(`group:"%s"`, "healths")),
		),
	)

	mongoInvokes = fx.Invoke(registerHooks)
)

// registerHooks registers hooks for the mongodb.
func registerHooks(
	lc fx.Lifecycle,
	client *mongo.Client,
	logger logger.Logger,
) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			err := client.Ping(ctx, nil)
			if err != nil {
				logger.Error("failed to ping mongo", zap.Error(err))

				return err
			}

			logger.Info("successfully pinged mongo")

			return nil
		},
		OnStop: func(ctx context.Context) error {
			if err := client.Disconnect(ctx); err != nil {
				logger.Errorf("error in disconnecting mongo: %v", err)
			} else {
				logger.Info("mongo disconnected gracefully")
			}

			return nil
		},
	})
}
