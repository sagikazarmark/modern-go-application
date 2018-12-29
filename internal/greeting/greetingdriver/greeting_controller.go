package greetingdriver

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/goph/emperror"
	"github.com/sagikazarmark/modern-go-application/.gen/openapi/go"
	"github.com/sagikazarmark/modern-go-application/internal/greeting"
)

type HelloWorld interface {
	// HelloWorld says hello to the world.
	HelloWorld(ctx context.Context, output greeting.HelloWorldOutput)
}

type SayHello interface {
	// SayHello says hello to someone.
	SayHello(ctx context.Context, to greeting.SayHelloTo, output greeting.SayHelloOutput)
}

// GreetingController collects the greeting use cases and exposes them as HTTP handlers.
type GreetingController struct {
	helloWorld HelloWorld
	sayHello   SayHello

	errorHandler emperror.Handler
}

// NewGreetingController returns a new GreetingController instance.
func NewGreetingController(
	helloWorld HelloWorld,
	sayHello SayHello,
	errorHandler emperror.Handler,
) *GreetingController {
	return &GreetingController{
		helloWorld:   helloWorld,
		sayHello:     sayHello,
		errorHandler: errorHandler,
	}
}

// HelloWorld says hello to the world.
func (c *GreetingController) HelloWorld(w http.ResponseWriter, r *http.Request) {
	output := newGreetingWebOutput(w, &jsonView{}, "application/json; charset=UTF-8", c.errorHandler)

	c.helloWorld.HelloWorld(r.Context(), output)
}

// SayHello says hello to someone.
func (c *GreetingController) SayHello(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)

	var request api.HelloRequest

	if err := decoder.Decode(&request); err != nil {
		c.errorHandler.Handle(err)

		http.Error(w, "invalid request", http.StatusBadRequest)

		return
	}

	sayHelloTo := greeting.SayHelloTo{
		Who: request.Who,
	}

	output := newGreetingWebOutput(w, &jsonView{}, "application/json; charset=UTF-8", c.errorHandler)

	c.sayHello.SayHello(r.Context(), sayHelloTo, output)
}
