package grpc

import (
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/transport/grpc"
)

// ServerFactory constructs a new server, which wraps the provided endpoint and implements the grpc.Handler interface.
type ServerFactory interface {
	// NewServer constructs a new server, which wraps the provided endpoint and implements the grpc.Handler interface.
	NewServer(
		e endpoint.Endpoint,
		dec grpc.DecodeRequestFunc,
		enc grpc.EncodeResponseFunc,
		options ...grpc.ServerOption,
	) *grpc.Server
}

// NewServerFactory returns a new ServerFactory.
func NewServerFactory(options ...grpc.ServerOption) ServerFactory {
	return serverFactory{
		options: options,
	}
}

type serverFactory struct {
	options []grpc.ServerOption
}

func (f serverFactory) NewServer(
	e endpoint.Endpoint,
	dec grpc.DecodeRequestFunc,
	enc grpc.EncodeResponseFunc,
	options ...grpc.ServerOption,
) *grpc.Server {
	return grpc.NewServer(
		e,
		dec,
		enc,
		append(f.options, options...)...,
	)
}
