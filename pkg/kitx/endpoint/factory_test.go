package endpoint

import (
	"context"
	"fmt"
	"testing"

	"github.com/go-kit/kit/endpoint"
)

func TestNewFactory(t *testing.T) {
	tests := []struct {
		mf []MiddlewareFactory
	}{
		{},
		{
			mf: []MiddlewareFactory{
				func(name string) endpoint.Middleware {
					return func(i endpoint.Endpoint) endpoint.Endpoint {
						return func(ctx context.Context, request interface{}) (response interface{}, err error) {
							return i(ctx, request)
						}
					}
				},
			},
		},
		{
			mf: []MiddlewareFactory{
				func(name string) endpoint.Middleware {
					return func(i endpoint.Endpoint) endpoint.Endpoint {
						return func(ctx context.Context, request interface{}) (response interface{}, err error) {
							return i(ctx, request)
						}
					}
				},
				func(name string) endpoint.Middleware {
					return func(i endpoint.Endpoint) endpoint.Endpoint {
						return func(ctx context.Context, request interface{}) (response interface{}, err error) {
							return i(ctx, request)
						}
					}
				},
			},
		},
	}

	for _, test := range tests {
		test := test

		t.Run(fmt.Sprintf("%d_factories", len(test.mf)), func(t *testing.T) {
			factory := NewFactory(test.mf...)

			var endpointCalled bool
			ep := func(ctx context.Context, request interface{}) (interface{}, error) {
				endpointCalled = true

				return nil, nil
			}

			wrappedEndpoint := factory.NewEndpoint("endpoint", ep)

			_, _ = wrappedEndpoint(context.Background(), nil)

			if !endpointCalled {
				t.Error("original endpoint is supposed to be wrapped by the factory")
			}
		})
	}
}
