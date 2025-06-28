// Package repositories contains the mongo order read repository.
package repositories

import (
	"context"
	"fmt"

	"emperror.dev/errors"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/logger"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/mongodb"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/otel/tracing"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/otel/tracing/attribute"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	utils2 "github.com/raphaeldiscky/go-food-micro/internal/pkg/otel/tracing/utils"
	goUuid "github.com/satori/go.uuid"
	attribute2 "go.opentelemetry.io/otel/attribute"

	"github.com/raphaeldiscky/go-food-micro/internal/services/orderservice/internal/orders/contracts/repositories"
	"github.com/raphaeldiscky/go-food-micro/internal/services/orderservice/internal/orders/models/orders/readmodels"
)

const (
	orderCollection = "orders"
)

// mongoOrderReadRepository is the mongo order read repository.
type mongoOrderReadRepository struct {
	log          logger.Logger
	mongoOptions *mongodb.MongoDbOptions
	mongoClient  *mongo.Client
	tracer       tracing.AppTracer
}

// NewMongoOrderReadRepository creates a new mongo order read repository.
func NewMongoOrderReadRepository(
	log logger.Logger,
	cfg *mongodb.MongoDbOptions,
	mongoClient *mongo.Client,
	tracer tracing.AppTracer,
) repositories.OrderMongoRepository {
	return &mongoOrderReadRepository{
		log:          log,
		mongoOptions: cfg,
		mongoClient:  mongoClient,
		tracer:       tracer,
	}
}

// GetAllOrders gets all orders from the database.
func (m mongoOrderReadRepository) GetAllOrders(
	ctx context.Context,
	listQuery *utils.ListQuery,
) (*utils.ListResult[*readmodels.OrderReadModel], error) {
	ctx, span := m.tracer.Start(ctx, "mongoOrderReadRepository.GetAllOrders")
	defer span.End()

	collection := m.mongoClient.Database(m.mongoOptions.Database).Collection(orderCollection)

	result, err := mongodb.Paginate[*readmodels.OrderReadModel](ctx, listQuery, collection, nil)
	if err != nil {
		return nil, utils2.TraceStatusFromContext(
			ctx,
			errors.WrapIf(
				err,
				"[mongoOrderReadRepository_GetAllOrders.Paginate] error in the paginate",
			),
		)
	}

	m.log.Infow(
		"[mongoOrderReadRepository.GetAllOrders] orders loaded",
		logger.Fields{"OrdersResult": result},
	)

	span.SetAttributes(attribute.Object("OrdersResult", result))

	return result, nil
}

// SearchOrders searches for orders in the database.
func (m mongoOrderReadRepository) SearchOrders(
	ctx context.Context,
	searchText string,
	listQuery *utils.ListQuery,
) (*utils.ListResult[*readmodels.OrderReadModel], error) {
	ctx, span := m.tracer.Start(ctx, "mongoOrderReadRepository.SearchOrders")
	span.SetAttributes(attribute2.String("SearchText", searchText))
	defer span.End()

	collection := m.mongoClient.Database(m.mongoOptions.Database).Collection(orderCollection)

	filter := bson.D{
		{Key: "$or", Value: bson.A{
			bson.D{{Key: "name", Value: primitive.Regex{Pattern: searchText, Options: "gi"}}},
			bson.D{
				{Key: "description", Value: primitive.Regex{Pattern: searchText, Options: "gi"}},
			},
		}},
	}

	result, err := mongodb.Paginate[*readmodels.OrderReadModel](ctx, listQuery, collection, filter)
	if err != nil {
		return nil, utils2.TraceStatusFromContext(
			ctx,
			errors.WrapIf(
				err,
				"[mongoOrderReadRepository_SearchOrders.Paginate] error in the paginate",
			),
		)
	}
	span.SetAttributes(attribute.Object("OrdersResult", result))

	m.log.Infow(
		fmt.Sprintf(
			"[mongoOrderReadRepository.SearchOrders] orders loaded for search term '%s'",
			searchText,
		),
		logger.Fields{"OrdersResult": result},
	)

	return result, nil
}

// GetOrderByID gets an order by id from the database.
func (m mongoOrderReadRepository) GetOrderByID(
	ctx context.Context,
	id goUuid.UUID,
) (*readmodels.OrderReadModel, error) {
	ctx, span := m.tracer.Start(ctx, "mongoOrderReadRepository.GetOrderByID")
	span.SetAttributes(attribute2.String("ID", id.String()))
	defer span.End()

	collection := m.mongoClient.Database(m.mongoOptions.Database).Collection(orderCollection)

	var order readmodels.OrderReadModel
	if err := collection.FindOne(ctx, bson.M{"_id": id.String()}).Decode(&order); err != nil {
		// ErrNoDocuments means that the filter did not match any documents in the collection
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}

		return nil, utils2.TraceStatusFromContext(
			ctx,
			errors.WrapIf(
				err,
				fmt.Sprintf(
					"[mongoOrderReadRepository_GetOrderByID.FindOne] can't find the order with id %s into the database.",
					id,
				),
			),
		)
	}
	span.SetAttributes(attribute.Object("Order", order))

	m.log.Infow(
		fmt.Sprintf("[mongoOrderReadRepository.GetOrderByID] order with id %s loaded", id.String()),
		logger.Fields{"Order": order, "ID": id},
	)

	return &order, nil
}

