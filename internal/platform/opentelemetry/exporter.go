package opentelemetry

import (
	"time"

	"go.opentelemetry.io/otel/exporters/otlp"
)

// ExporterConfig configures an OpenCensus exporter.
type ExporterConfig struct {
	Address            string
	Insecure           bool
	ReconnectionPeriod time.Duration
}

// Options returns a set of OpenCensus exporter options used for configuring the exporter.
func (c ExporterConfig) Options() []otlp.ExporterOption {
	options := []otlp.ExporterOption{
		otlp.WithAddress(c.Address),
		otlp.WithReconnectionPeriod(c.ReconnectionPeriod),
	}

	if c.Insecure {
		options = append(options, otlp.WithInsecure())
	}

	return options
}
