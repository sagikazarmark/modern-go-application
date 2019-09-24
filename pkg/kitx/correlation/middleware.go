package correlation

import (
	"context"
	"math/rand"

	"github.com/go-kit/kit/endpoint"
)

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
