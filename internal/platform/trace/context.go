package trace

import (
	"context"

	"go.opencensus.io/trace"
)

type contextKey string

func (c contextKey) String() string {
	return "trace context key " + string(c)
}

type ContextExtractor struct{}

func (*ContextExtractor) Extract(ctx context.Context) map[string]interface{} {
	fields := make(map[string]interface{})

	if correlationID, ok := CorrelationID(ctx); ok {
		fields["correlation_id"] = correlationID
	}

	if span := trace.FromContext(ctx); span != nil {
		spanCtx := span.SpanContext()
		fields["trace_id"] = spanCtx.TraceID.String()
		fields["span_id"] = spanCtx.SpanID.String()
	}

	return fields
}
