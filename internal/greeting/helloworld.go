package greeting

import (
	"context"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
)

// HelloWorldOutput is the output channel for saying hello to the world.
type HelloWorldOutput interface {
	// Say outputs hello.
	Say(ctx context.Context, hello Hello)
}

// HelloWorld outputs Hello World.
type HelloWorld struct {
	logger log.Logger
}

// NewHelloWorld returns a new HelloWorld instance.
func NewHelloWorld(logger log.Logger) *HelloWorld {
	return &HelloWorld{
		logger: logger,
	}
}

// HelloWorld outputs Hello World.
func (hw *HelloWorld) HelloWorld(ctx context.Context, output HelloWorldOutput) {
	level.Info(hw.logger).Log("msg", "Hello, World!")

	hello := Hello{"Hello, World!"}

	output.Say(ctx, hello)
}
