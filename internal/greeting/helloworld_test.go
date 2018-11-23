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

func TestHelloWorld_HelloWorld(t *testing.T) {
	helloWorld := NewHelloWorld(greetingadapter.NewNopLogger())

	output := &helloWorldOutputStub{}

	helloWorld.HelloWorld(context.Background(), output)

	assert.Equal(t, Hello{Message: "Hello, World!"}, output.hello)
}
