package greeting_test

import (
	"context"
	"testing"

	"github.com/goph/emperror"
	. "github.com/sagikazarmark/modern-go-application/internal/greeting"
	"github.com/sagikazarmark/modern-go-application/internal/greeting/greetingadapter"
	"github.com/stretchr/testify/assert"
)

type sayHelloOutputStub struct {
	hello Hello
}

func (o *sayHelloOutputStub) Say(ctx context.Context, hello Hello) {
	o.hello = hello
}

type sayHelloEventsStub struct {
	saidHelloTo SaidHelloTo
}

func (e *sayHelloEventsStub) SaidHelloTo(ctx context.Context, event SaidHelloTo) error {
	e.saidHelloTo = event

	return nil
}

func TestSayHello_SayHello(t *testing.T) {
	events := &sayHelloEventsStub{}

	sayHello := NewSayHello(events, greetingadapter.NewNoopLogger(), emperror.NewNoopHandler())

	to := SayHelloTo{Who: "me"}
	output := &sayHelloOutputStub{}

	sayHello.SayHello(context.Background(), to, output)

	assert.Equal(t, Hello{Message: "Hello, me!"}, output.hello)
	assert.Equal(t, SaidHelloTo{Message: "Hello, me!", Who: to.Who}, events.saidHelloTo)
}
