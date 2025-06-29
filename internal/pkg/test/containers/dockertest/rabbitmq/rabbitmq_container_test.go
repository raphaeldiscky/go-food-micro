//go:build integration
// +build integration

// Package rabbitmq provides a rabbitmq docker test.
package rabbitmq

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"

	uuid "github.com/satori/go.uuid"

	"github.com/raphaeldiscky/go-food-micro/internal/pkg/config"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/config/environment"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/core"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/core/messaging/bus"
	messageConsumer "github.com/raphaeldiscky/go-food-micro/internal/pkg/core/messaging/consumer"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/core/messaging/types"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/logger"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/logger/external/fxlog"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/logger/zap"
	rabbitmq2 "github.com/raphaeldiscky/go-food-micro/internal/pkg/rabbitmq"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/rabbitmq/configurations"
	consumerConfigurations "github.com/raphaeldiscky/go-food-micro/internal/pkg/rabbitmq/consumer/configurations"
	typeMapper "github.com/raphaeldiscky/go-food-micro/internal/pkg/reflection/typemapper"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/test/messaging/consumer"
	testUtils "github.com/raphaeldiscky/go-food-micro/internal/pkg/test/utils"
)

// TestRabbitMQContainer tests the rabbitmq container.
func TestRabbitMQContainer(t *testing.T) {
	t.Skip(
		"Skipping RabbitMQ dockertest container test due to persistent port conflicts. See issue with dockertest port allocation.",
	)

	ctx := context.Background()
	fakeConsumer := consumer.NewRabbitMQFakeTestConsumerHandler[*ProducerConsumerMessage]()

	var rabbitmqbus bus.Bus

	fxtest.New(t,
		config.ModuleFunc(environment.Test),
		zap.Module,
		fxlog.FxLogger,
		core.Module,
		rabbitmq2.ModuleFunc(
			func(l logger.Logger) configurations.RabbitMQConfigurationBuilderFuc {
				return func(builder configurations.RabbitMQConfigurationBuilder) {
					builder.AddConsumer(
						&ProducerConsumerMessage{}, // Use pointer type for interface compatibility
						func(consumerBuilder consumerConfigurations.RabbitMQConsumerConfigurationBuilder) {
							consumerBuilder.WithHandlers(
								func(handlerBuilder messageConsumer.ConsumerHandlerConfigurationBuilder) {
									handlerBuilder.AddHandler(fakeConsumer)
								},
							)
						},
					)
				}
			},
		),
		fx.Decorate(RabbitmqDockerTestContainerOptionsDecorator(t, ctx)),
		fx.Populate(&rabbitmqbus),
	).RequireStart()

	require.NotNil(t, rabbitmqbus)

	err := rabbitmqbus.Start(ctx)
	require.NoError(t, err)

	//// wait for consumers ready to consume before publishing messages (for preventing messages lost)
	time.Sleep(time.Second * 2)

	err = rabbitmqbus.PublishMessage(
		context.Background(),
		&ProducerConsumerMessage{
			Data:    "ssssssssss",
			Message: *types.NewMessage(uuid.NewV4().String()), // Dereference to get value instead of pointer
		},
		nil,
	)
	if err != nil {
		return
	}

	err = testUtils.WaitUntilConditionMet(func() bool {
		return fakeConsumer.IsHandled()
	})

	t.Log("stopping test container")

	if err != nil {
		require.FailNow(t, err.Error())
	}
}

// ProducerConsumerMessage is a struct that represents a producer consumer message.
type ProducerConsumerMessage struct {
	types.Message // Remove pointer embedding - use value embedding instead
	Data          string
}

// GetMessageTypeName overrides the embedded method to return the correct type name
func (p *ProducerConsumerMessage) GetMessageTypeName() string {
	return typeMapper.GetTypeName(p)
}

// GetMessageFullTypeName overrides the embedded method to return the correct full type name
func (p *ProducerConsumerMessage) GetMessageFullTypeName() string {
	return typeMapper.GetFullTypeName(p)
}
