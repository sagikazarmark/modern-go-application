package otelgokit

import (
	"context"
	"strconv"

	otelcontrib "go.opentelemetry.io/contrib"
	otelglobal "go.opentelemetry.io/otel/api/global"
	"go.opentelemetry.io/otel/api/trace"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/label"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/sd/lb"
)

const (
	tracerName = "go.opentelemetry.io/contrib/instrumentation/github.com/go-kit/kit/otelkit"
)

// defaultSpanName is the default endpoint span name to use.
const defaultSpanName = "gokit/endpoint"

// EndpointMiddleware returns an Endpoint middleware, tracing a Go kit endpoint.
// This endpoint tracer should be used in combination with a Go kit Transport
// tracing middleware, generic OpenCensus transport middleware or custom before
// and after transport functions as service propagation of SpanContext is not
// provided in this middleware.
func EndpointMiddleware(options ...EndpointOption) endpoint.Middleware {
	cfg := &EndpointOptions{}

	for _, o := range options {
		o(cfg)
	}

	if cfg.TracerProvider == nil {
		cfg.TracerProvider = otelglobal.TracerProvider()
	}

	tracer := cfg.TracerProvider.Tracer(
		tracerName,
		trace.WithInstrumentationVersion(otelcontrib.SemVersion()),
	)

	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			operation := cfg.Operation
			if cfg.GetOperation != nil {
				if newOperation := cfg.GetOperation(ctx, operation); newOperation != "" {
					operation = newOperation
				}
			}

			spanName := operation
			if spanName == "" {
				spanName = defaultSpanName
			}

			opts := []trace.SpanOption{
				trace.WithAttributes(cfg.Attributes...),
				trace.WithSpanKind(trace.SpanKindServer),
			}

			if cfg.GetAttributes != nil {
				opts = append(opts, trace.WithAttributes(cfg.GetAttributes(ctx)...))
			}

			ctx, span := tracer.Start(ctx, spanName, opts...)
			defer span.End()

			defer func() {
				if err != nil {
					if lberr, ok := err.(lb.RetryError); ok {
						// handle errors originating from lb.Retry
						for idx, rawErr := range lberr.RawErrors {
							span.AddEvent(
								ctx, "error",
								label.String("error.type", "gokit.lb.retry"),
								label.String("gokit.lb.retry.count", strconv.Itoa(idx+1)),
								label.String("error.message", rawErr.Error()),
							)
						}

						span.RecordError(ctx, lberr.Final, trace.WithErrorStatus(codes.Error))

						return
					}

					// generic error
					span.RecordError(ctx, err, trace.WithErrorStatus(codes.Error))

					return
				}

				// test for business error
				if res, ok := response.(endpoint.Failer); ok && res.Failed() != nil {
					var opts []trace.ErrorOption

					// treating business error as real error in span.
					if !cfg.IgnoreBusinessError {
						opts = append(opts, trace.WithErrorStatus(codes.Error))
					}

					span.RecordError(ctx, res.Failed(), opts...)
					return
				}

				// no errors identified
			}()

			response, err = next(ctx, request)

			return
		}
	}
}
