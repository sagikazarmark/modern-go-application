package helloworld

import "fmt"

// UseCase outputs Hello World.
type UseCase struct{}

// HelloWorld outputs Hello World.
func (uc *UseCase) HelloWorld() string {
	return "Hello, World!"
}

// SayHello says hello to someone.
func (uc *UseCase) SayHello(who string) string {
	return fmt.Sprintf("Hello, %s!", who)
}
