package opencensus

import (
	"time"

	"contrib.go.opencensus.io/exporter/ocagent"
)

// ExporterConfig configures an OpenCensus exporter.
type ExporterConfig struct {
	Address         string
	Insecure        bool
	ReconnectPeriod time.Duration
}

// Options returns a set of OpenCensus exporter options used for configuring the exporter.
func (c ExporterConfig) Options() []ocagent.ExporterOption {
	options := []ocagent.ExporterOption{
		ocagent.WithAddress(c.Address),
		ocagent.WithReconnectionPeriod(c.ReconnectPeriod),
	}

	if c.Insecure {
		options = append(options, ocagent.WithInsecure())
	}

	return options
}
