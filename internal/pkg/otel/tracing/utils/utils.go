// Package utils provides a module for the utils.
package utils

import (
	"context"
	"net/http"
	"reflect"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/semconv/v1.20.0/httpconv"
	"go.opentelemetry.io/otel/trace"

	linq "github.com/ahmetb/go-linq/v3"
	trace2 "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"

	"github.com/raphaeldiscky/go-food-micro/internal/pkg/core/metadata"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/grpc/grpcerrors"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/http/httperrors/customerrors"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/http/httperrors/problemdetails"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/otel/constants/telemetrytags"
	errorUtils "github.com/raphaeldiscky/go-food-micro/internal/pkg/utils/errorutils"
)

// traceContextKeyType is a type for the trace context key.
type traceContextKeyType int

const parentSpanKey traceContextKeyType = iota + 1

// HttpTraceStatusFromSpan creates an error span if we have an error and a successful span when error is nil.
func HttpTraceStatusFromSpan(span trace.Span, err error) error {
	isError := err != nil

	if customerrors.IsCustomError(err) {
		httpError := problemdetails.ParseError(err)

		return HttpTraceStatusFromSpanWithCode(
			span,
			err,
			httpError.GetStatus(),
		)
	}

	var (
		status      int
		code        codes.Code
		description = ""
	)

	if isError {
		status = http.StatusInternalServerError
		code = codes.Error
		description = err.Error()
	} else {
		status = http.StatusOK
		code = codes.Ok
	}

	span.SetStatus(code, description)
	span.SetAttributes(
		semconv.HTTPStatusCode(status),
	)

	if isError {
		stackTraceError := errorUtils.ErrorsWithStack(err)

		// https://opentelemetry.io/docs/instrumentation/go/manual/#record-errors
		span.SetAttributes(
			attribute.String(telemetrytags.Exceptions.Message, err.Error()),
			attribute.String(telemetrytags.Exceptions.Stacktrace, stackTraceError),
		)
		span.RecordError(err)
	}

	return err
}

// TraceStatusFromSpan creates a trace status from a span.
func TraceStatusFromSpan(span trace.Span, err error) error {
	isError := err != nil

	var (
		code        codes.Code
		description = ""
	)

	if isError {
		code = codes.Error
		description = err.Error()
	} else {
		code = codes.Ok
	}

	span.SetStatus(code, description)

	if isError {
		stackTraceError := errorUtils.ErrorsWithStack(err)

		// https://opentelemetry.io/docs/instrumentation/go/manual/#record-errors
		span.SetAttributes(
			attribute.String(telemetrytags.Exceptions.Message, err.Error()),
			attribute.String(telemetrytags.Exceptions.Stacktrace, stackTraceError),
		)
		span.RecordError(err)
	}

	return err
}

// TraceErrStatusFromSpan creates a trace error status from a span.
func TraceErrStatusFromSpan(span trace.Span, err error) error {
	isError := err != nil

	span.SetStatus(codes.Error, err.Error())

	if isError {
		stackTraceError := errorUtils.ErrorsWithStack(err)

		// https://opentelemetry.io/docs/instrumentation/go/manual/#record-errors
		span.SetAttributes(
			attribute.String(telemetrytags.Exceptions.Message, err.Error()),
			attribute.String(telemetrytags.Exceptions.Stacktrace, stackTraceError),
		)
		span.RecordError(err)
	}

	return err
}

