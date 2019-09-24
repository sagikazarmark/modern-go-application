// Package correlation provides a set of tools to add correlation ID to the context at certain levels
// (transport, endpoint) of the application.
package correlation

import (
	"context"
)

type contextKey string

// correlationIDContextKey holds the key used to store a
// correlation ID in the context.
const correlationIDContextKey contextKey = "CorrelationID"

// FromContext returns the correlation ID from the context (if any).
// Returns false as the second parameter if none is found.
func FromContext(ctx context.Context) (string, bool) {
	id, ok := ctx.Value(correlationIDContextKey).(string)

	return id, ok
}
