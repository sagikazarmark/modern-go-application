package log

import (
	stdlog "log"

	kitlog "github.com/go-kit/kit/log"
)

// NewStandardLogger returns a standard library logger.
func NewStandardLogger(logger kitlog.Logger) *stdlog.Logger {
	return stdlog.New(kitlog.NewStdlibAdapter(logger), "", 0)
}
