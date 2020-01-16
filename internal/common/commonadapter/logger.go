package commonadapter

import (
	"context"

	"logur.dev/logur"

	"github.com/sagikazarmark/modern-go-application/internal/common"
)

// Logger wraps a logur logger and exposes it under a custom interface.
type Logger struct {
	logur.LoggerFacade

	extractor ContextExtractor
}

// ContextExtractor extracts log fields from a context.
type ContextExtractor func(ctx context.Context) map[string]interface{}

// NewLogger returns a new Logger instance.
func NewLogger(logger logur.LoggerFacade) *Logger {
	return &Logger{
		LoggerFacade: logger,
	}
}

// NewContextAwareLogger returns a new Logger instance that can extract information from a context.
func NewContextAwareLogger(logger logur.LoggerFacade, extractor ContextExtractor) *Logger {
	return &Logger{
		LoggerFacade: logur.WithContextExtractor(logger, logur.ContextExtractor(extractor)),
		extractor:    extractor,
	}
}

// WithFields annotates a logger with key-value pairs.
func (l *Logger) WithFields(fields map[string]interface{}) common.Logger {
	return &Logger{
		LoggerFacade: logur.WithFields(l.LoggerFacade, fields),
		extractor:    l.extractor,
	}
}

// WithContext annotates a logger with a context.
func (l *Logger) WithContext(ctx context.Context) common.Logger {
	if l.extractor == nil {
		return l
	}

	return l.WithFields(l.extractor(ctx))
}
