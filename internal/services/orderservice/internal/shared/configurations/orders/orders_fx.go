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
var OrderServiceModule = fx.Module(
	"ordersfx",
	// Shared Modules
	config.NewModule(),
	infrastructure.Module(),

	// Features Modules
	orders.Module(),

	// Other provides
	fx.Provide(configOrdersMetrics),
)

// ref: https://github.com/open-telemetry/opentelemetry-go/blob/main/example/prometheus/main.go

// configGrpcMetrics configures the gRPC metrics.
func configGrpcMetrics(meter metric.Meter, serviceName string) (*contracts.GrpcMetrics, error) {
	if meter == nil {
		return nil, nil
	}

	successGrpcRequests, err := meter.Float64Counter(
		fmt.Sprintf("%s_success_grpc_requests_total", serviceName),
		metric.WithDescription("The total number of success grpc requests"),
	)
	if err != nil {
		return nil, err
	}

	errorGrpcRequests, err := meter.Float64Counter(
		fmt.Sprintf("%s_error_grpc_requests_total", serviceName),
		metric.WithDescription("The total number of error grpc requests"),
	)
	if err != nil {
		return nil, err
	}

	createOrderGrpcRequests, err := meter.Float64Counter(
		fmt.Sprintf("%s_create_order_grpc_requests_total", serviceName),
		metric.WithDescription("The total number of create order grpc requests"),
	)
	if err != nil {
		return nil, err
	}

	updateOrderGrpcRequests, err := meter.Float64Counter(
		fmt.Sprintf("%s_update_order_grpc_requests_total", serviceName),
		metric.WithDescription("The total number of update order grpc requests"),
	)
	if err != nil {
		return nil, err
	}

	payOrderGrpcRequests, err := meter.Float64Counter(
		fmt.Sprintf("%s_pay_order_grpc_requests_total", serviceName),
		metric.WithDescription("The total number of pay order grpc requests"),
	)
	if err != nil {
		return nil, err
	}

	submitOrderGrpcRequests, err := meter.Float64Counter(
		fmt.Sprintf("%s_submit_order_grpc_requests_total", serviceName),
		metric.WithDescription("The total number of submit order grpc requests"),
	)
	if err != nil {
		return nil, err
	}

	getOrderByIDGrpcRequests, err := meter.Float64Counter(
		fmt.Sprintf("%s_get_order_by_id_grpc_requests_total", serviceName),
		metric.WithDescription("The total number of get order by id grpc requests"),
	)
	if err != nil {
		return nil, err
	}

	getOrdersGrpcRequests, err := meter.Float64Counter(
		fmt.Sprintf("%s_get_orders_grpc_requests_total", serviceName),
		metric.WithDescription("The total number of get orders grpc requests"),
	)
	if err != nil {
		return nil, err
	}

	searchOrderGrpcRequests, err := meter.Float64Counter(
		fmt.Sprintf("%s_search_order_grpc_requests_total", serviceName),
		metric.WithDescription("The total number of search order grpc requests"),
	)
	if err != nil {
		return nil, err
	}

	return &contracts.GrpcMetrics{
		SuccessGrpcRequests:      successGrpcRequests,
		ErrorGrpcRequests:        errorGrpcRequests,
		CreateOrderGrpcRequests:  createOrderGrpcRequests,
		UpdateOrderGrpcRequests:  updateOrderGrpcRequests,
		PayOrderGrpcRequests:     payOrderGrpcRequests,
		SubmitOrderGrpcRequests:  submitOrderGrpcRequests,
		GetOrderByIDGrpcRequests: getOrderByIDGrpcRequests,
		GetOrdersGrpcRequests:    getOrdersGrpcRequests,
		SearchOrderGrpcRequests:  searchOrderGrpcRequests,
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