// HttpTraceStatusFromSpanWithCode creates an error span with specific status code if we have an error and a successful span when error is nil with a specific status.
func HttpTraceStatusFromSpanWithCode(
	span trace.Span,
	err error,
	code int,
) error {
	if err != nil {
		stackTraceError := errorUtils.ErrorsWithStack(err)

		// https://opentelemetry.io/docs/instrumentation/go/manual/#record-errors
		span.SetAttributes(
			attribute.String(telemetrytags.Exceptions.Message, err.Error()),
			attribute.String(telemetrytags.Exceptions.Stacktrace, stackTraceError),
		)
		span.RecordError(err)
	}

	if code > 0 {
		// httpconv doesn't exist in semconv v1.21.0, and it moved to `opentelemetry-go-contrib` pkg
		// https://github.com/open-telemetry/opentelemetry-go/pull/4362
		// https://github.com/open-telemetry/opentelemetry-go/issues/4081
		// using ClientStatus instead of ServerStatus for consideration of 4xx status as error
		span.SetStatus(httpconv.ClientStatus(code))
		span.SetAttributes(semconv.HTTPStatusCode(code))
	} else {
		span.SetStatus(codes.Error, "")
		span.SetAttributes(semconv.HTTPStatusCode(http.StatusInternalServerError))
	}

	return err
}

// HttpTraceStatusFromContext creates an error span if we have an error and a successful span when error is nil.
func HttpTraceStatusFromContext(ctx context.Context, err error) error {
	// https://opentelemetry.io/docs/instrumentation/go/manual/#record-errors
	span := trace.SpanFromContext(ctx)

	defer span.End()

	return HttpTraceStatusFromSpan(span, err)
}

// TraceStatusFromContext creates a trace status from a context.
func TraceStatusFromContext(ctx context.Context, err error) error {
	// https://opentelemetry.io/docs/instrumentation/go/manual/#record-errors
	span := trace.SpanFromContext(ctx)

	defer span.End()

	return TraceStatusFromSpan(span, err)
}

// TraceErrStatusFromContext creates a trace error status from a context.
func TraceErrStatusFromContext(ctx context.Context, err error) error {
	// https://opentelemetry.io/docs/instrumentation/go/manual/#record-errors
	span := trace.SpanFromContext(ctx)

	defer span.End()

	return TraceErrStatusFromSpan(span, err)
}

// GrpcTraceErrFromSpan sets span with status error with error message.
func GrpcTraceErrFromSpan(span trace.Span, err error) error {
	isError := err != nil

	span.SetStatus(codes.Error, err.Error())

	if isError {
		stackTraceError := errorUtils.ErrorsWithStack(err)
		// https://opentelemetry.io/docs/instrumentation/go/manual/#record-errors
		span.SetAttributes(
			attribute.String(telemetrytags.Exceptions.Message, err.Error()),
			attribute.String(telemetrytags.Exceptions.Stacktrace, stackTraceError),
		)

		if customerrors.IsCustomError(err) {
			grpcErr := grpcerrors.ParseError(err)
			span.SetAttributes(
				semconv.RPCGRPCStatusCodeKey.Int(int(grpcErr.GetStatus())),
			)
		}

		span.RecordError(err)
	}

	return err
}

// GrpcTraceErrFromSpanWithCode sets span with status error with error message.
func GrpcTraceErrFromSpanWithCode(span trace.Span, err error, code int) error {
	isError := err != nil

	span.SetStatus(codes.Error, err.Error())

	if isError {
		stackTraceError := errorUtils.ErrorsWithStack(err)
		// https://opentelemetry.io/docs/instrumentation/go/manual/#record-errors
		span.SetAttributes(
			attribute.String(telemetrytags.Exceptions.Message, err.Error()),
			attribute.String(telemetrytags.Exceptions.Stacktrace, stackTraceError),
		)
		span.SetAttributes(semconv.RPCGRPCStatusCodeKey.Int(code))
		span.RecordError(err)
	}

	return err
}

// GetParentSpanContext gets the parent span context.
func GetParentSpanContext(span trace.Span) trace.SpanContext {
	readWriteSpan, ok := span.(trace2.ReadWriteSpan)
	if !ok {
		return *new(trace.SpanContext)
	}

	return readWriteSpan.Parent()
}

// ContextWithParentSpan creates a context with a parent span.
func ContextWithParentSpan(
	parent context.Context,
	span trace.Span,
) context.Context {
	return context.WithValue(parent, parentSpanKey, span)
}

