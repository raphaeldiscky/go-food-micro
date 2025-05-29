// Package repositories contains the elastic order read repository.
package repositories

import (
	"context"
	"encoding/json"
	"fmt"
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

type elasticOrderReadRepository struct {
	log           logger.Logger
	elasticClient *elasticsearch.Client
	tracer        tracing.AppTracer
}

func NewElasticOrderReadRepository(
	log logger.Logger,
	elasticClient *elasticsearch.Client,
	tracer tracing.AppTracer,
) repositories.OrderElasticRepository {
	return &elasticOrderReadRepository{log: log, elasticClient: elasticClient, tracer: tracer}
}

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
	defer res.Body.Close()

	if res.IsError() {
		return nil, errors.New(fmt.Sprintf("search error: %s", res.String()))
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
	for i, hit := range result.Hits.Hits {
		orders[i] = &hit.Source
	}

	return &utils.ListResult[*readmodels.OrderReadModel]{
		Items:      orders,
		TotalItems: int64(result.Hits.Total.Value),
		Page:       listQuery.Page,
		Size:       listQuery.Size,
		TotalPage:  (result.Hits.Total.Value + listQuery.Size - 1) / listQuery.Size,
	}, nil
}

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
	defer res.Body.Close()

	if res.IsError() {
		return nil, errors.New(fmt.Sprintf("search error: %s", res.String()))
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
	for i, hit := range result.Hits.Hits {
		orders[i] = &hit.Source
	}

	return &utils.ListResult[*readmodels.OrderReadModel]{
		Items:      orders,
		TotalItems: int64(result.Hits.Total.Value),
		Page:       listQuery.Page,
		Size:       listQuery.Size,
		TotalPage:  (result.Hits.Total.Value + listQuery.Size - 1) / listQuery.Size,
	}, nil
}

func (e elasticOrderReadRepository) GetOrderById(
	ctx context.Context,
	id uuid.UUID,
) (*readmodels.OrderReadModel, error) {
	ctx, span := e.tracer.Start(ctx, "elasticOrderReadRepository.GetOrderById")
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
	defer res.Body.Close()

	if res.StatusCode == 404 {
		return nil, nil
	}

	if res.IsError() {
		return nil, errors.New(fmt.Sprintf("get error: %s", res.String()))
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
			"[elasticOrderReadRepository.GetOrderById] order with id %s loaded",
			id.String(),
		),
		logger.Fields{"Order": result.Source, "ID": id},
	)

	return &result.Source, nil
}

func (e elasticOrderReadRepository) GetOrderByOrderId(
	ctx context.Context,
	orderId uuid.UUID,
) (*readmodels.OrderReadModel, error) {
	ctx, span := e.tracer.Start(ctx, "elasticOrderReadRepository.GetOrderByOrderId")
	span.SetAttributes(attribute2.String("OrderId", orderId.String()))
	defer span.End()

	// Build the search query
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"term": map[string]interface{}{
				"orderId": orderId.String(),
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
	defer res.Body.Close()

	if res.IsError() {
		return nil, errors.New(fmt.Sprintf("search error: %s", res.String()))
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
			"[elasticOrderReadRepository.GetOrderByOrderId] order with orderId %s loaded",
			orderId.String(),
		),
		logger.Fields{"Order": order, "OrderId": orderId},
	)

	return &order, nil
}

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
	defer res.Body.Close()

	if res.IsError() {
		return nil, errors.New(fmt.Sprintf("index error: %s", res.String()))
	}

	e.log.Infow(
		fmt.Sprintf("[elasticOrderReadRepository.CreateOrder] order with id %s created", order.ID),
		logger.Fields{"Order": order},
	)

	return order, nil
}

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
	defer res.Body.Close()

	if res.IsError() {
		return nil, errors.New(fmt.Sprintf("index error: %s", res.String()))
	}

	e.log.Infow(
		fmt.Sprintf("[elasticOrderReadRepository.UpdateOrder] order with id %s updated", order.ID),
		logger.Fields{"Order": order},
	)

	return order, nil
}

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
	defer res.Body.Close()

	if res.StatusCode == 404 {
		return nil
	}

	if res.IsError() {
		return errors.New(fmt.Sprintf("delete error: %s", res.String()))
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
