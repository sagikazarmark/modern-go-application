package landingdriver

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/sagikazarmark/modern-go-application/internal/app/mga/landing"
)

// RegisterHTTPHandlers mounts the HTTP handler for the landing page in a router.
func RegisterHTTPHandlers(router *mux.Router) {
	router.Path("/").Methods(http.MethodGet).HandlerFunc(Landing)
}

// Landing is the landing page for Modern Go Application.
func Landing(w http.ResponseWriter, _ *http.Request) {
	w.Header().Add("Content-Type", "text/html")

	_, _ = w.Write([]byte(landing.Template))
}
