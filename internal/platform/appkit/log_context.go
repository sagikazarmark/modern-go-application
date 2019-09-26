package appkit

import (
	"context"

	"github.com/sagikazarmark/kitx/correlation"
	"go.opencensus.io/trace"
)

// ContextExtractor extracts values from a context.
type ContextExtractor struct{}

// Extract extracts values from a context.
func (ContextExtractor) Extract(ctx context.Context) map[string]interface{} {
	fields := make(map[string]interface{})

	if correlationID, ok := correlation.FromContext(ctx); ok {
		fields["correlation_id"] = correlationID
	}

	if span := trace.FromContext(ctx); span != nil {
		spanCtx := span.SpanContext()
		fields["trace_id"] = spanCtx.TraceID.String()
		fields["span_id"] = spanCtx.SpanID.String()
	}

	return fields
}
