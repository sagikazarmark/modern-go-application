package helloworld

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHelloWorldUseCase_HelloWorld(t *testing.T) {
	uc := new(UseCase)

	assert.Equal(t, "Hello, World!", uc.HelloWorld())
}

func TestHelloWorldUseCase_SayHello(t *testing.T) {
	uc := new(UseCase)

	assert.Equal(t, "Hello, John!", uc.SayHello("John"))
}
