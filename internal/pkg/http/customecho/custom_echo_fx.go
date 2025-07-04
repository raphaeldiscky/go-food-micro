// Package customecho provides a custom echo server.
package customecho

import (
	"context"
	"errors"
	"net/http"

	"go.uber.org/fx"

	"github.com/raphaeldiscky/go-food-micro/internal/pkg/http/customecho/config"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/http/customecho/contracts"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/logger"
)

// Module provided to fxlog
// https://uber-go.github.io/fx/modules.html
var (
	Module = fx.Module(
		"customechofx",

		echoProviders,
		echoInvokes,
	)

	// - order is not important in provide
	// - provide can have parameter and will resolve if registered
	// - execute its func only if it requested.
	echoProviders = fx.Options(fx.Provide(
		config.ProvideConfig,
		// https://uber-go.github.io/fx/value-groups/consume.html#with-annotated-functions
		// https://uber-go.github.io/fx/annotate.html
		fx.Annotate(
			NewEchoHTTPServer,
			fx.ParamTags(``, ``, `optional:"true"`),
		),
	))

	// - execute after registering all of our provided
	// - they execute by their orders
	// - invokes always execute its func compare to provides that only run when we request for them.
	// - return value will be discarded and can not be provided.
	echoInvokes = fx.Options(fx.Invoke(registerHooks))
)

// registerHooks registers the hooks.
// we don't want to register any dependencies here, its func body should execute always even we don't request for that, so we should use `invoke`.
func registerHooks(
	lc fx.Lifecycle,
	echoServer contracts.EchoHTTPServer,
	logger logger.Logger,
) {
	lc.Append(fx.Hook{
		OnStart: func(_ context.Context) error {
			// https://github.com/uber-go/fx/blob/v1.20.0/app.go#L573
			// this ctx is just for startup dependencies setup and OnStart callbacks, and it has short timeout 15s, and it is not alive in whole lifetime app
			// if we need an app context which is alive until the app context done we should create it manually here

			go func() {
				// https://medium.com/@mokiat/proper-http-shutdown-in-go-bd3bfaade0f2
				// When Shutdown is called, Serve, ListenAndServe, and ListenAndServeTLS immediately return ErrServerClosed. Make sure the program doesn’t exit and waits instead for Shutdown to return.
				if err := echoServer.RunHTTPServer(); !errors.Is(
					err,
					http.ErrServerClosed,
				) {
					// do a fatal for going to OnStop process
					logger.Fatalf(
						"(EchoHTTPServer.RunHTTPServer) error in running server: {%v}",
						err,
					)
				}
			}()
			echoServer.Logger().Infof(
				"%s is listening on Host:{%s} Http PORT: {%s}",
				echoServer.Cfg().Name,
				echoServer.Cfg().Host,
				echoServer.Cfg().Port,
			)

			return nil
		},
		OnStop: func(ctx context.Context) error {
			// https://github.com/uber-go/fx/blob/v1.20.0/app.go#L573
			// this ctx is just for stopping callbacks or OnStop callbacks, and it has short timeout 15s, and it is not alive in whole lifetime app
			// https://medium.com/@mokiat/proper-http-shutdown-in-go-bd3bfaade0f2
			// When Shutdown is called, Serve, ListenAndServe, and ListenAndServeTLS immediately return ErrServerClosed. Make sure the program doesn’t exit and waits instead for Shutdown to return.
			if err := echoServer.GracefulShutdown(ctx); err != nil {
				echoServer.Logger().
					Errorf("error shutting down echo server: %v", err)
			} else {
				echoServer.Logger().Info("echo server shutdown gracefully")
			}

			return nil
		},
	})
}
