package greeting

import (
	"context"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
)

// HelloWorld says hello to the world.
type HelloWorld interface {
	// HelloWorld says hello to the world.
	HelloWorld(ctx context.Context, output HelloWorldOutput)
}

// HelloWorldOutput is the output channel for saying hello to the world.
type HelloWorldOutput interface { // nolint: golint
	// Say outputs hello.
	Say(ctx context.Context, hello Hello)
}

type helloWorld struct {
	logger log.Logger
}

// NewHelloWorld returns a new HelloWorld instance.
func NewHelloWorld(logger log.Logger) HelloWorld {
	return &helloWorld{
		logger: logger,
	}
}

// HelloWorld outputs Hello World.
func (hw *helloWorld) HelloWorld(ctx context.Context, output HelloWorldOutput) {
	level.Info(hw.logger).Log("msg", "Hello, World!")

	hello := Hello{"Hello, World!"}

	output.Say(ctx, hello)
}
