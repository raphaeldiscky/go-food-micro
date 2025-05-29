// Package orders contains the orderservice module.
package orders

import (
	"fmt"

	"go.opentelemetry.io/otel/metric"
	"go.uber.org/fx"

	"github.com/raphaeldiscky/go-food-micro/internal/services/orderservice/config"
	"github.com/raphaeldiscky/go-food-micro/internal/services/orderservice/internal/orders"
	"github.com/raphaeldiscky/go-food-micro/internal/services/orderservice/internal/shared/configurations/orders/infrastructure"
	"github.com/raphaeldiscky/go-food-micro/internal/services/orderservice/internal/shared/contracts"
)

// https://pmihaylov.com/shared-components-go-microservices/

// OrderServiceModule is the module for the orderservice.
func OrderServiceModule() fx.Option {
	return fx.Module(
		"ordersfx",
		// Shared Modules
		config.NewModule(),
		infrastructure.Module(),

		// Features Modules
		orders.Module(),

		// Other provides
		fx.Provide(configOrdersMetrics),
	)
}

// ref: https://github.com/open-telemetry/opentelemetry-go/blob/main/example/prometheus/main.go

// createCounter creates a new counter with the given name and description.
func createCounter(meter metric.Meter, name, description string) (metric.Float64Counter, error) {
	if meter == nil {
		return nil, nil
	}

	return meter.Float64Counter(name, metric.WithDescription(description))
}

// configGrpcMetrics configures the gRPC metrics.
func configGrpcMetrics(meter metric.Meter, serviceName string) (*contracts.GrpcMetrics, error) {
	if meter == nil {
		return nil, nil
	}

	counters := map[string]string{
		"success_grpc_requests_total":         "The total number of success grpc requests",
		"error_grpc_requests_total":           "The total number of error grpc requests",
		"create_order_grpc_requests_total":    "The total number of create order grpc requests",
		"update_order_grpc_requests_total":    "The total number of update order grpc requests",
		"pay_order_grpc_requests_total":       "The total number of pay order grpc requests",
		"submit_order_grpc_requests_total":    "The total number of submit order grpc requests",
		"get_order_by_id_grpc_requests_total": "The total number of get order by id grpc requests",
		"get_orders_grpc_requests_total":      "The total number of get orders grpc requests",
		"search_order_grpc_requests_total":    "The total number of search order grpc requests",
	}

	metrics := make(map[string]metric.Float64Counter)
	for name, desc := range counters {
		counter, err := createCounter(meter, fmt.Sprintf("%s_%s", serviceName, name), desc)
		if err != nil {
			return nil, err
		}
		metrics[name] = counter
	}

	return &contracts.GrpcMetrics{
		SuccessGrpcRequests:      metrics["success_grpc_requests_total"],
		ErrorGrpcRequests:        metrics["error_grpc_requests_total"],
		CreateOrderGrpcRequests:  metrics["create_order_grpc_requests_total"],
		UpdateOrderGrpcRequests:  metrics["update_order_grpc_requests_total"],
		PayOrderGrpcRequests:     metrics["pay_order_grpc_requests_total"],
		SubmitOrderGrpcRequests:  metrics["submit_order_grpc_requests_total"],
		GetOrderByIDGrpcRequests: metrics["get_order_by_id_grpc_requests_total"],
		GetOrdersGrpcRequests:    metrics["get_orders_grpc_requests_total"],
		SearchOrderGrpcRequests:  metrics["search_order_grpc_requests_total"],
	}, nil
}

