package helloworld

// UseCase outputs Hello World.
type UseCase struct{}

// HelloWorld outputs Hello World.
func (uc *UseCase) HelloWorld() string {
	return "Hello, World!"
}
