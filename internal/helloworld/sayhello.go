package helloworld

import (
	"context"
	"fmt"
)

// SayHello defines an interface for saying hello to someone.
type SayHello interface {
	// SayHello says hello to someone.
	SayHello(ctx context.Context, to SayHelloTo, output SayHelloOutput)
}

// SayHelloTo is the input model of the say hello use case.
type SayHelloTo struct {
	Who string
}

// SayHelloOutput is the output channel for the say hello use case.
type SayHelloOutput interface {
	// Say outputs hello.
	Say(hello Hello)
}

// Hello is the common greeting for hello related use cases.
type Hello struct {
	Message string
}

// SayHello says hello to someone.
func (uc *UseCase) SayHello(ctx context.Context, to SayHelloTo, output SayHelloOutput) {
	hello := Hello{fmt.Sprintf("Hello, %s!", to.Who)}

	output.Say(hello)
}
