package trace

import (
	"net/http"

	"github.com/goph/idgen"
	"github.com/gorilla/mux"
)

// IDGenerator generates a new ID.
type IDGenerator interface {
	// Generate generates a new ID.
	Generate() (string, error)
}

// HTTPCorrelationIDMiddleware attaches a correlation ID to the request.
func HTTPCorrelationIDMiddleware(idgenerator IDGenerator) mux.MiddlewareFunc {
	idgen := idgen.NewGenerator(idgenerator)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			var id string

			if header := req.Header.Get("Correlation-ID"); header != "" {
				id = header
			} else {
				id = idgen.Generate()
			}

			req = req.WithContext(WithCorrelationID(req.Context(), id))

			next.ServeHTTP(w, req)
		})
	}
}