// configHTTPMetrics configures the HTTP metrics.
func configHTTPMetrics(meter metric.Meter, serviceName string) (*contracts.HTTPMetrics, error) {
	if meter == nil {
		return nil, nil
	}

	getOrdersHTTPRequests, err := meter.Float64Counter(
		fmt.Sprintf("%s_get_orders_http_requests_total", serviceName),
		metric.WithDescription("The total number of get orders http requests"),
	)
	if err != nil {
		return nil, err
	}

	createOrderHTTPRequests, err := meter.Float64Counter(
		fmt.Sprintf("%s_create_order_http_requests_total", serviceName),
		metric.WithDescription("The total number of create order http requests"),
	)
	if err != nil {
		return nil, err
	}

	updateOrderHTTPRequests, err := meter.Float64Counter(
		fmt.Sprintf("%s_update_order_http_requests_total", serviceName),
		metric.WithDescription("The total number of update order http requests"),
	)
	if err != nil {
		return nil, err
	}

	payOrderHTTPRequests, err := meter.Float64Counter(
		fmt.Sprintf("%s_pay_order_http_requests_total", serviceName),
		metric.WithDescription("The total number of pay order http requests"),
	)
	if err != nil {
		return nil, err
	}

	submitOrderHTTPRequests, err := meter.Float64Counter(
		fmt.Sprintf("%s_submit_order_http_requests_total", serviceName),
		metric.WithDescription("The total number of submit order http requests"),
	)
	if err != nil {
		return nil, err
	}

	getOrderByIDHTTPRequests, err := meter.Float64Counter(
		fmt.Sprintf("%s_get_order_by_id_http_requests_total", serviceName),
		metric.WithDescription("The total number of get order by id http requests"),
	)
	if err != nil {
		return nil, err
	}

	searchOrderHTTPRequests, err := meter.Float64Counter(
		fmt.Sprintf("%s_search_order_http_requests_total", serviceName),
		metric.WithDescription("The total number of search order http requests"),
	)
	if err != nil {
		return nil, err
	}

	return &contracts.HTTPMetrics{
		GetOrdersHTTPRequests:    getOrdersHTTPRequests,
		CreateOrderHTTPRequests:  createOrderHTTPRequests,
		UpdateOrderHTTPRequests:  updateOrderHTTPRequests,
		PayOrderHTTPRequests:     payOrderHTTPRequests,
		SubmitOrderHTTPRequests:  submitOrderHTTPRequests,
		GetOrderByIDHTTPRequests: getOrderByIDHTTPRequests,
		SearchOrderHTTPRequests:  searchOrderHTTPRequests,
	}, nil
}

// configRabbitMQMetrics configures the RabbitMQ metrics.
func configRabbitMQMetrics(
	meter metric.Meter,
	serviceName string,
) (*contracts.RabbitMQMetrics, error) {
	if meter == nil {
		return nil, nil
	}

	deleteOrderRabbitMQMessages, err := meter.Float64Counter(
		fmt.Sprintf("%s_delete_order_rabbitmq_messages_total", serviceName),
		metric.WithDescription("The total number of delete order rabbirmq messages"),
	)
	if err != nil {
		return nil, err
	}

	createOrderRabbitMQMessages, err := meter.Float64Counter(
		fmt.Sprintf("%s_create_order_rabbitmq_messages_total", serviceName),
		metric.WithDescription("The total number of create order rabbirmq messages"),
	)
	if err != nil {
		return nil, err
	}

	updateOrderRabbitMQMessages, err := meter.Float64Counter(
		fmt.Sprintf("%s_update_order_rabbitmq_messages_total", serviceName),
		metric.WithDescription("The total number of update order rabbirmq messages"),
	)
	if err != nil {
		return nil, err
	}

	return &contracts.RabbitMQMetrics{
		DeleteOrderRabbitMQMessages: deleteOrderRabbitMQMessages,
		CreateOrderRabbitMQMessages: createOrderRabbitMQMessages,
		UpdateOrderRabbitMQMessages: updateOrderRabbitMQMessages,
	}, nil
}

// configOrdersMetrics configures the orderservice metrics.
func configOrdersMetrics(
	cfg *config.Config,
	meter metric.Meter,
) (*contracts.OrdersMetrics, error) {
	if meter == nil {
		return nil, nil
	}

	serviceName := cfg.AppOptions.ServiceName

	grpcMetrics, err := configGrpcMetrics(meter, serviceName)
	if err != nil {
		return nil, err
	}

	httpMetrics, err := configHTTPMetrics(meter, serviceName)
	if err != nil {
		return nil, err
	}

	rabbitMQMetrics, err := configRabbitMQMetrics(meter, serviceName)
	if err != nil {
		return nil, err
	}

	return &contracts.OrdersMetrics{
		GrpcMetrics:     grpcMetrics,
		HTTPMetrics:     httpMetrics,
		RabbitMQMetrics: rabbitMQMetrics,
	}, nil
}
