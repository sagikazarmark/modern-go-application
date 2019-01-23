package greeting_test

import (
	"context"
	"testing"

	"github.com/goph/emperror"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	. "github.com/sagikazarmark/modern-go-application/internal/greeting"
	"github.com/sagikazarmark/modern-go-application/internal/greeting/greetingadapter"
)

type sayHelloEventsStub struct {
	saidHello SaidHello
}

func (e *sayHelloEventsStub) SaidHello(ctx context.Context, event SaidHello) error {
	e.saidHello = event

	return nil
}

func TestHelloService_SayHello(t *testing.T) {
	events := &sayHelloEventsStub{}

	sayHello := NewHelloService(events, greetingadapter.NewNoopLogger(), emperror.NewNoopHandler())

	req := HelloRequest{Greeting: "welcome"}

	resp, err := sayHello.SayHello(context.Background(), req)

	require.NoError(t, err)
	assert.Equal(t, &HelloResponse{Reply: "hello"}, resp)
	assert.Equal(t, SaidHello{Greeting: req.Greeting, Reply: resp.Reply}, events.saidHello)
}
