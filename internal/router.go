package internal

import (
	"net/http"

	"github.com/gorilla/mux"
)

type helloWorldDriver interface {
	HelloWorld(rw http.ResponseWriter, r *http.Request)
}

// NewRouter returns a new HTTP handler for the application.
func NewRouter(helloWorld helloWorldDriver) http.Handler {
	router := mux.NewRouter()

	router.Path("/helloworld").HandlerFunc(helloWorld.HelloWorld)

	return router
}
