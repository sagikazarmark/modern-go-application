package greetingdriver

import (
	"encoding/json"
	"net/http"

	"github.com/goph/emperror"
	"github.com/sagikazarmark/modern-go-application/.gen/openapi/go"
	"github.com/sagikazarmark/modern-go-application/internal/greeting"
)

type helloWorldController struct {
	helloWorld greeting.HelloWorld
	sayHello   greeting.SayHello

	errorHandler emperror.Handler
}

// nolint: golint
func NewHelloWorldController(
	helloWorld greeting.HelloWorld,
	sayHello greeting.SayHello,
	errorHandler emperror.Handler,
) *helloWorldController {
	return &helloWorldController{
		helloWorld:   helloWorld,
		sayHello:     sayHello,
		errorHandler: errorHandler,
	}
}

func (c *helloWorldController) HelloWorld(w http.ResponseWriter, r *http.Request) {
	output := newHelloWorldWebOutput(w, &jsonView{}, "application/json; charset=UTF-8", c.errorHandler)

	c.helloWorld.HelloWorld(r.Context(), output)
}

func (c *helloWorldController) SayHello(w http.ResponseWriter, r *http.Request) {
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

	output := newHelloWorldWebOutput(w, &jsonView{}, "application/json; charset=UTF-8", c.errorHandler)

	c.sayHello.SayHello(r.Context(), sayHelloTo, output)
}
