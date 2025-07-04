// Package producer provides a producer tracing.
package producer

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
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/core/messaging/types"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/core/metadata"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/otel/constants"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/otel/tracing/utils"
)

// https://devandchill.com/posts/2021/12/go-step-by-step-guide-for-implementing-tracing-on-a-microservices-architecture-2/2/
// https://github.com/open-telemetry/opentelemetry-go-contrib/blob/v0.12.0/instrumentation/github.com/Shopify/sarama/otelsarama/producer.go
// https://opentelemetry.io/docs/reference/specification/trace/semantic_conventions/messaging/
// https://opentelemetry.io/docs/instrumentation/go/manual/#semantic-attributes
// https://github.com/open-telemetry/opentelemetry-specification/blob/main/specification/trace/semantic_conventions/messaging.md#messaging-attributes
// https://trstringer.com/otel-part5-propagation/

// StartProducerSpan is a function that starts a producer span.
func StartProducerSpan(
	ctx context.Context,
	message types.IMessage,
	meta *metadata.Metadata,
	payload string,
	producerTracingOptions *ProducerTracingOptions,
) (context.Context, trace.Span) {
	ctx = addAfterBaggage(ctx, message, meta)

	// If there's a span context in the message, use that as the parent context.
	// extracts the tracing from the header and puts it into the context
	carrier := tracing.NewMessageCarrier(meta)
	parentSpanContext := otel.GetTextMapPropagator().Extract(ctx, carrier)

	opts := getTraceOptions(meta, message, payload, producerTracingOptions)

	// https://github.com/open-telemetry/opentelemetry-specification/blob/main/specification/trace/semantic_conventions/messaging.md#span-name
	// SpanName = Destination ShortTypeName + Operation ShortTypeName
	ctx, span := tracing.MessagingTracer.Start(
		parentSpanContext,
		fmt.Sprintf("%s %s", producerTracingOptions.Destination, "send"),
		opts...)

	span.AddEvent(
		fmt.Sprintf(
			"start publishing message '%s' to the broker",
			messageHeader.GetMessageName(*meta),
		),
	)

	// Injects current span context, so consumers can use it to propagate span.
	// injects the tracing from the context into the header map
	otel.GetTextMapPropagator().Inject(ctx, carrier)

	// we don't want next trace (AfterProduce) becomes child of this span, so we should not use new ctx for (AfterProducer) span. if already exists a span on ctx next span will be a child of that span
	return ctx, span
}

// FinishProducerSpan is a function that finishes a producer span.
func FinishProducerSpan(span trace.Span, err error) error {
	messageName := utils.GetSpanAttribute(
		span,
		tracing.MessageName,
	).Value.AsString()

	if err != nil {
		span.AddEvent(
			fmt.Sprintf(
				"failed to publish message '%s' to the broker",
				messageName,
			),
		)
		if traceErr := utils.TraceErrStatusFromSpan(span, err); traceErr != nil {
			log.Printf("Failed to trace error status: %v", traceErr)
		}
	}
	span.SetAttributes(
		attribute.Key(constants.TraceId).
			String(span.SpanContext().TraceID().String()),
		attribute.Key(constants.SpanId).
			String(span.SpanContext().SpanID().String()), // current span id
	)

	span.AddEvent(
		fmt.Sprintf(
			"message '%s' published to the broker successfully",
			messageName,
		),
	)
	span.End()

	return err
}

// getTraceOptions is a function that gets the trace options.
func getTraceOptions(
	meta *metadata.Metadata,
	message types.IMessage,
	payload string,
	producerTracingOptions *ProducerTracingOptions,
) []trace.SpanStartOption {
	correlationId := messageHeader.GetCorrelationId(*meta)

	// https://github.com/open-telemetry/opentelemetry-specification/blob/main/specification/trace/semantic_conventions/messaging.md#topic-with-multiple-consumers
	// https://github.com/open-telemetry/opentelemetry-specification/blob/main/specification/trace/semantic_conventions/messaging.md#batch-receiving
	attrs := []attribute.KeyValue{
		semconv.MessageIDKey.String(message.GeMessageId()),
		semconv.MessagingMessageConversationID(correlationId),
		attribute.Key(tracing.MessageType).
			String(message.GetMessageTypeName()),
		attribute.Key(tracing.MessageName).
			String(messageHeader.GetMessageName(*meta)),
		attribute.Key(tracing.Payload).String(payload),
		attribute.String(tracing.Headers, meta.ToJSON()),
		attribute.Key(constants.Timestamp).Int64(time.Now().UnixMilli()),
		semconv.MessagingDestinationName(producerTracingOptions.Destination),
		semconv.MessagingSystemKey.String(
			producerTracingOptions.MessagingSystem,
		),
		semconv.MessagingOperationKey.String("send"),
	}

	if len(producerTracingOptions.OtherAttributes) > 0 {
		attrs = append(attrs, producerTracingOptions.OtherAttributes...)
	}

	opts := []trace.SpanStartOption{
		trace.WithAttributes(attrs...),
		trace.WithSpanKind(trace.SpanKindProducer),
	}

	return opts
}

// addAfterBaggage is a function that adds after baggage.
func addAfterBaggage(
	ctx context.Context,
	message types.IMessage,
	meta *metadata.Metadata,
) context.Context {
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
		message.GeMessageId(),
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
