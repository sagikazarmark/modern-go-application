package greeting

import (
	"context"

	"github.com/pkg/errors"
)

// GreeterEvents is the dispatcher for hello events.
type GreeterEvents interface {
	// SaidHello dispatches a SaidHello event.
	SaidHello(ctx context.Context, event SaidHello) error
}

// SaidHello indicates that hello was said to someone.
type SaidHello struct {
	Name  string
	Reply string
}

// Greeter responds to greetings.
type Greeter struct {
	events       GreeterEvents
	logger       Logger
	errorHandler ErrorHandler
}

// NewGreeter returns a new Greeter instance.
func NewGreeter(events GreeterEvents, logger Logger, errorHandler ErrorHandler) *Greeter {
	return &Greeter{
		events:       events,
		logger:       logger,
		errorHandler: errorHandler,
	}
}

// HelloRequest contains a greeting that the service needs to respond to.
type HelloRequest struct {
	Name string
}

// HelloResponse is the the response to a greeting.
type HelloResponse struct {
	Reply string
}

// SayHello says hello to someone.
func (sh *Greeter) SayHello(ctx context.Context, req HelloRequest) (*HelloResponse, error) {
	sh.logger.Info("Hello!", map[string]interface{}{"greeting": req.Name})

	resp := &HelloResponse{
		Reply: "hello",
	}

	err := sh.events.SaidHello(ctx, SaidHello{Name: req.Name, Reply: resp.Reply})
	if err != nil {
		sh.errorHandler.Handle(errors.WithMessage(err, "failed to dispatch event"))
	}

	return resp, nil
}
