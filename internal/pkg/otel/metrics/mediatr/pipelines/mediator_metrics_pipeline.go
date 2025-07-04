// Package pipelines provides a mediator metrics pipeline.
package pipelines

import (
	"context"
	"fmt"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"

	mediatr "github.com/mehdihadeli/go-mediatr"

	"github.com/raphaeldiscky/go-food-micro/internal/pkg/core/cqrs"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/core/events"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/otel/constants/telemetrytags"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/otel/metrics"
	customAttribute "github.com/raphaeldiscky/go-food-micro/internal/pkg/otel/tracing/attribute"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/reflection/typemapper"
)

// mediatorMetricsPipeline is a mediator metrics pipeline.
type mediatorMetricsPipeline struct {
	config *config
	meter  metrics.AppMetrics
}

// NewMediatorMetricsPipeline creates a new mediator metrics pipeline.
func NewMediatorMetricsPipeline(
	appMetrics metrics.AppMetrics,
	opts ...Option,
) mediatr.PipelineBehavior {
	cfg := defaultConfig
	for _, opt := range opts {
		opt.apply(cfg)
	}

	return &mediatorMetricsPipeline{
		config: cfg,
		meter:  appMetrics,
	}
}

// Handle handles a request.
func (r *mediatorMetricsPipeline) Handle(
	ctx context.Context,
	request interface{},
	next mediatr.RequestHandlerFunc,
) (interface{}, error) {
	payloadSnakeTypeName := typemapper.GetSnakeTypeName(request)
	typeName := typemapper.GetTypeName(request)

	nameTag := telemetrytags.App.RequestName
	typeNameTag := telemetrytags.App.RequestType
	payloadTag := telemetrytags.App.Request
	resultSnakeTypeNameTag := telemetrytags.App.RequestResultName
	resultTag := telemetrytags.App.RequestResult

	if cqrs.IsCommand(request) {
		nameTag = telemetrytags.App.CommandName
		typeNameTag = telemetrytags.App.CommandType
		payloadTag = telemetrytags.App.Command
		resultSnakeTypeNameTag = telemetrytags.App.CommandResultName
		resultTag = telemetrytags.App.CommandResult
	} else if cqrs.IsQuery(request) {
		nameTag = telemetrytags.App.QueryName
		typeNameTag = telemetrytags.App.QueryType
		payloadTag = telemetrytags.App.Query
		resultSnakeTypeNameTag = telemetrytags.App.QueryResultName
		resultTag = telemetrytags.App.QueryResult
	} else if events.IsEvent(request) {
		nameTag = telemetrytags.App.EventName
		typeNameTag = telemetrytags.App.EventType
		payloadTag = telemetrytags.App.Event
		resultSnakeTypeNameTag = telemetrytags.App.EventResultName
		resultTag = telemetrytags.App.EventResult
	}

	successRequestsCounter, err := r.meter.Int64Counter(
		fmt.Sprintf("%s.success_total", payloadSnakeTypeName),
		metric.WithUnit("count"),
		metric.WithDescription(
			fmt.Sprintf(
				"Measures the number of success '%s' (%s)",
				payloadSnakeTypeName,
				typeName,
			),
		),
	)
	if err != nil {
		return nil, err
	}

	failedRequestsCounter, err := r.meter.Int64Counter(
		fmt.Sprintf("%s.failed_total", payloadSnakeTypeName),
		metric.WithUnit("count"),
		metric.WithDescription(
			fmt.Sprintf(
				"Measures the number of failed '%s' (%s)",
				payloadSnakeTypeName,
				typeName,
			),
		),
	)
	if err != nil {
		return nil, err
	}

	totalRequestsCounter, err := r.meter.Int64Counter(
		fmt.Sprintf("%s.total", payloadSnakeTypeName),
		metric.WithUnit("count"),
		metric.WithDescription(
			fmt.Sprintf(
				"Measures the total number of '%s' (%s)",
				payloadSnakeTypeName,
				typeName,
			),
		),
	)
	if err != nil {
		return nil, err
	}

	durationValueRecorder, err := r.meter.Int64Histogram(
		fmt.Sprintf("%s.duration", payloadSnakeTypeName),
		metric.WithUnit("ms"),
		metric.WithDescription(
			fmt.Sprintf(
				"Measures the duration of '%s' (%s)",
				payloadSnakeTypeName,
				typeName,
			),
		),
	)
	if err != nil {
		return nil, err
	}

	// start recording the duration
	startTime := time.Now()

	response, err := next(ctx)

	// calculate the duration
	duration := time.Since(startTime).Milliseconds()

	// response will be nil if we have an error, so we need to check for that
	responseSnakeName := typemapper.GetSnakeTypeName(response)

	opt := metric.WithAttributes(
		attribute.String(nameTag, payloadSnakeTypeName),
		attribute.String(typeNameTag, typeName),
		customAttribute.Object(payloadTag, request),
		attribute.String(resultSnakeTypeNameTag, responseSnakeName),
		customAttribute.Object(resultTag, response),
	)

	// record metrics
	totalRequestsCounter.Add(ctx, 1, opt)

	if err == nil {
		successRequestsCounter.Add(ctx, 1, opt)
	} else {
		failedRequestsCounter.Add(ctx, 1, opt)
	}

	durationValueRecorder.Record(ctx, duration, opt)

	return response, err
}