// GetOrderByOrderID gets an order by order id from the database.
func (m mongoOrderReadRepository) GetOrderByOrderID(
	ctx context.Context,
	orderID goUuid.UUID,
) (*readmodels.OrderReadModel, error) {
	ctx, span := m.tracer.Start(ctx, "mongoOrderReadRepository.GetOrderByOrderID")
	span.SetAttributes(attribute2.String("OrderID", orderID.String()))
	defer span.End()

	collection := m.mongoClient.Database(m.mongoOptions.Database).Collection(orderCollection)

	var order readmodels.OrderReadModel
	if err := collection.FindOne(ctx, bson.M{"orderId": orderID.String()}).Decode(&order); err != nil {
		// ErrNoDocuments means that the filter did not match any documents in the collection
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}

		return nil, utils2.TraceStatusFromContext(
			ctx,
			errors.WrapIf(
				err,
				fmt.Sprintf(
					"[mongoOrderReadRepository_GetOrderByID.FindOne] can't find the order with orderId %s into the database.",
					orderID.String(),
				),
			),
		)
	}
	span.SetAttributes(attribute.Object("Order", order))

	m.log.Infow(
		fmt.Sprintf(
			"[mongoOrderReadRepository.GetOrderByID] order with orderId %s loaded",
			orderID.String(),
		),
		logger.Fields{"Order": order, "orderId": orderID},
	)

	return &order, nil
}

// CreateOrder creates an order in the database.
func (m mongoOrderReadRepository) CreateOrder(
	ctx context.Context,
	order *readmodels.OrderReadModel,
) (*readmodels.OrderReadModel, error) {
	ctx, span := m.tracer.Start(ctx, "mongoOrderReadRepository.CreateOrder")
	defer span.End()

	collection := m.mongoClient.Database(m.mongoOptions.Database).Collection(orderCollection)
	_, err := collection.InsertOne(ctx, order, &options.InsertOneOptions{})
	if err != nil {
		return nil, utils2.TraceStatusFromContext(
			ctx,
			errors.WrapIf(
				err,
				"[mongoOrderReadRepository_CreateOrder.InsertOne] error in the inserting order into the database.",
			),
		)
	}
	span.SetAttributes(attribute.Object("Order", order))

	m.log.Infow(
		fmt.Sprintf(
			"[mongoOrderReadRepository.CreateOrder] order with id '%s' created",
			order.OrderID,
		),
		logger.Fields{"Order": order, "ID": order.OrderID},
	)

	return order, nil
}

// UpdateOrder updates an order in the database.
func (m mongoOrderReadRepository) UpdateOrder(
	ctx context.Context,
	order *readmodels.OrderReadModel,
) (*readmodels.OrderReadModel, error) {
	ctx, span := m.tracer.Start(ctx, "mongoOrderReadRepository.UpdateOrder")
	defer span.End()

	collection := m.mongoClient.Database(m.mongoOptions.Database).Collection(orderCollection)

	ops := options.FindOneAndUpdate()
	ops.SetReturnDocument(options.After)
	ops.SetUpsert(true)

	var updated readmodels.OrderReadModel
	if err := collection.FindOneAndUpdate(ctx, bson.M{"_id": order.OrderID}, bson.M{"$set": order}, ops).Decode(&updated); err != nil {
		return nil, utils2.TraceStatusFromContext(
			ctx,
			errors.WrapIf(
				err,
				fmt.Sprintf(
					"[mongoOrderReadRepository_UpdateOrder.FindOneAndUpdate] error in updating order with id %s into the database.",
					order.OrderID,
				),
			),
		)
	}
	span.SetAttributes(attribute.Object("Order", order))

	m.log.Infow(
		fmt.Sprintf(
			"[mongoOrderReadRepository.UpdateOrder] order with id '%s' updated",
			order.OrderID,
		),
		logger.Fields{"Order": order, "ID": order.OrderID},
	)

	return &updated, nil
}

// DeleteOrderByID deletes an order by id from the database.
func (m mongoOrderReadRepository) DeleteOrderByID(ctx context.Context, uuid goUuid.UUID) error {
	ctx, span := m.tracer.Start(ctx, "mongoOrderReadRepository.DeleteOrderByID")
	span.SetAttributes(attribute2.String("ID", uuid.String()))
	defer span.End()

	collection := m.mongoClient.Database(m.mongoOptions.Database).Collection(orderCollection)

	if err := collection.FindOneAndDelete(ctx, bson.M{"_id": uuid.String()}).Err(); err != nil {
		return utils2.TraceStatusFromContext(ctx, errors.WrapIf(err, fmt.Sprintf(
			"[mongoOrderReadRepository_DeleteOrderByID.FindOneAndDelete] error in deleting order with id %d from the database.",
			uuid,
		)))
	}

	m.log.Infow(
		fmt.Sprintf("[mongoOrderReadRepository.DeleteOrderByID] order with id %s deleted", uuid),
		logger.Fields{"ID": uuid},
	)

	return nil
}
