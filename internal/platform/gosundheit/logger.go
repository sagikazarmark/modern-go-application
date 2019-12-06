package gosundheit

import (
	health "github.com/AppsFlyer/go-sundheit"
	"logur.dev/logur"
)

type checkListener struct {
	logger logur.Logger
}

func NewLogger(logger logur.Logger) health.CheckListener {
	return checkListener{
		logger: logger,
	}
}

func (c checkListener) OnCheckStarted(name string) {
	c.logger.Trace("starting check", map[string]interface{}{"check": name})
}

func (c checkListener) OnCheckCompleted(name string, result health.Result) {
	if result.Error != nil {
		c.logger.Trace("check failed", map[string]interface{}{"check": name, "error": result.Error.Error()})

		return
	}

	c.logger.Trace("check completed", map[string]interface{}{"check": name})
}
