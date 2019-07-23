package httpbin

import (
	"net/http"

	"github.com/mccutchen/go-httpbin/httpbin"
)

// MakeHTTPHandler returns a new HTTP handler serving HTTPBin.
func MakeHTTPHandler(logger Logger) http.Handler {
	return httpbin.New(
		httpbin.WithObserver(func(result httpbin.Result) {
			logger.Info(
				"httpbin call",
				map[string]interface{}{
					"status":      result.Status,
					"method":      result.Method,
					"uri":         result.URI,
					"size_bytes":  result.Size,
					"duration_ms": result.Duration.Seconds() * 1e3, // https://github.com/golang/go/issues/5491#issuecomment-66079585
				},
			)
		}),
	).Handler()
}
