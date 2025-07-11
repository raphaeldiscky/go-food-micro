// Package grpc provides a grpc module.
package grpc

import (
	"context"

	"go.uber.org/fx"

	"github.com/raphaeldiscky/go-food-micro/internal/pkg/grpc/config"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/logger"
)

// Module provided to fxlog.
var (
	// Module provided to fxlog
	// https://uber-go.github.io/fx/modules.html
	Module = fx.Module(
		"grpcfx",
		grpcProviders,
		grpcInvokes,
	)

	// grpcProviders is a fx.Options that provides the grpc module.
	// - order is not important in provide
	// - provide can have parameter and will resolve if registered
	// - execute its func only if it requested.
	grpcProviders = fx.Options(fx.Provide(
		config.ProvideConfig,
		// https://uber-go.github.io/fx/value-groups/consume.html#with-annotated-functions
		// https://uber-go.github.io/fx/annotate.html
		fx.Annotate(
			NewGrpcServer,
			fx.ParamTags(``, ``),
		),
		NewGrpcClient,
	))

	// grpcInvokes is a fx.Options that invokes the grpc module.
	// - execute after registering all of our provided
	// - they execute by their orders
	// - invokes always execute its func compare to provides that only run when we request for them.
	// - return value will be discarded and can not be provided.
	grpcInvokes = fx.Options(fx.Invoke(registerHooks))
)

// registerHooks is a function that registers the grpc module.
func registerHooks(
	lc fx.Lifecycle,
	grpcServer GrpcServer,
	grpcClient GrpcClient,
	logger logger.Logger,
	options *config.GrpcOptions,
) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			// https://github.com/uber-go/fx/blob/v1.20.0/app.go#L573
			// this ctx is just for startup dependencies setup and OnStart callbacks, and it has short timeout 15s, and it is not alive in whole lifetime app
			// if we need an app context which is alive until the app context done we should create it manually here
			go func() {
				// if (ctx.Err() == nil), context not canceled or deadlined
				if err := grpcServer.RunGrpcServer(nil); err != nil {
					// do a fatal for going to OnStop process
					logger.Fatalf(
						"(GrpcServer.RunGrpcServer) error in running server: {%v}",
						err,
					)
				}
			}()
			logger.Infof(
				"%s is listening on Host:{%s} Grpc PORT: {%s}",
				options.Name,
				options.Host,
				options.Port,
			)

			return nil
		},
		OnStop: func(ctx context.Context) error {
			// https://github.com/uber-go/fx/blob/v1.20.0/app.go#L573
			// this ctx is just for stopping callbacks or OnStop callbacks, and it has short timeout 15s, and it is not alive in whole lifetime app
			grpcServer.GracefulShutdown()
			logger.Info("server shutdown gracefully")

			if err := grpcClient.Close(); err != nil {
				logger.Errorf("error in closing grpc-client: %v", err)
			} else {
				logger.Info("grpc-client closed gracefully")
			}

			return nil
		},
	})
}
