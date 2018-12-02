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

// SayHello says hello to someone.
type SayHello struct {
	logger Logger
}

// NewSayHello returns a new SayHello instance.
func NewSayHello(logger Logger) *SayHello {
	return &SayHello{
		logger: logger,
	}
}

// SayHello says hello to someone.
func (sh *SayHello) SayHello(ctx context.Context, to SayHelloTo, output SayHelloOutput) {
	sh.logger.WithFields(map[string]interface{}{"who": to.Who}).Info("Said hello!")

	hello := Hello{fmt.Sprintf("Hello, %s!", to.Who)}

	output.Say(ctx, hello)
}
