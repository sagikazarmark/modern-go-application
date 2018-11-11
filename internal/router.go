package internal

import (
	"net/http"

	"github.com/gorilla/mux"
)

type helloWorldController interface {
	HelloWorld(rw http.ResponseWriter, r *http.Request)
	SayHello(rw http.ResponseWriter, r *http.Request)
}

// NewRouter returns a new HTTP handler for the application.
func NewRouter(helloWorld helloWorldController) http.Handler {
	router := mux.NewRouter()

	router.Path("/hello").Methods("GET").HandlerFunc(helloWorld.HelloWorld)
	router.Path("/hello").Methods("POST").HandlerFunc(helloWorld.SayHello)

	return router
}
