package greetingdriver

import (
	"context"

	"github.com/goph/emperror"

	greetingpb "github.com/sagikazarmark/modern-go-application/.gen/proto/greeting"
	"github.com/sagikazarmark/modern-go-application/internal/greeting"
)

// GRPCController collects the greeting use cases and exposes them as HTTP handlers.
type GRPCController struct {
	helloService Greeter

	errorHandler emperror.Handler
}

// NewGRPCController returns a new GRPCController instance.
func NewGRPCController(helloService Greeter, errorHandler emperror.Handler) *GRPCController {
	return &GRPCController{
		helloService: helloService,
		errorHandler: errorHandler,
	}
}

// SayHello says hello to someone.
func (c *GRPCController) SayHello(
	ctx context.Context,
	rpcReq *greetingpb.HelloRequest,
) (*greetingpb.HelloResponse, error) {
	req := greeting.HelloRequest{
		Name: rpcReq.GetName(),
	}

	resp, err := c.helloService.SayHello(ctx, req)
	if err != nil {
		return nil, nil
	}

	return &greetingpb.HelloResponse{
		Reply: resp.Reply,
	}, nil
}
