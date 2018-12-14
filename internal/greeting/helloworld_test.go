package greeting_test

import (
	"context"
	"testing"

	. "github.com/sagikazarmark/modern-go-application/internal/greeting"
	"github.com/sagikazarmark/modern-go-application/internal/greeting/greetingadapter"
	"github.com/stretchr/testify/assert"
)

type helloWorldOutputStub struct {
	hello Hello
}

func (o *helloWorldOutputStub) Say(ctx context.Context, hello Hello) {
	o.hello = hello
}

type helloWorldEventsStub struct {
	saidHello SaidHello
}

func (e *helloWorldEventsStub) SaidHello(ctx context.Context, event SaidHello) {
	e.saidHello = event
}

func TestHelloWorld_HelloWorld(t *testing.T) {
	events := &helloWorldEventsStub{}

	helloWorld := NewHelloWorld(events, greetingadapter.NewNopLogger())

	output := &helloWorldOutputStub{}

	helloWorld.HelloWorld(context.Background(), output)

	assert.Equal(t, Hello{Message: "Hello, World!"}, output.hello)
	assert.Equal(t, SaidHello{Message: "Hello, World!"}, events.saidHello)
}
