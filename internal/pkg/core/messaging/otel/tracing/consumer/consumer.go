// Package consumer provides a consumer tracing.
package consumer

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/baggage"
	"go.opentelemetry.io/otel/trace"

	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"

	messageHeader "github.com/raphaeldiscky/go-food-micro/internal/pkg/core/messaging/messageheader"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/core/messaging/otel/tracing"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/core/metadata"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/otel/constants"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/otel/tracing/tracingheaders"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/otel/tracing/utils"
)

// https://devandchill.com/posts/2021/12/go-step-by-step-guide-for-implementing-tracing-on-a-microservices-architecture-2/2/
// https://github.com/open-telemetry/opentelemetry-go-contrib/blob/e84d6d6575e3c3eabcf3204ac88550258673ed3c/instrumentation/github.com/Shopify/sarama/otelsarama/dispatcher.go
// https://opentelemetry.io/docs/reference/specification/trace/semantic_conventions/messaging/
// https://github.com/open-telemetry/opentelemetry-specification/blob/main/specification/trace/semantic_conventions/messaging.md#messaging-attributes
// https://opentelemetry.io/docs/instrumentation/go/manual/#semantic-attributes
// https://trstringer.com/otel-part5-propagation/

// StartConsumerSpan is a function that starts a consumer span.
func StartConsumerSpan(
	ctx context.Context,
	meta *metadata.Metadata,
	payload string,
	consumerTracingOptions *ConsumerTracingOptions,
) (context.Context, trace.Span) {
	ctx = addAfterBaggage(ctx, meta)

	// If there's a span context in the message, use that as the parent context.
	// extracts the tracing from the header and puts it into the context
	carrier := tracing.NewMessageCarrier(meta)
	parentSpanContext := otel.GetTextMapPropagator().Extract(ctx, carrier)

	opts := getTraceOptions(meta, payload, consumerTracingOptions)

	// https://github.com/open-telemetry/opentelemetry-specification/blob/main/specification/trace/semantic_conventions/messaging.md#span-name
	// SpanName = Destination ShortTypeName + Operation ShortTypeName
	ctx, span := tracing.MessagingTracer.Start(
		parentSpanContext,
		fmt.Sprintf("%s %s", consumerTracingOptions.Destination, "receive"),
		opts...)

	span.AddEvent(
		fmt.Sprintf(
			"start consuming message '%s' from the broker",
			messageHeader.GetMessageName(*meta),
		),
	)

	// Emulate Work loads
	time.Sleep(1 * time.Second)

	// we don't want next trace (AfterConsume) becomes child of this span, so we should not use new ctx for (AfterConsume) span. if already exists a span on ctx next span will be a child of that span
	return ctx, span
}

// FinishConsumerSpan is a function that finishes a consumer span.
func FinishConsumerSpan(span trace.Span, err error) error {
	messageName := utils.GetSpanAttribute(span, tracing.MessageName).Value.AsString()

	if err != nil {
		span.AddEvent(fmt.Sprintf("failed to consume message '%s' from the broker", messageName))
		if traceErr := utils.TraceErrStatusFromSpan(span, err); traceErr != nil {
			log.Printf("Failed to trace error status: %v", traceErr)
		}
	}

	span.SetAttributes(
		attribute.Key(constants.SpanId).
			String(span.SpanContext().SpanID().String()), // current span id
	)

	span.AddEvent(fmt.Sprintf("message '%s' consumed from the broker successfully", messageName))
	span.End()

	return err
}

// getTraceOptions is a function that gets the trace options.
func getTraceOptions(
	meta *metadata.Metadata,
	payload string,
	consumerTracingOptions *ConsumerTracingOptions,
) []trace.SpanStartOption {
	correlationId := messageHeader.GetCorrelationId(*meta)

	// https://github.com/open-telemetry/opentelemetry-specification/blob/main/specification/trace/semantic_conventions/messaging.md#topic-with-multiple-consumers
	// https://github.com/open-telemetry/opentelemetry-specification/blob/main/specification/trace/semantic_conventions/messaging.md#batch-receiving
	attrs := []attribute.KeyValue{
		semconv.MessageIDKey.String(messageHeader.GetMessageId(*meta)),
		semconv.MessagingMessageConversationID(correlationId),
		semconv.MessagingOperationReceive,
		attribute.Key(constants.TraceId).String(tracingheaders.GetTracingTraceID(*meta)),
		attribute.Key(constants.Traceparent).String(tracingheaders.GetTracingTraceParent(*meta)),
		attribute.Key(constants.ParentSpanId).String(tracingheaders.GetTracingParentSpanID(*meta)),
		attribute.Key(constants.Timestamp).Int64(time.Now().UnixMilli()),
		attribute.Key(tracing.MessageType).String(messageHeader.GetMessageType(*meta)),
		attribute.Key(tracing.MessageName).String(messageHeader.GetMessageName(*meta)),
		attribute.Key(tracing.Payload).String(payload),
		attribute.String(tracing.Headers, meta.ToJSON()),
		semconv.MessagingDestinationName(consumerTracingOptions.Destination),
		semconv.MessagingSystemKey.String(consumerTracingOptions.MessagingSystem),
	}

	if len(consumerTracingOptions.OtherAttributes) > 0 {
		attrs = append(attrs, consumerTracingOptions.OtherAttributes...)
	}

	opts := []trace.SpanStartOption{
		trace.WithAttributes(attrs...),
		trace.WithSpanKind(trace.SpanKindConsumer),
	}

	return opts
}

// addAfterBaggage is a function that adds after baggage.
func addAfterBaggage(ctx context.Context, meta *metadata.Metadata) context.Context {
	correlationId := messageHeader.GetCorrelationId(*meta)

	correlationIdBag, err := baggage.NewMember(
		string(semconv.MessagingMessageConversationIDKey),
		correlationId,
	)
	if err != nil {
		log.Printf("Failed to create correlationIdBag: %v", err)
	}
	messageIdBag, err := baggage.NewMember(
		string(semconv.MessageIDKey),
		messageHeader.GetMessageId(*meta),
	)
	if err != nil {
		log.Printf("Failed to create messageIdBag: %v", err)
	}
	b, err := baggage.New(correlationIdBag, messageIdBag)
	if err != nil {
		// Log the error but continue with the original context
		// since baggage is optional for tracing
		log.Printf("Failed to create baggage: %v", err)

		return ctx
	}
	ctx = baggage.ContextWithBaggage(ctx, b)

	// new context including baggage
	return ctx
}
