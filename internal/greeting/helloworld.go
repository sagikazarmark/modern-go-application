package greeting

import (
	"context"
)

// HelloWorldOutput is the output channel for saying hello to the world.
type HelloWorldOutput interface {
	// Say outputs hello.
	Say(ctx context.Context, hello Hello)
}

// HelloWorld outputs Hello World.
type HelloWorld struct {
	logger Logger
}

// NewHelloWorld returns a new HelloWorld instance.
func NewHelloWorld(logger Logger) *HelloWorld {
	return &HelloWorld{
		logger: logger,
	}
}

// HelloWorld outputs Hello World.
func (hw *HelloWorld) HelloWorld(ctx context.Context, output HelloWorldOutput) {
	hw.logger.Info("Hello, World!")

	hello := Hello{"Hello, World!"}

	output.Say(ctx, hello)
}
