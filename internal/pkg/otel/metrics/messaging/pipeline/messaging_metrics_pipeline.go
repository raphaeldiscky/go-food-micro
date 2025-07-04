// Package pipelines provides a messaging metrics pipeline.
package pipelines

import (
	"context"
	"fmt"
	"time"

	"github.com/iancoleman/strcase"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"

	"github.com/raphaeldiscky/go-food-micro/internal/pkg/core/messaging/pipeline"
	types2 "github.com/raphaeldiscky/go-food-micro/internal/pkg/core/messaging/types"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/otel/constants/telemetrytags"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/otel/metrics"
	attribute2 "github.com/raphaeldiscky/go-food-micro/internal/pkg/otel/tracing/attribute"
)

// messagingMetricsPipeline is a messaging metrics pipeline.
type messagingMetricsPipeline struct {
	config *config
	meter  metrics.AppMetrics
}

// NewMessagingMetricsPipeline creates a new messaging metrics pipeline.
func NewMessagingMetricsPipeline(
	appMetrics metrics.AppMetrics,
	opts ...Option,
) pipeline.ConsumerPipeline {
	cfg := defaultConfig
	for _, opt := range opts {
		opt.apply(cfg)
	}

	return &messagingMetricsPipeline{
		config: cfg,
		meter:  appMetrics,
	}
}

// Handle handles a request.
func (m *messagingMetricsPipeline) Handle(
	ctx context.Context,
	consumerContext types2.MessageConsumeContext,
	next pipeline.ConsumerHandlerFunc,
) error {
	message := consumerContext.Message()
	messageTypeName := message.GetMessageTypeName()
	snakeTypeName := strcase.ToSnake(messageTypeName)

	successRequestsCounter, err := m.meter.Int64Counter(
		fmt.Sprintf("%s.success_total", snakeTypeName),
		metric.WithUnit("count"),
		metric.WithDescription(
			fmt.Sprintf(
				"Measures the number of success '%s' (%s)",
				snakeTypeName,
				messageTypeName,
			),
		),
	)
	if err != nil {
		return err
	}

	failedRequestsCounter, err := m.meter.Int64Counter(
		fmt.Sprintf("%s.failed_total", snakeTypeName),
		metric.WithUnit("count"),
		metric.WithDescription(
			fmt.Sprintf(
				"Measures the number of failed '%s' (%s)",
				snakeTypeName,
				messageTypeName,
			),
		),
	)
	if err != nil {
		return err
	}

	totalRequestsCounter, err := m.meter.Int64Counter(
		fmt.Sprintf("%s.total", snakeTypeName),
		metric.WithUnit("count"),
		metric.WithDescription(
			fmt.Sprintf(
				"Measures the total number of '%s' (%s)",
				snakeTypeName,
				messageTypeName,
			),
		),
	)
	if err != nil {
		return err
	}

	durationValueRecorder, err := m.meter.Int64Histogram(
		fmt.Sprintf("%s.duration", snakeTypeName),
		metric.WithUnit("ms"),
		metric.WithDescription(
			fmt.Sprintf(
				"Measures the duration of '%s' (%s)",
				snakeTypeName,
				messageTypeName,
			),
		),
	)
	if err != nil {
		return err
	}

	// start recording the duration
	startTime := time.Now()

	err = next(ctx)

	// calculate the duration
	duration := time.Since(startTime).Milliseconds()

	opt := metric.WithAttributes(
		attribute.String(telemetrytags.App.MessageType, messageTypeName),
		attribute.String(telemetrytags.App.MessageName, snakeTypeName),
		attribute2.Object(telemetrytags.App.Message, message),
	)

	// record metrics
	totalRequestsCounter.Add(ctx, 1, opt)

	if err == nil {
		successRequestsCounter.Add(ctx, 1, opt)
	} else {
		failedRequestsCounter.Add(ctx, 1, opt)
	}

	durationValueRecorder.Record(ctx, duration, opt)

	return nil
}
