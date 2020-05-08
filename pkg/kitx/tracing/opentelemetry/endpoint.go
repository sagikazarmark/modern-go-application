package opentelemetry

import (
	"context"
	"strconv"

	"go.opentelemetry.io/otel/api/core"
	"go.opentelemetry.io/otel/api/global"
	"go.opentelemetry.io/otel/api/key"
	"go.opentelemetry.io/otel/api/trace"
	"google.golang.org/grpc/codes"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/sd/lb"
)

// TraceEndpointDefaultName is the default endpoint span name to use.
const TraceEndpointDefaultName = "gokit/endpoint"

// TraceEndpoint returns an Endpoint middleware, tracing a Go kit endpoint.
// This endpoint tracer should be used in combination with a Go kit Transport
// tracing middleware, generic OpenTelemetry transport middleware or custom before
// and after transport functions as service propagation of SpanContext is not
// provided in this middleware.
func TraceEndpoint(options ...EndpointOption) endpoint.Middleware {
	cfg := &EndpointOptions{}

	global.Tracer("")

	for _, o := range options {
		o(cfg)
	}

	if cfg.Tracer == nil {
		cfg.Tracer = global.Tracer("")
	}

	if cfg.DefaultName == "" {
		cfg.DefaultName = TraceEndpointDefaultName
	}

	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			name := cfg.DefaultName

			if cfg.GetName != nil {
				if newName := cfg.GetName(ctx, name); newName != "" {
					name = newName
				}
			}

			ctx, span := cfg.Tracer.Start(
				ctx, name,
				trace.WithAttributes(cfg.Attributes...),
				trace.WithAttributes(cfg.getAttributes(ctx)...),
			)
			defer span.End()

			defer func() {
				if err != nil {
					if lberr, ok := err.(lb.RetryError); ok {
						// handle errors originating from lb.Retry
						attrs := make([]core.KeyValue, 0, len(lberr.RawErrors))
						for idx, rawErr := range lberr.RawErrors {
							attrs = append(attrs, key.String("gokit.retry.error."+strconv.Itoa(idx+1), rawErr.Error()))
						}

						span.SetAttributes(attrs...)
						span.SetStatus(codes.Unknown, lberr.Final.Error())

						return
					}

					// generic error
					span.SetStatus(codes.Unknown, err.Error())

					return
				}

				// test for business error
				if res, ok := response.(endpoint.Failer); ok && res.Failed() != nil {
					span.SetAttributes(key.String("gokit.business.error", res.Failed().Error()))

					if cfg.IgnoreBusinessError {
						// status ok

						return
					}

					// treating business error as real error in span.
					span.SetStatus(codes.Unknown, res.Failed().Error())

					return
				}

				// no errors identified
				// status ok
			}()

			response, err = next(ctx, request)

			return
		}
	}
}
