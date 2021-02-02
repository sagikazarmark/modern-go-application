package otelgokit

import (
	"context"

	"go.opentelemetry.io/otel/api/trace"
	"go.opentelemetry.io/otel/label"
)

// EndpointOptions holds the options for tracing an endpoint.
type EndpointOptions struct {
	// TracerProvider provides access to instrumentation Tracers.
	TracerProvider trace.TracerProvider

	// IgnoreBusinessError if set to true will not treat a business error
	// identified through the endpoint.Failer interface as a span error.
	IgnoreBusinessError bool

	// Operation identifies the current operation and serves as a span name.
	Operation string

	// GetOperation is an optional function that can set the span name based on the existing operation
	// for the endpoint and information in the context.
	//
	// If the function is nil, or the returned operation is empty, the existing operation for the endpoint is used.
	GetOperation func(ctx context.Context, operation string) string

	// Attributes holds the default attributes for each span created by this middleware.
	Attributes []label.KeyValue

	// GetAttributes is an optional function that can extract trace attributes
	// from the context and add them to the span.
	GetAttributes func(ctx context.Context) []label.KeyValue
}

// EndpointOption allows for functional options to our OpenCensus endpoint
// tracing middleware.
type EndpointOption func(*EndpointOptions)

// WithEndpointOptions sets all configuration options at once.
func WithEndpointOptions(options EndpointOptions) EndpointOption {
	return func(o *EndpointOptions) {
		*o = options
	}
}

// WithTracerProvider specifies a tracer provider to use for creating a tracer.
// If none is specified, the global provider is used.
func WithTracerProvider(provider trace.TracerProvider) EndpointOption {
	return func(o *EndpointOptions) {
		o.TracerProvider = provider
	}
}

// WithIgnoreBusinessError if set to true will not treat a business error
// identified through the endpoint.Failer interface as a span error.
func WithIgnoreBusinessError(val bool) EndpointOption {
	return func(o *EndpointOptions) {
		o.IgnoreBusinessError = val
	}
}

// WithOperation sets an operation name for an endpoint.
// Use this when you register a middleware for each endpoint.
func WithOperation(operation string) EndpointOption {
	return func(o *EndpointOptions) {
		o.Operation = operation
	}
}

// WithOperationGetter sets an operation name getter function in EndpointOptions.
func WithOperationGetter(fn func(ctx context.Context, name string) string) EndpointOption {
	return func(o *EndpointOptions) {
		o.GetOperation = fn
	}
}

// WithAttributes sets the default attributes for the spans created by the Endpoint tracer.
func WithAttributes(attrs ...label.KeyValue) EndpointOption {
	return func(o *EndpointOptions) {
		o.Attributes = attrs
	}
}

// WithAttributeGetter extracts additional attributes from the context.
func WithAttributeGetter(fn func(ctx context.Context) []label.KeyValue) EndpointOption {
	return func(o *EndpointOptions) {
		o.GetAttributes = fn
	}
}
