package greetingdriver

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/goph/emperror"

	"github.com/sagikazarmark/modern-go-application/.gen/openapi/go"
	"github.com/sagikazarmark/modern-go-application/internal/greeting"
)

type SayHello interface {
	// SayHello says hello to someone.
	SayHello(ctx context.Context, to greeting.SayHelloTo, output greeting.SayHelloOutput)
}

// GreetingController collects the greeting use cases and exposes them as HTTP handlers.
type GreetingController struct {
	sayHello SayHello

	errorHandler emperror.Handler
}

// NewGreetingController returns a new GreetingController instance.
func NewGreetingController(
	sayHello SayHello,
	errorHandler emperror.Handler,
) *GreetingController {
	return &GreetingController{
		sayHello:     sayHello,
		errorHandler: errorHandler,
	}
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
