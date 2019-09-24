package correlation

import (
	"context"
	"math/rand"

	"github.com/go-kit/kit/endpoint"
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

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func generate() string {
	b := make([]byte, 32)

	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}

	return string(b)
}

// Middleware creates a new middleware that ensures the context contains a correlation ID.
func Middleware() endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			cid, ok := ctx.Value(correlationIDContextKey).(string)
			if !ok || cid == "" {
				ctx = context.WithValue(ctx, correlationIDContextKey, generate())
			}

			return next(ctx, request)
		}
	}
}
