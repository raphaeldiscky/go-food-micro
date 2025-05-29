// Package repositories contains the elastic order read repository.
package repositories

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"emperror.dev/errors"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/logger"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/otel/tracing"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/otel/tracing/attribute"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/utils"

	elasticsearch "github.com/elastic/go-elasticsearch/v8"
	uuid "github.com/satori/go.uuid"
	attribute2 "go.opentelemetry.io/otel/attribute"

	"github.com/raphaeldiscky/go-food-micro/internal/services/orderservice/internal/orders/contracts/repositories"
	"github.com/raphaeldiscky/go-food-micro/internal/services/orderservice/internal/orders/models/orders/readmodels"
)

const (
	orderIndex = "orders"
)

// closeResponseBody is a helper function to close the response body and handle any errors.
func closeResponseBody(body io.ReadCloser) error {
	if err := body.Close(); err != nil {
		return errors.WrapIf(err, "failed to close response body")
	}

	return nil
}

// elasticOrderReadRepository is the repository for the order read model.
type elasticOrderReadRepository struct {
	log           logger.Logger
	elasticClient *elasticsearch.Client
	tracer        tracing.AppTracer
}

// NewElasticOrderReadRepository creates a new elastic order read repository.
func NewElasticOrderReadRepository(
	log logger.Logger,
	elasticClient *elasticsearch.Client,
	tracer tracing.AppTracer,
) repositories.OrderElasticRepository {
	return &elasticOrderReadRepository{log: log, elasticClient: elasticClient, tracer: tracer}
}

// GetAllOrders gets all orders.
func (e elasticOrderReadRepository) GetAllOrders(
	ctx context.Context,
	listQuery *utils.ListQuery,
) (*utils.ListResult[*readmodels.OrderReadModel], error) {
	ctx, span := e.tracer.Start(ctx, "elasticOrderReadRepository.GetAllOrders")
	defer span.End()

	// Build the search query
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"match_all": map[string]interface{}{},
		},
		"from": (listQuery.Page - 1) * listQuery.Size,
		"size": listQuery.Size,
		"sort": []map[string]interface{}{
			{"createdAt": "desc"},
		},
	}

	// Convert query to JSON
	queryJSON, err := json.Marshal(query)
	if err != nil {
		return nil, errors.WrapIf(err, "failed to marshal search query")
	}

	// Execute search
	res, err := e.elasticClient.Search(
		e.elasticClient.Search.WithContext(ctx),
		e.elasticClient.Search.WithIndex(orderIndex),
		e.elasticClient.Search.WithBody(strings.NewReader(string(queryJSON))),
	)
	if err != nil {
		return nil, errors.WrapIf(err, "failed to execute search")
	}
	defer func() {
		if closeErr := closeResponseBody(res.Body); closeErr != nil {
			e.log.Error(closeErr)
		}
	}()

	if res.IsError() {
		return nil, fmt.Errorf("search error: %s", res.String())
	}

	// Parse response
	var result struct {
		Hits struct {
			Total struct {
				Value int `json:"value"`
			} `json:"total"`
			Hits []struct {
				Source readmodels.OrderReadModel `json:"_source"`
			} `json:"hits"`
		} `json:"hits"`
	}

	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, errors.WrapIf(err, "failed to decode search response")
	}

	// Convert hits to orders
	orders := make([]*readmodels.OrderReadModel, len(result.Hits.Hits))
	for i := range result.Hits.Hits {
		orders[i] = &result.Hits.Hits[i].Source
	}

	return &utils.ListResult[*readmodels.OrderReadModel]{
		Items:      orders,
		TotalItems: int64(result.Hits.Total.Value),
		Page:       listQuery.Page,
		Size:       listQuery.Size,
		TotalPage:  (result.Hits.Total.Value + listQuery.Size - 1) / listQuery.Size,
	}, nil
}

