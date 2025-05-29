// Package bus provides a module for the bus.
package bus

import (
	consumer2 "github.com/raphaeldiscky/go-food-micro/internal/pkg/core/messaging/consumer"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/core/messaging/producer"
)

// Bus is a bus.
type Bus interface {
	producer.Producer
	consumer2.BusControl
	consumer2.ConsumerConnector
}
