// Package contracts contains the orders metrics.
package contracts

import (
	"go.opentelemetry.io/otel/metric"
)

// GrpcMetrics contains the gRPC metrics.
type GrpcMetrics struct {
	SuccessGrpcRequests      metric.Float64Counter
	ErrorGrpcRequests        metric.Float64Counter
	CreateOrderGrpcRequests  metric.Float64Counter
	UpdateOrderGrpcRequests  metric.Float64Counter
	PayOrderGrpcRequests     metric.Float64Counter
	SubmitOrderGrpcRequests  metric.Float64Counter
	GetOrderByIDGrpcRequests metric.Float64Counter
	GetOrdersGrpcRequests    metric.Float64Counter
	SearchOrderGrpcRequests  metric.Float64Counter
}

// HTTPMetrics contains the HTTP metrics.
type HTTPMetrics struct {
	GetOrdersHTTPRequests    metric.Float64Counter
	CreateOrderHTTPRequests  metric.Float64Counter
	UpdateOrderHTTPRequests  metric.Float64Counter
	PayOrderHTTPRequests     metric.Float64Counter
	SubmitOrderHTTPRequests  metric.Float64Counter
	GetOrderByIDHTTPRequests metric.Float64Counter
	SearchOrderHTTPRequests  metric.Float64Counter
}

// RabbitMQMetrics contains the RabbitMQ metrics.
type RabbitMQMetrics struct {
	DeleteOrderRabbitMQMessages metric.Float64Counter
	CreateOrderRabbitMQMessages metric.Float64Counter
	UpdateOrderRabbitMQMessages metric.Float64Counter
}

// OrdersMetrics is the metrics for the orders.
type OrdersMetrics struct {
	GrpcMetrics     *GrpcMetrics
	HTTPMetrics     *HTTPMetrics
	RabbitMQMetrics *RabbitMQMetrics
}
