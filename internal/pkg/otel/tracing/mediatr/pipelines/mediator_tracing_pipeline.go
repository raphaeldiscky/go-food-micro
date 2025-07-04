// Package pipelines provides a mediator tracing pipeline.
package pipelines

import (
	"context"
	"fmt"
	"strings"

	"go.opentelemetry.io/otel/attribute"

	mediatr "github.com/mehdihadeli/go-mediatr"

	"github.com/raphaeldiscky/go-food-micro/internal/pkg/otel/constants/telemetrytags"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/otel/constants/tracing/components"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/otel/tracing"
	customAttribute "github.com/raphaeldiscky/go-food-micro/internal/pkg/otel/tracing/attribute"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/otel/tracing/utils"
	typeMapper "github.com/raphaeldiscky/go-food-micro/internal/pkg/reflection/typemapper"
)

// mediatorTracingPipeline is a mediator tracing pipeline.
type mediatorTracingPipeline struct {
	config *config
	tracer tracing.AppTracer
}

// NewMediatorTracingPipeline creates a new mediator tracing pipeline.
func NewMediatorTracingPipeline(
	appTracer tracing.AppTracer,
	opts ...Option,
) mediatr.PipelineBehavior {
	cfg := defaultConfig
	for _, opt := range opts {
		opt.apply(cfg)
	}

	return &mediatorTracingPipeline{
		config: cfg,
		tracer: appTracer,
	}
}

// Handle handles a request.
func (r *mediatorTracingPipeline) Handle(
	ctx context.Context,
	request interface{},
	next mediatr.RequestHandlerFunc,
) (interface{}, error) {
	requestName := typeMapper.GetSnakeTypeName(request)

	componentName := components.RequestHandler
	requestNameTag := telemetrytags.App.RequestName
	requestTag := telemetrytags.App.Request
	requestResultNameTag := telemetrytags.App.RequestResultName
	requestResultTag := telemetrytags.App.RequestResult

	switch {
	case strings.Contains(typeMapper.GetPackageName(request), "command") || strings.Contains(typeMapper.GetPackageName(request), "commands"):
		componentName = components.CommandHandler
		requestNameTag = telemetrytags.App.CommandName
		requestTag = telemetrytags.App.Command
		requestResultNameTag = telemetrytags.App.CommandResultName
		requestResultTag = telemetrytags.App.CommandResult
	case strings.Contains(typeMapper.GetPackageName(request), "query") || strings.Contains(typeMapper.GetPackageName(request), "queries"):
		componentName = components.QueryHandler
		requestNameTag = telemetrytags.App.QueryName
		requestTag = telemetrytags.App.Query
		requestResultNameTag = telemetrytags.App.QueryResultName
		requestResultTag = telemetrytags.App.QueryResult
	case strings.Contains(typeMapper.GetPackageName(request), "event") || strings.Contains(typeMapper.GetPackageName(request), "events"):
		componentName = components.EventHandler
		requestNameTag = telemetrytags.App.EventName
		requestTag = telemetrytags.App.Event
		requestResultNameTag = telemetrytags.App.EventResultName
		requestResultTag = telemetrytags.App.EventResult
	}

	operationName := fmt.Sprintf("%s_handler", requestName)
	spanName := fmt.Sprintf(
		"%s.%s/%s",
		componentName,
		operationName,
		requestName,
	) // by convention, we use this format to identify the span name

	// https://golang.org/pkg/context/#Context
	newCtx, span := r.tracer.Start(ctx, spanName)

	defer span.End()

	span.SetAttributes(
		attribute.String(requestNameTag, requestName),
		customAttribute.Object(requestTag, request),
	)

	response, err := next(newCtx)

	responseName := typeMapper.GetSnakeTypeName(response)
	span.SetAttributes(
		attribute.String(requestResultNameTag, responseName),
		customAttribute.Object(requestResultTag, response),
	)

	err = utils.TraceStatusFromSpan(
		span,
		err,
	)

	return response, err
}
