// Package messaging provides a messaging utils.
package messaging

import (
	"context"

	"github.com/onsi/ginkgo/v2"

	"github.com/raphaeldiscky/go-food-micro/internal/pkg/core/messaging/bus"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/core/messaging/types"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/core/messaging/utils"
	typeMapper "github.com/raphaeldiscky/go-food-micro/internal/pkg/reflection/typemapper"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/test/hypothesis"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/test/messaging/consumer"
)

// ShouldProduced tests if a message is produced.
func ShouldProduced[T types.IMessage](
	ctx context.Context,
	bus bus.Bus,
	condition func(T) bool,
) hypothesis.Hypothesis[T] {
	hypo := hypothesis.ForT[T](condition)

	bus.IsProduced(func(message types.IMessage) {
		defer ginkgo.GinkgoRecover()
		typ := utils.GetMessageBaseReflectType(typeMapper.GenericInstanceByT[T]())
		if utils.GetMessageBaseReflectType(message) == typ {
			m, ok := message.(T)
			if !ok {
				hypo.Test(ctx, *new(T))
			}
			hypo.Test(ctx, m)
		}
	})

	return hypo
}

// ShouldConsume tests if a message is consumed.
func ShouldConsume[T types.IMessage](
	ctx context.Context,
	bus bus.Bus,
	condition func(T) bool,
) hypothesis.Hypothesis[T] {
	hypo := hypothesis.ForT[T](condition)

	bus.IsConsumed(func(message types.IMessage) {
		defer ginkgo.GinkgoRecover()
		typ := utils.GetMessageBaseReflectType(typeMapper.GenericInstanceByT[T]())
		if utils.GetMessageBaseReflectType(message) == typ {
			m, ok := message.(T)
			if !ok {
				hypo.Test(ctx, *new(T))
			}
			hypo.Test(ctx, m)
		}
	})

	return hypo
}

// ShouldConsumeNewConsumer creates a new consumer and tests if a message is consumed.
func ShouldConsumeNewConsumer[T types.IMessage](bus bus.Bus) (hypothesis.Hypothesis[T], error) {
	hypo := hypothesis.ForT[T](nil)
	testConsumer := consumer.NewRabbitMQFakeTestConsumerHandlerWithHypothesis(hypo)
	err := bus.ConnectConsumerHandler(typeMapper.GenericInstanceByT[T](), testConsumer)
	if err != nil {
		return nil, err
	}

	return hypo, nil
}
