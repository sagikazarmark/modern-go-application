package greeting

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

type helloWorldOutputStub struct {
	hello Hello
}

func (o *helloWorldOutputStub) Say(ctx context.Context, hello Hello) {
	o.hello = hello
}

func TestHelloWorld_HelloWorld(t *testing.T) {
	helloWorld := NewHelloWorld(NewNopLogger())

	output := &helloWorldOutputStub{}

	helloWorld.HelloWorld(context.Background(), output)

	assert.Equal(t, Hello{"Hello, World!"}, output.hello)
}
