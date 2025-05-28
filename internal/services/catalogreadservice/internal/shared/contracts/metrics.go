package contracts

import "go.opentelemetry.io/otel/metric"

// GrpcMetrics contains all gRPC-related metrics
type GrpcMetrics struct {
	CreateProduct  metric.Float64Counter
	UpdateProduct  metric.Float64Counter
	DeleteProduct  metric.Float64Counter
	GetProductByID metric.Float64Counter
	SearchProduct  metric.Float64Counter
}

// RabbitMQMetrics contains all RabbitMQ-related metrics
type RabbitMQMetrics struct {
	CreateProduct metric.Float64Counter
	UpdateProduct metric.Float64Counter
	DeleteProduct metric.Float64Counter
	Success       metric.Float64Counter
	Error         metric.Float64Counter
}
