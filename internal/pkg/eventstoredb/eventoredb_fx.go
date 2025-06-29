// Package eventstoredb provides a serializer for EventStoreDB.
package eventstoredb

import (
	"context"
	"time"

	"go.uber.org/fx"

	kdb "github.com/kurrent-io/KurrentDB-Client-Go/kurrentdb"

	"github.com/raphaeldiscky/go-food-micro/internal/pkg/eventstoredb/config"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/logger"
)

var (
	// ModuleFunc provided to fxlog
	// https://uber-go.github.io/fx/modules.html
	ModuleFunc = func(projectionBuilderConstructor interface{}) fx.Option {
		return fx.Module(
			"kurrentdbfx",
			fx.Provide(projectionBuilderConstructor),
			eventstoreProviders,
			eventstoreInvokes,
		)
	}

	// - order is not important in provide
	// - provide can have parameter and will resolve if registered
	// - execute its func only if it requested.
	eventstoreProviders = fx.Options(fx.Provide(
		config.ProvideConfig,
		NewEsdbSerializer,
		NewEventStoreDB,
		NewEventStoreDbEventStore,
		NewEsdbSubscriptionCheckpointRepository,
		NewEsdbSubscriptionAllWorker,
	))

	// FiberInvokes - execute after registering all of our provided
	// - they execute by their orders
	// - invokes always execute its func compare to provides that only run when we request for them.
	// - return value will be discarded and can not be provided.
	eventstoreInvokes = fx.Options(fx.Invoke(registerHooks))
)

// registerHooks registers hooks for the event store db.
// we don't want to register any dependencies here, its func body should execute always even we don't request for that, so we should use `invoke`.
func registerHooks(
	lc fx.Lifecycle,
	worker EsdbSubscriptionAllWorker,
	logger logger.Logger,
	cfg *config.EventStoreDbOptions,
) {
	lifetimeCtx := context.Background()

	lc.Append(fx.Hook{
		OnStart: func(_ context.Context) error {
			// https://github.com/uber-go/fx/blob/v1.20.0/app.go#L573
			// this ctx is just for startup dependencies setup and OnStart callbacks, and it has short timeout 15s, and it is not alive in whole lifetime app
			// if we need an app context which is alive until the app context done we should create it manually here
			go func() {
				option := &EventStoreDBSubscriptionToAllOptions{
					FilterOptions: &kdb.SubscriptionFilter{
						Type:     kdb.StreamFilterType,
						Prefixes: cfg.Subscription.Prefix,
					},
					SubscriptionId: cfg.Subscription.SubscriptionId,
				}
				if err := worker.SubscribeAll(lifetimeCtx, option); err != nil {
					logger.Errorf(
						"(worker.SubscribeAll) error in running esdb subscription worker: {%v}",
						err,
					)

					return
				}
			}()
			logger.Info("esdb subscription worker is listening.")

			return nil
		},
		OnStop: func(ctx context.Context) error {
			// https://github.com/uber-go/fx/blob/v1.20.0/app.go#L573
			// this ctx is just for stopping callbacks or OnStop callbacks, and it has short timeout 15s, and it is not alive in whole lifetime app
			_, cancel := context.WithTimeout(lifetimeCtx, 5*time.Second)
			defer cancel()

			return nil
		},
	})
}
