package jaeger

import (
	"github.com/goph/emperror"
	"github.com/pkg/errors"
	"go.opencensus.io/exporter/jaeger"
)

// NewExporter creates a new, configured Jaeger exporter.
func NewExporter(config Config, errorHandler emperror.Handler) (*jaeger.Exporter, error) {
	exporter, err := jaeger.NewExporter(jaeger.Options{
		CollectorEndpoint: config.CollectorEndpoint,
		AgentEndpoint:     config.AgentEndpoint,
		Username:          config.Username,
		Password:          config.Password,
		OnError: emperror.HandlerWith(
			errorHandler,
			"component", "opencensus",
			"exporter", "jaeger",
		).Handle,
		Process: jaeger.Process{
			ServiceName: config.ServiceName,
		},
	})

	return exporter, errors.Wrap(err, "failed to create jaeger exporter")
}
