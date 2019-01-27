package landingdriver

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sagikazarmark/modern-go-application/internal/landing"
)

// NewHTTPHandler returns a new HTTP handler for the landing page.
func NewHTTPHandler() http.Handler {
	router := mux.NewRouter()

	router.Path("/").Methods("GET").HandlerFunc(Landing)

	return router
}

// Landing is the landing page for Modern Go Application.
func Landing(w http.ResponseWriter, _ *http.Request) {
	w.Header().Add("Content-Type", "text/html")

	_, _ = w.Write([]byte(landing.Template))
}