// ParentSpanFromContext gets the parent span from a context.
func ParentSpanFromContext(ctx context.Context) trace.Span {
	_, nopSpan := trace.NewNoopTracerProvider().Tracer("").Start(ctx, "")
	if ctx == nil {
		return nopSpan
	}

	if span, ok := ctx.Value(parentSpanKey).(trace.Span); ok {
		return span
	}

	return nopSpan
}

// CopyFromParentSpanAttribute copies a parent span attribute to a span.
func CopyFromParentSpanAttribute(
	ctx context.Context,
	span trace.Span,
	attributeName string,
	parentAttributeName string,
) {
	parentAtt := GetParentSpanAttribute(ctx, parentAttributeName)
	if reflect.ValueOf(parentAtt).IsZero() {
		return
	}

	span.SetAttributes(
		attribute.String(attributeName, parentAtt.Value.AsString()),
	)
}

// CopyFromParentSpanAttributeIfNotSet copies a parent span attribute to a span if it is not set.
func CopyFromParentSpanAttributeIfNotSet(
	ctx context.Context,
	span trace.Span,
	attributeName string,
	attributeValue string,
	parentAttributeName string,
) {
	if attributeValue != "" {
		span.SetAttributes(attribute.String(attributeName, attributeValue))

		return
	}
	CopyFromParentSpanAttribute(ctx, span, attributeName, parentAttributeName)
}

// GetParentSpanAttribute gets the parent span attribute.
func GetParentSpanAttribute(
	ctx context.Context,
	parentAttributeName string,
) attribute.KeyValue {
	parentSpan := ParentSpanFromContext(ctx)
	readWriteSpan, ok := parentSpan.(trace2.ReadWriteSpan)
	if !ok {
		return *new(attribute.KeyValue)
	}
	att := linq.From(readWriteSpan.Attributes()).
		FirstWithT(func(att attribute.KeyValue) bool { return string(att.Key) == parentAttributeName })

	return att.(attribute.KeyValue)
}

// GetSpanAttributeFromCurrentContext gets the span attribute from the current context.
func GetSpanAttributeFromCurrentContext(
	ctx context.Context,
	attributeName string,
) attribute.KeyValue {
	span := trace.SpanFromContext(ctx)
	readWriteSpan, ok := span.(trace2.ReadWriteSpan)
	if !ok {
		return *new(attribute.KeyValue)
	}
	att := linq.From(readWriteSpan.Attributes()).
		FirstWithT(func(att attribute.KeyValue) bool { return string(att.Key) == attributeName })

	return att.(attribute.KeyValue)
}

// GetSpanAttribute gets the span attribute.
func GetSpanAttribute(
	span trace.Span,
	attributeName string,
) attribute.KeyValue {
	readWriteSpan, ok := span.(trace2.ReadWriteSpan)
	if !ok {
		return *new(attribute.KeyValue)
	}

	att := linq.From(readWriteSpan.Attributes()).
		FirstWithT(func(att attribute.KeyValue) bool { return string(att.Key) == attributeName })

	return att.(attribute.KeyValue)
}

// MapsToAttributes maps a maps to attributes.
func MapsToAttributes(maps map[string]interface{}) []attribute.KeyValue {
	var att []attribute.KeyValue

	for key, val := range maps {
		switch v := val.(type) {
		case string:
			att = append(att, attribute.String(key, v))
		case int64:
			att = append(att, attribute.Int64(key, v))
		case int:
			att = append(att, attribute.Int(key, v))
		case int32:
			att = append(att, attribute.Int(key, int(v)))
		case float64:
			att = append(att, attribute.Float64(key, v))
		case float32:
			att = append(att, attribute.Float64(key, float64(v)))
		case bool:
			att = append(att, attribute.Bool(key, v))
		}
	}

	return att
}

// MetadataToSet converts a metadata to a set.
func MetadataToSet(meta metadata.Metadata) attribute.Set {
	var keyValue []attribute.KeyValue
	for key, val := range meta {
		keyValue = append(keyValue, attribute.String(key, val.(string)))
	}

	return attribute.NewSet(keyValue...)
}
