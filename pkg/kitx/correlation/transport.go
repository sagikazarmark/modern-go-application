package correlation

import (
	"context"
	stdhttp "net/http"

	"github.com/go-kit/kit/transport/grpc"
	"github.com/go-kit/kit/transport/http"
	"google.golang.org/grpc/metadata"
)

// Capital letters are invalid in HTTP/2
const defaultCorrelationHeader = "correlation-id"

// HTTPToContext moves a correlation ID from request header to context (if any).
func HTTPToContext(headers ...string) http.RequestFunc {
	if len(headers) == 0 {
		headers = []string{defaultCorrelationHeader}
	}

	return func(ctx context.Context, r *stdhttp.Request) context.Context {
		for _, header := range headers {
			cid := r.Header.Get(header)
			if cid != "" {
				return context.WithValue(ctx, correlationIDContextKey, cid)
			}
		}

		return ctx
	}
}

// GRPCToContext moves a correlation ID from request header to context (if any).
func GRPCToContext(headers ...string) grpc.ServerRequestFunc {
	if len(headers) == 0 {
		headers = []string{defaultCorrelationHeader}
	}

	return func(ctx context.Context, md metadata.MD) context.Context {
		for _, header := range headers {
			cid, ok := md[header]
			if ok && cid[0] != "" {
				return context.WithValue(ctx, correlationIDContextKey, cid[0])
			}
		}

		return ctx
	}
}
