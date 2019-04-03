package correlation

import (
	"context"
)

// nolint: gochecknoglobals
var correlationID = contextKey("correlation-id")

// WithID returns a new context annotated with a correlation ID.
func WithID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, correlationID, id)
}

// ID is awesome.
func ID(ctx context.Context) (string, bool) {
	id, ok := ctx.Value(correlationID).(string)

	return id, ok
}
