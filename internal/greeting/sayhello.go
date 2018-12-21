package greeting

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
)

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

// SayHelloEvents is the dispatcher for say hello events.
type SayHelloEvents interface {
	// SaidHelloTo dispatches a SaidHelloTo event.
	SaidHelloTo(ctx context.Context, event SaidHelloTo) error
}

// SaidHelloTo indicates that hello was said to someone.
type SaidHelloTo struct {
	Message string
	Who     string
}

// SayHello says hello to someone.
type SayHello struct {
	events       SayHelloEvents
	logger       Logger
	errorHandler ErrorHandler
}

// NewSayHello returns a new SayHello instance.
func NewSayHello(events SayHelloEvents, logger Logger, errorHandler ErrorHandler) *SayHello {
	return &SayHello{
		events:       events,
		logger:       logger,
		errorHandler: errorHandler,
	}
}

// SayHello says hello to someone.
func (sh *SayHello) SayHello(ctx context.Context, to SayHelloTo, output SayHelloOutput) {
	sh.logger.Info("Said hello!", map[string]interface{}{"who": to.Who})

	hello := Hello{fmt.Sprintf("Hello, %s!", to.Who)}

	output.Say(ctx, hello)

	err := sh.events.SaidHelloTo(ctx, SaidHelloTo{Message: hello.Message, Who: to.Who})
	if err != nil {
		sh.errorHandler.Handle(errors.WithMessage(err, "failed to dispatch event"))
	}
}
