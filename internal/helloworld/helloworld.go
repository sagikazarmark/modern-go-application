package helloworld

import (
	"context"
)

// HelloWorld defines an interface for saying hello to the world.
type HelloWorld interface {
	// HelloWorld says hello to the world.
	HelloWorld(ctx context.Context, output HelloWorldOutput)
}

// HelloWorldOutput is the output channel for the say hello to the world use case.
type HelloWorldOutput interface { // nolint: golint
	// Say outputs hello.
	Say(hello Hello)
}

// HelloWorld outputs Hello World.
func (uc *UseCase) HelloWorld(ctx context.Context, output HelloWorldOutput) {
	hello := Hello{"Hello, World!"}

	output.Say(hello)
}
