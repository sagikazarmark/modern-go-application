package greetingdriver

import (
	"context"

	"github.com/sagikazarmark/modern-go-application/internal/greeting"
)

// HelloService responds to greetings
type HelloService interface {
	// SayHello says hello to someone.
	SayHello(ctx context.Context, req greeting.HelloRequest) (*greeting.HelloResponse, error)
}
