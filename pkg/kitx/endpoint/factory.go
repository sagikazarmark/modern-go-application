package endpoint

import (
	"github.com/go-kit/kit/endpoint"
)

// Factory returns an endpoint wrapped with preconfigured middleware.
type Factory interface {
	// NewEndpoint returns an endpoint wrapped with preconfigured middleware.
	// It also accepts an operation name for per operation middleware (eg. logging or tracing middleware).
	NewEndpoint(name string, e endpoint.Endpoint) endpoint.Endpoint
}

// MiddlewareFactory creates a middleware per operation.
type MiddlewareFactory func(name string) endpoint.Middleware

// Middleware wraps singleton middleware and wraps them in a MiddlewareFactory.
func Middleware(middleware endpoint.Middleware) MiddlewareFactory {
	return func(_ string) endpoint.Middleware {
		return middleware
	}
}

// NewFactory returns a new Factory.
func NewFactory(middlewareFactories ...MiddlewareFactory) Factory {
	return factory{
		middlewareFactories: middlewareFactories,
	}
}

type factory struct {
	middlewareFactories []MiddlewareFactory
}

func (f factory) NewEndpoint(name string, e endpoint.Endpoint) endpoint.Endpoint {
	if len(f.middlewareFactories) == 0 {
		return e
	}

	var m endpoint.Middleware

	if len(f.middlewareFactories) == 1 {
		m = f.middlewareFactories[0](name)
	} else {
		mc := make([]endpoint.Middleware, 0, len(f.middlewareFactories))

		for _, mf := range f.middlewareFactories {
			mc = append(mc, mf(name))
		}

		m = endpoint.Chain(mc[0], mc[1:]...)
	}

	return m(e)
}
