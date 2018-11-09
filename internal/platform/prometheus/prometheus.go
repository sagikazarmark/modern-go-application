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
		OnError:   errorHandler.Handle,
	})

	return exporter, errors.Wrap(err, "failed to create prometheus exporter")
}

// Validate checks that the configuration is valid.
func (c Config) Validate() error {
	return nil
}