// SearchOrders searches for orders.
func (e elasticOrderReadRepository) SearchOrders(
	ctx context.Context,
	searchText string,
	listQuery *utils.ListQuery,
) (*utils.ListResult[*readmodels.OrderReadModel], error) {
	ctx, span := e.tracer.Start(ctx, "elasticOrderReadRepository.SearchOrders")
	span.SetAttributes(attribute2.String("SearchText", searchText))
	defer span.End()

	// Build the search query
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"multi_match": map[string]interface{}{
				"query": searchText,
				"fields": []string{
					"accountEmail",
					"deliveryAddress",
					"shopItems.title",
					"shopItems.description",
				},
				"type": "best_fields",
			},
		},
		"from": (listQuery.Page - 1) * listQuery.Size,
		"size": listQuery.Size,
		"sort": []map[string]interface{}{
			{"createdAt": "desc"},
		},
	}

	// Convert query to JSON
	queryJSON, err := json.Marshal(query)
	if err != nil {
		return nil, errors.WrapIf(err, "failed to marshal search query")
	}

	// Execute search
	res, err := e.elasticClient.Search(
		e.elasticClient.Search.WithContext(ctx),
		e.elasticClient.Search.WithIndex(orderIndex),
		e.elasticClient.Search.WithBody(strings.NewReader(string(queryJSON))),
	)
	if err != nil {
		return nil, errors.WrapIf(err, "failed to execute search")
	}
	defer func() {
		if closeErr := closeResponseBody(res.Body); closeErr != nil {
			e.log.Error(closeErr)
		}
	}()

	if res.IsError() {
		return nil, fmt.Errorf("search error: %s", res.String())
	}

	// Parse response
	var result struct {
		Hits struct {
			Total struct {
				Value int `json:"value"`
			} `json:"total"`
			Hits []struct {
				Source readmodels.OrderReadModel `json:"_source"`
			} `json:"hits"`
		} `json:"hits"`
	}

	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, errors.WrapIf(err, "failed to decode search response")
	}

	// Convert hits to orders
	orders := make([]*readmodels.OrderReadModel, len(result.Hits.Hits))
	for i := range result.Hits.Hits {
		orders[i] = &result.Hits.Hits[i].Source
	}

	return &utils.ListResult[*readmodels.OrderReadModel]{
		Items:      orders,
		TotalItems: int64(result.Hits.Total.Value),
		Page:       listQuery.Page,
		Size:       listQuery.Size,
		TotalPage:  (result.Hits.Total.Value + listQuery.Size - 1) / listQuery.Size,
	}, nil
}

// GetOrderByID gets an order by id.
func (e elasticOrderReadRepository) GetOrderByID(
	ctx context.Context,
	id uuid.UUID,
) (*readmodels.OrderReadModel, error) {
	ctx, span := e.tracer.Start(ctx, "elasticOrderReadRepository.GetOrderByID")
	span.SetAttributes(attribute2.String("ID", id.String()))
	defer span.End()

	res, err := e.elasticClient.Get(
		orderIndex,
		id.String(),
		e.elasticClient.Get.WithContext(ctx),
	)
	if err != nil {
		return nil, errors.WrapIf(err, "failed to get order")
	}
	defer func() {
		if closeErr := closeResponseBody(res.Body); closeErr != nil {
			e.log.Error(closeErr)
		}
	}()

	if res.StatusCode == 404 {
		return nil, nil
	}

	if res.IsError() {
		return nil, fmt.Errorf("get error: %s", res.String())
	}

	var result struct {
		Source readmodels.OrderReadModel `json:"_source"`
	}

	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, errors.WrapIf(err, "failed to decode get response")
	}

	span.SetAttributes(attribute.Object("Order", result.Source))

	e.log.Infow(
		fmt.Sprintf(
			"[elasticOrderReadRepository.GetOrderByID] order with id %s loaded",
			id.String(),
		),
		logger.Fields{"Order": result.Source, "ID": id},
	)

	return &result.Source, nil
}

// GetOrderByOrderID gets an order by order id.
func (e elasticOrderReadRepository) GetOrderByOrderID(
	ctx context.Context,
	orderID uuid.UUID,
) (*readmodels.OrderReadModel, error) {
	ctx, span := e.tracer.Start(ctx, "elasticOrderReadRepository.GetOrderByOrderID")
	span.SetAttributes(attribute2.String("OrderID", orderID.String()))
	defer span.End()

	// Build the search query
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"term": map[string]interface{}{
				"orderId": orderID.String(),
			},
		},
	}

	// Convert query to JSON
	queryJSON, err := json.Marshal(query)
	if err != nil {
		return nil, errors.WrapIf(err, "failed to marshal search query")
	}

	// Execute search
	res, err := e.elasticClient.Search(
		e.elasticClient.Search.WithContext(ctx),
		e.elasticClient.Search.WithIndex(orderIndex),
		e.elasticClient.Search.WithBody(strings.NewReader(string(queryJSON))),
	)
	if err != nil {
		return nil, errors.WrapIf(err, "failed to execute search")
	}
	defer func() {
		if closeErr := closeResponseBody(res.Body); closeErr != nil {
			e.log.Error(closeErr)
		}
	}()

	if res.IsError() {
		return nil, fmt.Errorf("search error: %s", res.String())
	}

	// Parse response
	var result struct {
		Hits struct {
			Hits []struct {
				Source readmodels.OrderReadModel `json:"_source"`
			} `json:"hits"`
		} `json:"hits"`
	}

	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, errors.WrapIf(err, "failed to decode search response")
	}

	if len(result.Hits.Hits) == 0 {
		return nil, nil
	}

	order := result.Hits.Hits[0].Source
	span.SetAttributes(attribute.Object("Order", order))

	e.log.Infow(
		fmt.Sprintf(
			"[elasticOrderReadRepository.GetOrderByOrderID] order with orderId %s loaded",
			orderID.String(),
		),
		logger.Fields{"Order": order, "OrderID": orderID},
	)

	return &order, nil
}

