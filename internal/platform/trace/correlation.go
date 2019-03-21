package trace

import (
	"context"
)

// nolint: gochecknoglobals
var correlationID = contextKey("correlation-id")

// WithCorrelationID returns a new context annotated with a correlation ID.
func WithCorrelationID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, correlationID, id)
}

// CorrelationID is awesome.
func CorrelationID(ctx context.Context) (string, bool) {
	id, ok := ctx.Value(correlationID).(string)

	return id, ok
}
