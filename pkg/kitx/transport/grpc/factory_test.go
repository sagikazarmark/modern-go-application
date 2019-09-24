package grpc

import (
	"context"
	"testing"

	kitgrpc "github.com/go-kit/kit/transport/grpc"
	"google.golang.org/grpc/metadata"
)

func TestNewServerFactory(t *testing.T) {
	var beforeCalled bool
	factory := NewServerFactory(
		kitgrpc.ServerBefore(func(i context.Context, mds metadata.MD) context.Context {
			beforeCalled = true

			return i
		}),
	)

	var endpointCalled bool
	ep := func(ctx context.Context, request interface{}) (interface{}, error) {
		endpointCalled = true

		return nil, nil
	}

	server := factory.NewServer(
		ep,
		func(i context.Context, i2 interface{}) (request interface{}, err error) {
			return nil, nil
		},
		func(i context.Context, i2 interface{}) (response interface{}, err error) {
			return nil, nil
		},
	)

	_, _, _ = server.ServeGRPC(context.Background(), nil)

	if !beforeCalled {
		t.Error("global before function is supposed to be called")
	}

	if !endpointCalled {
		t.Error("endpoint is supposed to be called")
	}
}
