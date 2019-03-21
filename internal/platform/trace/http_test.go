package trace

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/goph/idgen"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHTTPCorrelationIDMiddleware(t *testing.T) {
	router := mux.NewRouter()

	router.Use(HTTPCorrelationIDMiddleware(nil))

	router.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id, ok := CorrelationID(r.Context())
		require.True(t, ok)
		assert.Equal(t, "id", id)
	}))

	rec := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("Correlation-ID", "id")

	router.ServeHTTP(rec, req)
}

func TestHTTPCorrelationIDMiddleware_Generate(t *testing.T) {
	router := mux.NewRouter()

	router.Use(HTTPCorrelationIDMiddleware(idgen.NewConstantGenerator("id")))

	router.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id, ok := CorrelationID(r.Context())
		require.True(t, ok)
		assert.Equal(t, "id", id)
	}))

	rec := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/", nil)

	router.ServeHTTP(rec, req)
}
