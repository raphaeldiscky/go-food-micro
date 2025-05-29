// Package projections contains the elastic order projection.
package projections

import (
	"context"
	"fmt"
	"time"

	"emperror.dev/errors"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/es/contracts/projection"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/es/models"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/logger"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/mapper"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/otel/tracing"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/otel/tracing/attribute"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/otel/tracing/utils"

	attribute2 "go.opentelemetry.io/otel/attribute"

	"github.com/raphaeldiscky/go-food-micro/internal/services/orderservice/internal/orders/contracts/repositories"
	createOrderDomainEventsV1 "github.com/raphaeldiscky/go-food-micro/internal/services/orderservice/internal/orders/features/creatingorder/v1/events/domainevents"
	submitOrderDomainEventsV1 "github.com/raphaeldiscky/go-food-micro/internal/services/orderservice/internal/orders/features/submittingorder/v1/events/domainevents"
	updateOrderDomainEventsV1 "github.com/raphaeldiscky/go-food-micro/internal/services/orderservice/internal/orders/features/updatingshoppingcard/v1/events"
	"github.com/raphaeldiscky/go-food-micro/internal/services/orderservice/internal/orders/models/orders/readmodels"
)

// elasticOrderProjection is the projection for the order.
type elasticOrderProjection struct {
	elasticOrderReadRepository repositories.OrderElasticRepository
	log                        logger.Logger
	tracer                     tracing.AppTracer
}

// NewElasticOrderProjection creates a new elastic order projection.
func NewElasticOrderProjection(
	elasticOrderReadRepository repositories.OrderElasticRepository,
	log logger.Logger,
	tracer tracing.AppTracer,
) projection.IProjection {
	return &elasticOrderProjection{
		elasticOrderReadRepository: elasticOrderReadRepository,
		log:                        log,
		tracer:                     tracer,
	}
}

// ProcessEvent processes the event.
func (e *elasticOrderProjection) ProcessEvent(
	ctx context.Context,
	streamEvent *models.StreamEvent,
) error {
	// Handling and projecting event to elastic read model
	switch evt := streamEvent.Event.(type) {
	case *createOrderDomainEventsV1.OrderCreatedV1:
		return e.onOrderCreated(ctx, evt)
	case *updateOrderDomainEventsV1.ShoppingCartUpdatedV1:
		return e.onShoppingCartUpdated(ctx, evt)
	case *submitOrderDomainEventsV1.OrderSubmittedV1:
		return e.onOrderSubmitted(ctx, evt)
	default:
		return nil
	}
}

// onOrderCreated handles the order created event.
func (e *elasticOrderProjection) onOrderCreated(
	ctx context.Context,
	evt *createOrderDomainEventsV1.OrderCreatedV1,
) error {
	ctx, span := e.tracer.Start(ctx, "elasticOrderProjection.onOrderCreated")
	span.SetAttributes(attribute.Object("Event", evt))
	span.SetAttributes(attribute2.String("OrderId", evt.OrderId.String()))
	defer span.End()

	items, err := mapper.Map[[]*readmodels.ShopItemReadModel](evt.ShopItems)
	if err != nil {
		return errors.WrapIf(
			err,
			"[elasticOrderProjection_onOrderCreated.Map] error in mapping shopItems",
		)
	}

	orderRead := readmodels.NewOrderReadModel(
		evt.OrderId,
		items,
		evt.AccountEmail,
		evt.DeliveryAddress,
		evt.DeliveredTime,
	)

	_, err = e.elasticOrderReadRepository.CreateOrder(ctx, orderRead)
	if err != nil {
		return utils.TraceStatusFromSpan(
			span,
			errors.WrapIf(
				err,
				"[elasticOrderProjection_onOrderCreated.CreateOrder] error in creating order with elasticOrderRepository",
			),
		)
	}

	e.log.Infow(
		fmt.Sprintf(
			"[elasticOrderProjection.onOrderCreated] order with id '%s' created",
			orderRead.OrderId,
		),
		logger.Fields{"ID": orderRead.OrderId},
	)

	return nil
}

// onShoppingCartUpdated handles the shopping cart updated event.
func (e *elasticOrderProjection) onShoppingCartUpdated(
	ctx context.Context,
	evt *updateOrderDomainEventsV1.ShoppingCartUpdatedV1,
) error {
	ctx, span := e.tracer.Start(ctx, "elasticOrderProjection.onShoppingCartUpdated")
	span.SetAttributes(attribute.Object("Event", evt))
	defer span.End()

	items, err := mapper.Map[[]*readmodels.ShopItemReadModel](evt.ShopItems)
	if err != nil {
		return errors.WrapIf(
			err,
			"[elasticOrderProjection_onShoppingCartUpdated.Map] error in mapping shopItems",
		)
	}

	// Get existing order
	order, err := e.elasticOrderReadRepository.GetOrderByID(ctx, evt.GetAggregateId())
	if err != nil {
		return utils.TraceStatusFromSpan(
			span,
			errors.WrapIf(
				err,
				"[elasticOrderProjection_onShoppingCartUpdated.GetOrderByID] error in getting order",
			),
		)
	}

	// Update order with new items
	order.ShopItems = items
	order.TotalPrice = getShopItemsTotalPrice(items)
	order.UpdatedAt = time.Now()

	_, err = e.elasticOrderReadRepository.UpdateOrder(ctx, order)
	if err != nil {
		return utils.TraceStatusFromSpan(
			span,
			errors.WrapIf(
				err,
				"[elasticOrderProjection_onShoppingCartUpdated.UpdateOrder] error in updating order",
			),
		)
	}

	e.log.Infow(
		fmt.Sprintf(
			"[elasticOrderProjection.onShoppingCartUpdated] order with id '%s' updated",
			order.OrderId,
		),
		logger.Fields{"ID": order.OrderId},
	)

	return nil
}

// onOrderSubmitted handles the order submitted event.
func (e *elasticOrderProjection) onOrderSubmitted(
	ctx context.Context,
	evt *submitOrderDomainEventsV1.OrderSubmittedV1,
) error {
	ctx, span := e.tracer.Start(ctx, "elasticOrderProjection.onOrderSubmitted")
	span.SetAttributes(attribute.Object("Event", evt))
	span.SetAttributes(attribute2.String("OrderId", evt.OrderID.String()))
	defer span.End()

	// Get existing order
	order, err := e.elasticOrderReadRepository.GetOrderByID(ctx, evt.OrderID)
	if err != nil {
		return utils.TraceStatusFromSpan(
			span,
			errors.WrapIf(
				err,
				"[elasticOrderProjection_onOrderSubmitted.GetOrderByID] error in getting order",
			),
		)
	}

	// Update order status
	order.Submitted = true
	order.UpdatedAt = time.Now()

	_, err = e.elasticOrderReadRepository.UpdateOrder(ctx, order)
	if err != nil {
		return utils.TraceStatusFromSpan(
			span,
			errors.WrapIf(
				err,
				"[elasticOrderProjection_onOrderSubmitted.UpdateOrder] error in updating order",
			),
		)
	}

	e.log.Infow(
		fmt.Sprintf(
			"[elasticOrderProjection.onOrderSubmitted] order with id '%s' submitted",
			order.OrderId,
		),
		logger.Fields{"ID": order.OrderId},
	)

	return nil
}

// getShopItemsTotalPrice gets the total price of the shop items.
func getShopItemsTotalPrice(shopItems []*readmodels.ShopItemReadModel) float64 {
	var totalPrice float64 = 0
	for _, item := range shopItems {
		totalPrice += item.Price * float64(item.Quantity)
	}

	return totalPrice
}
