package appkit

import (
	"context"
	"time"

	"github.com/go-kit/kit/endpoint"
	kitxendpoint "github.com/sagikazarmark/kitx/endpoint"

	"github.com/sagikazarmark/modern-go-application/internal/common"
)

// EndpointLoggerFactory logs trace information about a request.
func EndpointLoggerFactory(logger common.Logger) kitxendpoint.MiddlewareFactory {
	return func(name string) endpoint.Middleware {
		return EndpointLogger(logger.WithFields(map[string]interface{}{"operation": name}))
	}
}

// EndpointLogger logs trace information about a request.
func EndpointLogger(logger common.Logger) endpoint.Middleware {
	return func(e endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (interface{}, error) {
			logger := logger.WithContext(ctx)

			logger.Trace("processing request")

			defer func(begin time.Time) {
				logger.Trace("processing request finished", map[string]interface{}{
					"took": time.Since(begin),
				})
			}(time.Now())

			return e(ctx, request)
		}
	}
}
