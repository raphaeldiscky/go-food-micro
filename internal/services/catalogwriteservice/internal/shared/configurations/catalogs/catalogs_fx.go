package catalogs

import (
	"fmt"

	"go.opentelemetry.io/otel/metric"
	"go.uber.org/fx"

	"github.com/raphaeldiscky/go-food-micro/internal/services/catalogwriteservice/config"
	"github.com/raphaeldiscky/go-food-micro/internal/services/catalogwriteservice/internal/products"
	"github.com/raphaeldiscky/go-food-micro/internal/services/catalogwriteservice/internal/shared/configurations/catalogs/infrastructure"
	"github.com/raphaeldiscky/go-food-micro/internal/services/catalogwriteservice/internal/shared/contracts"
	"github.com/raphaeldiscky/go-food-micro/internal/services/catalogwriteservice/internal/shared/data"
)

// https://pmihaylov.com/shared-components-go-microservices/

// NewCatalogsServiceModule is a module that contains the catalogs service module.
func NewCatalogsServiceModule() fx.Option {
	return fx.Module(
		"catalogsfx",
		// Shared Modules
		config.NewModule(),
		infrastructure.NewModule(),
		data.NewModule(),

		// Features Modules
		products.NewModule(),

		// Other provides
		fx.Provide(provideCatalogsMetrics),
	)
}

// metricDefinition holds the metadata for creating a metric counter.
type metricDefinition struct {
	name        string
	description string
}

// metricBuilder helps create metrics with consistent naming and error handling.
type metricBuilder struct {
	meter       metric.Meter
	serviceName string
}

// newMetricBuilder creates a new metricBuilder instance.
func newMetricBuilder(meter metric.Meter, serviceName string) *metricBuilder {
	return &metricBuilder{
		meter:       meter,
		serviceName: serviceName,
	}
}

// createMetrics creates a slice of metric counters from a slice of definitions.
func (b *metricBuilder) createMetrics(
	definitions []metricDefinition,
) ([]metric.Float64Counter, error) {
	counters := make([]metric.Float64Counter, 0, len(definitions))
	for _, def := range definitions {
		counter, err := b.meter.Float64Counter(
			fmt.Sprintf("%s_%s", b.serviceName, def.name),
			metric.WithDescription(def.description),
		)
		if err != nil {
			return nil, fmt.Errorf("failed to create metric %s: %w", def.name, err)
		}
		counters = append(counters, counter)
	}

	return counters, nil
}

// ref: https://github.com/open-telemetry/opentelemetry-go/blob/main/example/prometheus/main.go
// provideCatalogsMetrics is a function that provides the catalogs metrics.
func provideCatalogsMetrics(
	cfg *config.AppOptions,
	meter metric.Meter,
) (*contracts.CatalogsMetrics, error) {
	if meter == nil {
		return nil, nil
	}

	builder := newMetricBuilder(meter, cfg.ServiceName)

	// Define all metrics upfront
	grpcMetrics := []metricDefinition{
		{
			name:        "create_product_grpc_requests_total",
			description: "The total number of create product grpc requests",
		},
		{
			name:        "update_product_grpc_requests_total",
			description: "The total number of update product grpc requests",
		},
		{
			name:        "delete_product_grpc_requests_total",
			description: "The total number of delete product grpc requests",
		},
		{
			name:        "get_product_by_id_grpc_requests_total",
			description: "The total number of get product by id grpc requests",
		},
		{
			name:        "search_product_grpc_requests_total",
			description: "The total number of search product grpc requests",
		},
	}

	rabbitMQMetrics := []metricDefinition{
		{
			name:        "create_product_rabbitmq_messages_total",
			description: "The total number of create product rabbirmq messages",
		},
		{
			name:        "update_product_rabbitmq_messages_total",
			description: "The total number of update product rabbirmq messages",
		},
		{
			name:        "delete_product_rabbitmq_messages_total",
			description: "The total number of delete product rabbirmq messages",
		},
		{
			name:        "search_product_rabbitmq_messages_total",
			description: "The total number of success rabbitmq processed messages",
		},
		{
			name:        "error_rabbitmq_processed_messages_total",
			description: "The total number of error rabbitmq processed messages",
		},
	}

	// Create all metrics in batches
	grpcCounters, err := builder.createMetrics(grpcMetrics)
	if err != nil {
		return nil, fmt.Errorf("failed to create gRPC metrics: %w", err)
	}

	rabbitMQCounters, err := builder.createMetrics(rabbitMQMetrics)
	if err != nil {
		return nil, fmt.Errorf("failed to create RabbitMQ metrics: %w", err)
	}

	// Map the counters to their respective fields
	return &contracts.CatalogsMetrics{
		CreateProductGrpcRequests:     grpcCounters[0],
		UpdateProductGrpcRequests:     grpcCounters[1],
		DeleteProductGrpcRequests:     grpcCounters[2],
		GetProductByIDGrpcRequests:    grpcCounters[3],
		SearchProductGrpcRequests:     grpcCounters[4],
		CreateProductRabbitMQMessages: rabbitMQCounters[0],
		UpdateProductRabbitMQMessages: rabbitMQCounters[1],
		DeleteProductRabbitMQMessages: rabbitMQCounters[2],
		SuccessRabbitMQMessages:       rabbitMQCounters[3],
		ErrorRabbitMQMessages:         rabbitMQCounters[4],
	}, nil
}
