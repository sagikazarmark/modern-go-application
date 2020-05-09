package othttp

import (
	"go.opentelemetry.io/otel/api/key"
	"go.opentelemetry.io/otel/api/metric"
	"go.opentelemetry.io/otel/api/unit"
)

// The following labels are applied to metrics recorded by this package. Host, Path
// and Method are applied to all measures.
var (
	// Host is the value of the HTTP Host header.
	//
	// The value of this label can be controlled by the HTTP client, so you need
	// to watch out for potentially generating high-cardinality labels in your
	// metrics backend if you use this tag in views.
	Host = key.New("http.host")

	// StatusCode is the numeric HTTP response status code,
	// or "error" if a transport error occurred and no status code was read.
	StatusCode = key.New("http.status")

	// Path is the URL path (not including query string) in the request.
	//
	// The value of this tag can be controlled by the HTTP client, so you need
	// to watch out for potentially generating high-cardinality labels in your
	// metrics backend if you use this tag in views.
	Path = key.New("http.path")

	// Method is the HTTP method of the request, capitalized (GET, POST, etc.).
	Method = key.New("http.method")

	// KeyServerRoute is a low cardinality string representing the logical
	// handler of the request. This is usually the pattern registered on the a
	// ServeMux (or similar string).
	KeyServerRoute = key.New("http_server_route")
)

type instruments struct {
	serverRequestCount  metric.Int64Counter
	serverRequestBytes  metric.Int64Measure
	serverResponseBytes metric.Int64Measure
	serverLatency       metric.Float64Measure
}

func newInstruments(meter metric.Meter) instruments {
	return instruments{
		serverRequestCount: metric.Must(meter).NewInt64Counter(
			"opentelemetry.io/http/server/request_count",
			metric.WithDescription("Count of HTTP requests started"),
			metric.WithUnit(unit.Dimensionless),
		),
		serverRequestBytes: metric.Must(meter).NewInt64Measure(
			"opencensus.io/http/server/request_bytes",
			metric.WithDescription("HTTP request body size if set as ContentLength (uncompressed)"),
			metric.WithUnit(unit.Bytes),
		),
		serverResponseBytes: metric.Must(meter).NewInt64Measure(
			"opencensus.io/http/server/response_bytes",
			metric.WithDescription("HTTP response body size (uncompressed)"),
			metric.WithUnit(unit.Bytes),
		),
		serverLatency: metric.Must(meter).NewFloat64Measure(
			"opencensus.io/http/server/latency",
			metric.WithDescription("End-to-end latency"),
			metric.WithUnit(unit.Milliseconds),
		),
	}
}
