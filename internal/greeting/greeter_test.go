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

type greeterEventsStub struct {
	saidHello SaidHello
}

func (e *greeterEventsStub) SaidHello(ctx context.Context, event SaidHello) error {
	e.saidHello = event

	return nil
}

func TestGreeter_SayHello(t *testing.T) {
	events := &greeterEventsStub{}

	sayHello := NewGreeter(events, greetingadapter.NewNoopLogger(), emperror.NewNoopHandler())

	req := HelloRequest{Name: "welcome"}

	resp, err := sayHello.SayHello(context.Background(), req)

	require.NoError(t, err)
	assert.Equal(t, &HelloResponse{Reply: "hello"}, resp)
	assert.Equal(t, SaidHello{Name: req.Name, Reply: resp.Reply}, events.saidHello)
}
