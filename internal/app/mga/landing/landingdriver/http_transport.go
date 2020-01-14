package landingdriver

import (
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/markbates/pkger"
)

// RegisterHTTPHandlers mounts the HTTP handler for the landing page in a router.
func RegisterHTTPHandlers(router *mux.Router) {
	router.Path("/").Methods(http.MethodGet).Handler(Landing())
}

// Landing is the landing page for Modern Go Application.
func Landing() http.Handler {
	file, err := pkger.Open("/static/templates/landing.html")
	if err != nil {
		panic(err)
	}

	body, err := ioutil.ReadAll(file)
	if err != nil {
		panic(err)
	}

	return http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Add("Content-Type", "text/html")

		_, _ = w.Write(body)
	})
}
