package producercontracts

import (
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/core/messaging/producer"
	types2 "github.com/raphaeldiscky/go-food-micro/internal/pkg/core/messaging/types"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/rabbitmq/producer/configurations"
)

type ProducerFactory interface {
	CreateProducer(
		rabbitmqProducersConfiguration map[string]*configurations.RabbitMQProducerConfiguration,
		isProducedNotifications ...func(message types2.IMessage),
	) (producer.Producer, error)
}
