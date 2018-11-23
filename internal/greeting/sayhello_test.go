package greeting_test

import (
	"context"
	"testing"

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

func TestSayHello_SayHello(t *testing.T) {
	sayHello := NewSayHello(greetingadapter.NewNopLogger())

	to := SayHelloTo{Who: "me"}
	output := &sayHelloOutputStub{}

	sayHello.SayHello(context.Background(), to, output)

	assert.Equal(t, Hello{Message: "Hello, me!"}, output.hello)
}
