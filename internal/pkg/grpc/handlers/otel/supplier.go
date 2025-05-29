// Package otel provides a otel supplier.
package otel

import (
	"context"

	"go.opentelemetry.io/otel/propagation"
	"google.golang.org/grpc/metadata"
)

// https://github.com/open-telemetry/opentelemetry-go-contrib/blob/main/instrumentation/google.golang.org/grpc/otelgrpc/metadata_supplier.go#L27
// metadataSupplier is a struct that represents a metadata supplier.
type metadataSupplier struct {
	metadata *metadata.MD
}

// var _ propagation.TextMapCarrier = &metadataSupplier{}

// Get is a function that returns the value of a key.
func (s *metadataSupplier) Get(key string) string {
	values := s.metadata.Get(key)
	if len(values) == 0 {
		return ""
	}

	return values[0]
}

// Set is a function that sets the value of a key.
func (s *metadataSupplier) Set(key string, value string) {
	if s.metadata != nil {
		s.metadata.Set(key, value)
	}
}

// Keys is a function that returns the keys of the metadata.
func (s *metadataSupplier) Keys() []string {
	out := make([]string, 0, len(*s.metadata))
	for key := range *s.metadata {
		out = append(out, key)
	}

	return out
}

// extract is a function that extracts the metadata from a context.
func extract(
	ctx context.Context,
	propagators propagation.TextMapPropagator,
) context.Context {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		md = metadata.MD{}
	}

	return propagators.Extract(ctx, &metadataSupplier{
		metadata: &md,
	})
}

// inject is a function that injects the metadata into a context.
func inject(
	ctx context.Context,
	propagators propagation.TextMapPropagator,
) context.Context {
	md, ok := metadata.FromOutgoingContext(ctx)
	if !ok {
		md = metadata.MD{}
	}

	propagators.Inject(ctx, &metadataSupplier{
		metadata: &md,
	})

	return metadata.NewOutgoingContext(ctx, md)
}
