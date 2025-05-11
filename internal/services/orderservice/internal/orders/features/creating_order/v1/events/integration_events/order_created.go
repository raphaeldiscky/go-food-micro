package integrationEvents

import (
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/core/messaging/types"
	dtosV1 "github.com/raphaeldiscky/go-food-micro/internal/services/orderservice/internal/orders/dtos/v1"

	uuid "github.com/satori/go.uuid"
)

type OrderCreatedV1 struct {
	*types.Message
	*dtosV1.OrderReadDto
}

func NewOrderCreatedV1(orderReadDto *dtosV1.OrderReadDto) *OrderCreatedV1 {
	return &OrderCreatedV1{
		OrderReadDto: orderReadDto,
		Message:      types.NewMessage(uuid.NewV4().String()),
	}
}