// CreateOrder creates an order.
func (e elasticOrderReadRepository) CreateOrder(
	ctx context.Context,
	order *readmodels.OrderReadModel,
) (*readmodels.OrderReadModel, error) {
	ctx, span := e.tracer.Start(ctx, "elasticOrderReadRepository.CreateOrder")
	span.SetAttributes(attribute.Object("Order", order))
	defer span.End()

	// Convert order to JSON
	orderJSON, err := json.Marshal(order)
	if err != nil {
		return nil, errors.WrapIf(err, "failed to marshal order")
	}

	// Create document
	res, err := e.elasticClient.Index(
		orderIndex,
		strings.NewReader(string(orderJSON)),
		e.elasticClient.Index.WithContext(ctx),
		e.elasticClient.Index.WithDocumentID(order.ID),
		e.elasticClient.Index.WithRefresh("true"),
	)
	if err != nil {
		return nil, errors.WrapIf(err, "failed to create order")
	}
	defer func() {
		if closeErr := closeResponseBody(res.Body); closeErr != nil {
			e.log.Error(closeErr)
		}
	}()

	if res.IsError() {
		return nil, fmt.Errorf("index error: %s", res.String())
	}

	e.log.Infow(
		fmt.Sprintf("[elasticOrderReadRepository.CreateOrder] order with id %s created", order.ID),
		logger.Fields{"Order": order},
	)

	return order, nil
}

// UpdateOrder updates an order.
func (e elasticOrderReadRepository) UpdateOrder(
	ctx context.Context,
	order *readmodels.OrderReadModel,
) (*readmodels.OrderReadModel, error) {
	ctx, span := e.tracer.Start(ctx, "elasticOrderReadRepository.UpdateOrder")
	span.SetAttributes(attribute.Object("Order", order))
	defer span.End()

	// Convert order to JSON
	orderJSON, err := json.Marshal(order)
	if err != nil {
		return nil, errors.WrapIf(err, "failed to marshal order")
	}

	// Update document
	res, err := e.elasticClient.Index(
		orderIndex,
		strings.NewReader(string(orderJSON)),
		e.elasticClient.Index.WithContext(ctx),
		e.elasticClient.Index.WithDocumentID(order.ID),
		e.elasticClient.Index.WithRefresh("true"),
	)
	if err != nil {
		return nil, errors.WrapIf(err, "failed to update order")
	}
	defer func() {
		if closeErr := closeResponseBody(res.Body); closeErr != nil {
			e.log.Error(closeErr)
		}
	}()

	if res.IsError() {
		return nil, fmt.Errorf("index error: %s", res.String())
	}

	e.log.Infow(
		fmt.Sprintf("[elasticOrderReadRepository.UpdateOrder] order with id %s updated", order.ID),
		logger.Fields{"Order": order},
	)

	return order, nil
}

// DeleteOrderByID deletes an order by id.
func (e elasticOrderReadRepository) DeleteOrderByID(ctx context.Context, id uuid.UUID) error {
	ctx, span := e.tracer.Start(ctx, "elasticOrderReadRepository.DeleteOrderByID")
	span.SetAttributes(attribute2.String("ID", id.String()))
	defer span.End()

	res, err := e.elasticClient.Delete(
		orderIndex,
		id.String(),
		e.elasticClient.Delete.WithContext(ctx),
		e.elasticClient.Delete.WithRefresh("true"),
	)
	if err != nil {
		return errors.WrapIf(err, "failed to delete order")
	}
	defer func() {
		if closeErr := closeResponseBody(res.Body); closeErr != nil {
			e.log.Error(closeErr)
		}
	}()

	if res.StatusCode == 404 {
		return nil
	}

	if res.IsError() {
		return fmt.Errorf("delete error: %s", res.String())
	}

	e.log.Infow(
		fmt.Sprintf(
			"[elasticOrderReadRepository.DeleteOrderByID] order with id %s deleted",
			id.String(),
		),
		logger.Fields{"ID": id},
	)

	return nil
}
