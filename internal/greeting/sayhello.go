package greeting

import (
	"context"

	"github.com/pkg/errors"
)

// HelloEvents is the dispatcher for hello events.
type HelloEvents interface {
	// SaidHello dispatches a SaidHello event.
	SaidHello(ctx context.Context, event SaidHello) error
}

// SaidHello indicates that hello was said to someone.
type SaidHello struct {
	Greeting string
	Reply    string
}

// HelloService responds to greetings.
type HelloService struct {
	events       HelloEvents
	logger       Logger
	errorHandler ErrorHandler
}

// NewHelloService returns a new HelloService instance.
func NewHelloService(events HelloEvents, logger Logger, errorHandler ErrorHandler) *HelloService {
	return &HelloService{
		events:       events,
		logger:       logger,
		errorHandler: errorHandler,
	}
}

// HelloRequest contains a greeting that the service needs to respond to.
type HelloRequest struct {
	Greeting string
}

// HelloResponse is the the response to a greeting.
type HelloResponse struct {
	Reply string
}

// SayHello says hello to someone.
func (sh *HelloService) SayHello(ctx context.Context, req HelloRequest) (*HelloResponse, error) {
	sh.logger.Info("Hello!", map[string]interface{}{"greeting": req.Greeting})

	resp := &HelloResponse{
		Reply: "hello",
	}

	err := sh.events.SaidHello(ctx, SaidHello{Greeting: req.Greeting, Reply: resp.Reply})
	if err != nil {
		sh.errorHandler.Handle(errors.WithMessage(err, "failed to dispatch event"))
	}

	return resp, nil
}
