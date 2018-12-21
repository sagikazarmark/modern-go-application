package greeting

import (
	"context"

	"github.com/pkg/errors"
)

// HelloWorldOutput is the output channel for saying hello to the world.
type HelloWorldOutput interface {
	// Say outputs hello.
	Say(ctx context.Context, hello Hello)
}

// HelloWorldEvents is the dispatcher for hello world events.
type HelloWorldEvents interface {
	// SaidHello dispatches a SaidHello event.
	SaidHello(ctx context.Context, event SaidHello) error
}

// SaidHello indicates that hello was said.
type SaidHello struct {
	Message string
}

// HelloWorld outputs Hello World.
type HelloWorld struct {
	events       HelloWorldEvents
	logger       Logger
	errorHandler ErrorHandler
}

// NewHelloWorld returns a new HelloWorld instance.
func NewHelloWorld(events HelloWorldEvents, logger Logger, errorHandler ErrorHandler) *HelloWorld {
	return &HelloWorld{
		events:       events,
		logger:       logger,
		errorHandler: errorHandler,
	}
}

// HelloWorld outputs Hello World.
func (hw *HelloWorld) HelloWorld(ctx context.Context, output HelloWorldOutput) {
	hw.logger.Info("Hello, World!", nil)

	hello := Hello{"Hello, World!"}

	output.Say(ctx, hello)

	saidHello := SaidHello(hello)

	err := hw.events.SaidHello(ctx, saidHello)
	if err != nil {
		hw.errorHandler.Handle(errors.WithMessage(err, "failed to dispatch event"))
	}
}
