// Package rabbitmq provides a set of functions for the rabbitmq package.
package rabbitmq

import (
	"context"
	"fmt"
	"time"

	"go.uber.org/fx"

	bus2 "github.com/raphaeldiscky/go-food-micro/internal/pkg/core/messaging/bus"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/core/messaging/producer"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/health/contracts"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/logger"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/rabbitmq/bus"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/rabbitmq/config"
	rabbitmqconsumer "github.com/raphaeldiscky/go-food-micro/internal/pkg/rabbitmq/consumer"
	rabbitmqproducer "github.com/raphaeldiscky/go-food-micro/internal/pkg/rabbitmq/producer"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/rabbitmq/types"
)

var (
	// ModuleFunc provided to fxlog
	// https://uber-go.github.io/fx/modules.html
	ModuleFunc = func(rabbitMQConfigurationConstructor interface{}) fx.Option {
		return fx.Module(
			"rabbitmqfx",
			fx.Provide(rabbitMQConfigurationConstructor),
			rabbitmqProviders,
			rabbitmqInvokes,
		)
	}

	// - order is not important in provide
	// - provide can have parameter and will resolve if registered
	// - execute its func only if it requested.
	rabbitmqProviders = fx.Options(
		fx.Provide(config.ProvideConfig),
		fx.Provide(types.NewRabbitMQConnection),
		fx.Provide(fx.Annotate(
			bus.NewRabbitmqBus,
			fx.ParamTags(``, ``, ``, `optional:"true"`),
			fx.As(new(producer.Producer)),
			fx.As(new(bus2.Bus)),
			fx.As(new(bus.RabbitmqBus)),
		)),
		fx.Provide(rabbitmqconsumer.NewConsumerFactory),
		fx.Provide(rabbitmqproducer.NewProducerFactory),
		fx.Provide(fx.Annotate(
			NewRabbitMQHealthChecker,
			fx.As(new(contracts.Health)),
			fx.ResultTags(fmt.Sprintf(`group:"%s"`, "healths")),
		)))

	// - execute after registering all of our provided
	// - they execute by their orders
	// - invokes always execute its func compare to provides that only run when we request for them.
	// - return value will be discarded and can not be provided.
	rabbitmqInvokes = fx.Options(
		fx.Invoke(registerHooks),
	)
)

// we don't want to register any dependencies here, its func body should execute always even we don't request for that, so we should use `invoke`.
func registerHooks(
	lc fx.Lifecycle,
	bus bus.RabbitmqBus,
	rabbitmqOptions *config.RabbitmqOptions,
	logger logger.Logger,
) {
	if !rabbitmqOptions.AutoStart {
		return
	}

	lifeTimeCtx := context.Background()

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			// https://github.com/uber-go/fx/blob/v1.20.0/app.go#L573
			// this ctx is just for startup dependencies setup and OnStart callbacks, and it has short timeout 15s, and it is not alive in whole lifetime app
			// if we need an app context which is alive until the app context done we should create it manually here

			go func() {
				// if (ctx.Err() == nil), context not canceled or deadlined
				if err := bus.Start(lifeTimeCtx); err != nil {
					logger.Errorf(
						"(bus.Start) error in running rabbitmq server: {%v}",
						err,
					)

					return
				}
			}()
			logger.Info("rabbitmq is listening.")

			return nil
		},
		OnStop: func(_ context.Context) error {
			// https://github.com/uber-go/fx/blob/v1.20.0/app.go#L573
			// this ctx is just for stopping callbacks or OnStop callbacks, and it has short timeout 15s, and it is not alive in whole lifetime app
			if err := bus.Stop(); err != nil {
				logger.Errorf("error shutting down rabbitmq server: %v", err)
			} else {
				logger.Info("rabbitmq server shutdown gracefully")
			}

			_, cancel := context.WithTimeout(lifeTimeCtx, 5*time.Second)
			defer cancel()

			return nil
		},
	})
}
