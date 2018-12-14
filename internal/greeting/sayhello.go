package greeting

import (
	"context"
	"fmt"
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
	SaidHelloTo(ctx context.Context, event SaidHelloTo)
}

// SaidHelloTo indicates an event of saying hello to someone.
type SaidHelloTo struct {
	Message string
	Who     string
}

// SayHello says hello to someone.
type SayHello struct {
	events SayHelloEvents
	logger Logger
}

// NewSayHello returns a new SayHello instance.
func NewSayHello(events SayHelloEvents, logger Logger) *SayHello {
	return &SayHello{
		events: events,
		logger: logger,
	}
}

// SayHello says hello to someone.
func (sh *SayHello) SayHello(ctx context.Context, to SayHelloTo, output SayHelloOutput) {
	sh.logger.WithFields(map[string]interface{}{"who": to.Who}).Info("Said hello!")

	hello := Hello{fmt.Sprintf("Hello, %s!", to.Who)}

	output.Say(ctx, hello)

	sh.events.SaidHelloTo(ctx, SaidHelloTo{Message: hello.Message, Who: to.Who})
}
