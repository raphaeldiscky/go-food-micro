package consumercontracts

import (
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/core/messaging/consumer"
	messagingTypes "github.com/raphaeldiscky/go-food-micro/internal/pkg/core/messaging/types"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/rabbitmq/consumer/configurations"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/rabbitmq/types"
)

type ConsumerFactory interface {
	CreateConsumer(
		consumerConfiguration *configurations.RabbitMQConsumerConfiguration,
		isConsumedNotifications ...func(message messagingTypes.IMessage),
	) (consumer.Consumer, error)

	Connection() types.IConnection
}
