package prometheus

import (
	"github.com/goph/emperror"
	"github.com/pkg/errors"
	"go.opencensus.io/exporter/prometheus"
)

// NewExporter creates a new, configured Prometheus exporter.
func NewExporter(config Config, errorHandler emperror.Handler) (*prometheus.Exporter, error) {
	exporter, err := prometheus.NewExporter(prometheus.Options{
		Namespace: config.Namespace,
		OnError: emperror.HandlerWith(
			errorHandler,
			"component", "opencensus",
			"exporter", "prometheus",
		).Handle,
	})

	return exporter, errors.Wrap(err, "failed to create prometheus exporter")
}
