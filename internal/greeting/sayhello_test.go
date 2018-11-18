package greeting

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

type sayHelloOutputStub struct {
	hello Hello
}

func (o *sayHelloOutputStub) Say(ctx context.Context, hello Hello) {
	o.hello = hello
}

func TestSayHello_SayHello(t *testing.T) {
	sayHello := NewSayHello(NewNopLogger())

	to := SayHelloTo{"me"}
	output := &sayHelloOutputStub{}

	sayHello.SayHello(context.Background(), to, output)

	assert.Equal(t, Hello{"Hello, me!"}, output.hello)
}
