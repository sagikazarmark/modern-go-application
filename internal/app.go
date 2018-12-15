package internal

import (
	"net/http"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/goph/emperror"
	"github.com/goph/logur"
	"github.com/gorilla/mux"
	"github.com/sagikazarmark/modern-go-application/internal/greeting"
	"github.com/sagikazarmark/modern-go-application/internal/greeting/greetingadapter"
	"github.com/sagikazarmark/modern-go-application/internal/greeting/greetingdriver"
)

// NewApp returns a new HTTP application.
func NewApp(logger logur.Logger, publisher message.Publisher, errorHandler emperror.Handler) http.Handler {
	helloWorld := greeting.NewHelloWorld(greetingadapter.NewHelloWorldEvents(publisher), greetingadapter.NewLogger(logger), errorHandler)
	sayHello := greeting.NewSayHello(greetingadapter.NewSayHelloEvents(publisher), greetingadapter.NewLogger(logger), errorHandler)
	helloWorldController := greetingdriver.NewGreetingController(helloWorld, sayHello, errorHandler)

	router := mux.NewRouter()

	router.Path("/hello").Methods("GET").HandlerFunc(helloWorldController.HelloWorld)
	router.Path("/hello").Methods("POST").HandlerFunc(helloWorldController.SayHello)

	return router
}
