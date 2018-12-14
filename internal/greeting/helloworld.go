package greeting

import (
	"context"
)

// HelloWorldOutput is the output channel for saying hello to the world.
type HelloWorldOutput interface {
	// Say outputs hello.
	Say(ctx context.Context, hello Hello)
}

// HelloWorldEvents is the dispatcher for hello world events.
type HelloWorldEvents interface {
	// SaidHello dispatches a SaidHello event.
	SaidHello(ctx context.Context, event SaidHello)
}

// SaidHello indicates an event of saying hello.
type SaidHello struct {
	Message string
}

// HelloWorld outputs Hello World.
type HelloWorld struct {
	events HelloWorldEvents
	logger Logger
}

// NewHelloWorld returns a new HelloWorld instance.
func NewHelloWorld(events HelloWorldEvents, logger Logger) *HelloWorld {
	return &HelloWorld{
		events: events,
		logger: logger,
	}
}

// HelloWorld outputs Hello World.
func (hw *HelloWorld) HelloWorld(ctx context.Context, output HelloWorldOutput) {
	hw.logger.Info("Hello, World!")

	hello := Hello{"Hello, World!"}

	output.Say(ctx, hello)

	saidHello := SaidHello{Message: hello.Message}

	hw.events.SaidHello(ctx, saidHello)
}
