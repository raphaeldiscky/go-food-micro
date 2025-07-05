package projections

import (
	"context"
	"fmt"

	"emperror.dev/errors"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/core/messaging/producer"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/es/contracts/projection"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/es/models"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/logger"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/mapper"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/otel/tracing"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/otel/tracing/attribute"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/otel/tracing/utils"

	customErrors "github.com/raphaeldiscky/go-food-micro/internal/pkg/http/httperrors/customerrors"
	attribute2 "go.opentelemetry.io/otel/attribute"

	"github.com/raphaeldiscky/go-food-micro/internal/services/orderservice/internal/orders/contracts/repositories"
	dtosV1 "github.com/raphaeldiscky/go-food-micro/internal/services/orderservice/internal/orders/dtos/v1"
	createOrderDomainEventsV1 "github.com/raphaeldiscky/go-food-micro/internal/services/orderservice/internal/orders/features/creatingorder/v1/events/domainevents"
	createOrderIntegrationEventsV1 "github.com/raphaeldiscky/go-food-micro/internal/services/orderservice/internal/orders/features/creatingorder/v1/events/integrationevents"
	"github.com/raphaeldiscky/go-food-micro/internal/services/orderservice/internal/orders/models/orders/readmodels"
)

// mongoOrderProjection is the mongo order projection.
type mongoOrderProjection struct {
	mongoOrderRepository repositories.OrderMongoRepository
	rabbitmqProducer     producer.Producer
	logger               logger.Logger
	tracer               tracing.AppTracer
}

// NewMongoOrderProjection creates a new mongo order projection.
func NewMongoOrderProjection(
	mongoOrderRepository repositories.OrderMongoRepository,
	rabbitmqProducer producer.Producer,
	log logger.Logger,
	tracer tracing.AppTracer,
) projection.IProjection {
	return &mongoOrderProjection{
		mongoOrderRepository: mongoOrderRepository,
		rabbitmqProducer:     rabbitmqProducer,
		logger:               log,
		tracer:               tracer,
	}
}

// ProcessEvent processes the event.
func (m *mongoOrderProjection) ProcessEvent(
	ctx context.Context,
	streamEvent *models.StreamEvent,
) error {
	// Handling and projecting event to elastic read model
	if evt, ok := streamEvent.Event.(*createOrderDomainEventsV1.OrderCreatedV1); ok {
		return m.onOrderCreated(ctx, evt)
	}

	return nil
}

// onOrderCreated handles the order created event.
func (m *mongoOrderProjection) onOrderCreated(
	ctx context.Context,
	evt *createOrderDomainEventsV1.OrderCreatedV1,
) error {
	ctx, span := m.tracer.Start(ctx, "mongoOrderProjection.onOrderCreated")
	span.SetAttributes(attribute.Object("Event", evt))
	span.SetAttributes(attribute2.String("OrderID", evt.OrderID.String()))
	defer span.End()

	items, err := mapper.Map[[]*readmodels.ShopItemReadModel](evt.ShopItems)
	if err != nil {
		return errors.WrapIf(
			err,
			"[mongoOrderProjection_onOrderCreated.Map] error in mapping shopItems",
		)
	}

	orderRead := readmodels.NewOrderReadModel(
		evt.OrderID,
		items,
		evt.AccountEmail,
		evt.DeliveryAddress,
		evt.DeliveredTime,
	)
	_, err = m.mongoOrderRepository.CreateOrder(ctx, orderRead)
	if err != nil {
		return utils.TraceStatusFromSpan(
			span,
			errors.WrapIf(
				err,
				"[mongoOrderProjection_onOrderCreated.CreateOrder] error in creating order with mongoOrderRepository",
			),
		)
	}

	orderReadDto, err := mapper.Map[*dtosV1.OrderReadDto](orderRead)
	if err != nil {
		return utils.TraceErrStatusFromSpan(
			span,
			customErrors.NewApplicationErrorWrap(
				err,
				"[mongoOrderProjection_onOrderCreated.Map] error in mapping OrderReadDto",
			),
		)
	}

	orderCreatedEvent := createOrderIntegrationEventsV1.NewOrderCreatedV1(orderReadDto)

	err = m.rabbitmqProducer.PublishMessage(ctx, orderCreatedEvent, nil)
	if err != nil {
		return utils.TraceErrStatusFromSpan(
			span,
			customErrors.NewApplicationErrorWrap(
				err,
				"[mongoOrderProjection_onOrderCreated.PublishMessage] error in publishing OrderCreated integration_events event",
			),
		)
	}

	m.logger.Infow(
		fmt.Sprintf(
			"[mongoOrderProjection.onOrderCreated] OrderCreated message with messageId `%s` published to the rabbitmq broker",
			orderCreatedEvent.MessageId,
		),
		logger.Fields{"MessageId": orderCreatedEvent.MessageId, "ID": orderCreatedEvent.OrderID},
	)

	m.logger.Infow(
		fmt.Sprintf(
			"[mongoOrderProjection.onOrderCreated] order with id '%s' created",
			orderCreatedEvent.ID,
		),
		logger.Fields{"ID": orderRead.ID, "MessageId": orderCreatedEvent.MessageId},
	)

	return nil
}
