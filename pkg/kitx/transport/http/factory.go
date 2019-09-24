package http

import (
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/transport/http"
)

// ServerFactory constructs a new server, which implements http.Handler and wraps the provided endpoint.
type ServerFactory interface {
	// NewServer constructs a new server, which implements http.Handler and wraps the provided endpoint.
	NewServer(
		e endpoint.Endpoint,
		dec http.DecodeRequestFunc,
		enc http.EncodeResponseFunc,
		options ...http.ServerOption,
	) *http.Server
}

// NewServerFactory returns a new ServerFactory.
func NewServerFactory(options ...http.ServerOption) ServerFactory {
	return serverFactory{
		options: options,
	}
}

type serverFactory struct {
	options []http.ServerOption
}

func (f serverFactory) NewServer(
	e endpoint.Endpoint,
	dec http.DecodeRequestFunc,
	enc http.EncodeResponseFunc,
	options ...http.ServerOption,
) *http.Server {
	return http.NewServer(
		e,
		dec,
		enc,
		append(f.options, options...)...,
	)
}
