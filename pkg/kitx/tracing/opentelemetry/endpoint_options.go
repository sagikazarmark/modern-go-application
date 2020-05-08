package opentelemetry

import (
	"context"

	"go.opentelemetry.io/otel/api/core"
	"go.opentelemetry.io/otel/api/trace"
)

// EndpointOptions holds the options for tracing an endpoint
type EndpointOptions struct {
	// Tracer (if specified) is used for starting new spans.
	// Falls back to a global tracer with an empty name.
	//
	// See https://github.com/open-telemetry/opentelemetry-specification/blob/master/specification/trace/api.md#obtaining-a-tracer
	Tracer trace.Tracer

	// DefaultName is used as a fallback if GetName is not specified.
	DefaultName string

	// IgnoreBusinessError if set to true will not treat a business error
	// identified through the endpoint.Failer interface as a span error.
	IgnoreBusinessError bool

	// Attributes holds the default attributes which will be set on span
	// creation by our Endpoint middleware.
	Attributes []core.KeyValue

	// GetName is an optional function that can set the span name based on the existing name
	// for the endpoint and information in the context.
	//
	// If the function is nil, or the returned name is empty, the existing name for the endpoint is used.
	GetName func(ctx context.Context, name string) string

	// GetAttributes is an optional function that can extract trace attributes
	// from the context and add them to the span.
	GetAttributes func(ctx context.Context) []core.KeyValue
}

func (o EndpointOptions) getAttributes(ctx context.Context) []core.KeyValue {
	if o.GetAttributes == nil {
		return nil
	}

	return o.GetAttributes(ctx)
}

// EndpointOption allows for functional options to our OpenTelemetry endpoint
// tracing middleware.
type EndpointOption func(*EndpointOptions)

// WithEndpointConfig sets all configuration options at once by use of the
// EndpointOptions struct.
func WithEndpointConfig(options EndpointOptions) EndpointOption {
	return func(o *EndpointOptions) {
		*o = options
	}
}

// WithTracer sets the tracer.
func WithTracer(tracer trace.Tracer) EndpointOption {
	return func(o *EndpointOptions) {
		o.Tracer = tracer
	}
}

// WithDefaultName sets the default name.
func WithDefaultName(defaultName string) EndpointOption {
	return func(o *EndpointOptions) {
		o.DefaultName = defaultName
	}
}

// WithEndpointAttributes sets the default attributes for the spans created by
// the Endpoint tracer.
func WithEndpointAttributes(attrs ...core.KeyValue) EndpointOption {
	return func(o *EndpointOptions) {
		o.Attributes = attrs
	}
}

// WithIgnoreBusinessError if set to true will not treat a business error
// identified through the endpoint.Failer interface as a span error.
func WithIgnoreBusinessError(val bool) EndpointOption {
	return func(o *EndpointOptions) {
		o.IgnoreBusinessError = val
	}
}

// WithSpanName extracts additional attributes from the request context.
func WithSpanName(fn func(ctx context.Context, name string) string) EndpointOption {
	return func(o *EndpointOptions) {
		o.GetName = fn
	}
}

// WithSpanAttributes extracts additional attributes from the request context.
func WithSpanAttributes(fn func(ctx context.Context) []core.KeyValue) EndpointOption {
	return func(o *EndpointOptions) {
		o.GetAttributes = fn
	}
}
