package opentelemetry

import (
	"strings"

	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

// TraceConfig configures OpenTelemetry tracing.
type TraceConfig struct {
	// Sampling describes the default sampler used when creating new spans.
	Sampling SamplingTraceConfig

	// MaxEventsPerSpan is max number of message events per span.
	MaxEventsPerSpan int

	// MaxAnnotationEventsPerSpan is max number of attributes per span.
	MaxAttributesPerSpan int

	// MaxLinksPerSpan is max number of links per span.
	MaxLinksPerSpan int
}

// SamplingTraceConfig configures OpenCensus trace sampling.
type SamplingTraceConfig struct {
	Sampler  string
	Fraction float64
}

// Config returns an OpenTelemetry trace configuration.
func (t TraceConfig) Config() sdktrace.Config {
	config := sdktrace.Config{
		DefaultSampler:       sdktrace.AlwaysSample(),
		MaxEventsPerSpan:     t.MaxEventsPerSpan,
		MaxAttributesPerSpan: t.MaxAttributesPerSpan,
		MaxLinksPerSpan:      t.MaxLinksPerSpan,
	}

	switch strings.ToLower(strings.TrimSpace(t.Sampling.Sampler)) {
	case "always":
		config.DefaultSampler = sdktrace.AlwaysSample()

	case "never":
		config.DefaultSampler = sdktrace.NeverSample()

	case "probability":
		config.DefaultSampler = sdktrace.TraceIDRatioBased(t.Sampling.Fraction)
	}

	return config
}
