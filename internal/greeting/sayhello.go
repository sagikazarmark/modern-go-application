package greeting

import (
	"context"
	"fmt"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
)

// SayHello says hello to someone.
type SayHello interface {
	// SayHello says hello to someone.
	SayHello(ctx context.Context, to SayHelloTo, output SayHelloOutput)
}

// SayHelloTo contains who to say hello to.
type SayHelloTo struct {
	Who string
}

// SayHelloOutput is the output channel for saying hello.
type SayHelloOutput interface {
	// Say outputs hello.
	Say(ctx context.Context, hello Hello)
}

// Hello is the common greeting.
type Hello struct {
	Message string
}

type sayHello struct {
	logger log.Logger
}

// NewSayHello returns a new SayHello instance.
func NewSayHello(logger log.Logger) SayHello {
	return &sayHello{
		logger: logger,
	}
}

// SayHello says hello to someone.
func (sh *sayHello) SayHello(ctx context.Context, to SayHelloTo, output SayHelloOutput) {
	level.Info(sh.logger).Log("msg", fmt.Sprintf("Hello, %s!", to.Who))

	hello := Hello{fmt.Sprintf("Hello, %s!", to.Who)}

	output.Say(ctx, hello)
}
