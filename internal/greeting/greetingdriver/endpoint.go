package greetingdriver

import (
	"context"

	"github.com/go-kit/kit/endpoint"
	"github.com/sagikazarmark/modern-go-application/internal/greeting"
)

// MakeSayHelloEndpoint constructs a SayHello endpoint wrapping the service.
func MakeSayHelloEndpoint(s Greeter) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		ereq := request.(HelloRequest)

		req := greeting.HelloRequest{
			Name: ereq.Name,
		}

		resp, err := s.SayHello(ctx, req)

		return HelloResponse{
			Reply: resp.Reply,
			Err:   err,
		}, nil
	}
}

// HelloRequest contains a greeting that the service needs to respond to.
type HelloRequest struct {
	Name string
}

// HelloResponse is the the response to a greeting.
type HelloResponse struct {
	Reply string

	Err error
}

func (r HelloResponse) Failed() error {
	return r.Err
}
