package opencensus

import (
	"strings"

	"go.opencensus.io/trace"
)

// TraceConfig configures OpenCensus tracing.
type TraceConfig struct {
	// Sampling describes the default sampler used when creating new spans.
	Sampling SamplingTraceConfig

	// MaxAnnotationEventsPerSpan is max number of annotation events per span.
	MaxAnnotationEventsPerSpan int

	// MaxMessageEventsPerSpan is max number of message events per span.
	MaxMessageEventsPerSpan int

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

// Config returns an OpenCensus trace configuration.
func (t TraceConfig) Config() trace.Config {
	config := trace.Config{
		MaxAnnotationEventsPerSpan: t.MaxAnnotationEventsPerSpan,
		MaxMessageEventsPerSpan:    t.MaxMessageEventsPerSpan,
		MaxAttributesPerSpan:       t.MaxAttributesPerSpan,
		MaxLinksPerSpan:            t.MaxLinksPerSpan,
	}

	switch strings.ToLower(strings.TrimSpace(t.Sampling.Sampler)) {
	case "always":
		config.DefaultSampler = trace.AlwaysSample()

	case "never":
		config.DefaultSampler = trace.NeverSample()

	case "probability":
		config.DefaultSampler = trace.ProbabilitySampler(t.Sampling.Fraction)
	}

	return config
}
