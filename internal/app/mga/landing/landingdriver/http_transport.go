package landingdriver

import (
	"io/fs"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
)

// RegisterHTTPHandlers mounts the HTTP handler for the landing page in a router.
func RegisterHTTPHandlers(router *mux.Router, fsys fs.FS) {
	router.Path("/").Methods(http.MethodGet).Handler(Landing(fsys))
}

// Landing is the landing page for Modern Go Application.
func Landing(fsys fs.FS) http.Handler {
	file, err := fsys.Open("landing.html")
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
