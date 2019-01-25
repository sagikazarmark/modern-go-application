package greetingdriver

import (
	"context"

	"github.com/sagikazarmark/modern-go-application/internal/greeting"
)

// Greeter responds to greetings
type Greeter interface {
	// SayHello says hello to someone.
	SayHello(ctx context.Context, req greeting.HelloRequest) (*greeting.HelloResponse, error)
}
