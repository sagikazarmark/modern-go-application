package invisionkitlog

import (
	"fmt"

	"github.com/InVisionApp/go-logger"
	kitlog "github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
)

type loggerShim struct {
	logger kitlog.Logger
}

// New returns a new kitlog shim for invision logger interface.
func New(logger kitlog.Logger) *loggerShim {
	return &loggerShim{logger}
}

// this will add a space between all elements in the slice
// this func is needed because fmt.Sprint will not separate
// inputs by a space in all cases, which makes the resulting
// output very hard to read
func spaceSep(a []interface{}) []interface{} {
	aLen := len(a)
	if aLen <= 1 {
		return a
	}

	// we only allocate enough room to add a single space between
	// all elements, so len(a) - 1
	spaceSlice := make([]interface{}, aLen-1)
	// add the empty space to the end of the original slice
	a = append(a, spaceSlice...)

	// stagger the values.  this will leave an empty slot between all
	// values to be filled with a space
	for i := aLen - 1; i > 0; i-- {
		a[i+i] = a[i]
		a[i+i-1] = " "
	}

	return a
}

func (s *loggerShim) Debug(msg ...interface{}) {
	level.Debug(s.logger).Log("msg", fmt.Sprint(spaceSep(msg)...))
}

func (s *loggerShim) Info(msg ...interface{}) {
	level.Info(s.logger).Log("msg", fmt.Sprint(spaceSep(msg)...))
}

func (s *loggerShim) Warn(msg ...interface{}) {
	level.Warn(s.logger).Log("msg", fmt.Sprint(spaceSep(msg)...))
}

func (s *loggerShim) Error(msg ...interface{}) {
	level.Error(s.logger).Log("msg", fmt.Sprint(spaceSep(msg)...))
}

func (s *loggerShim) Debugln(msg ...interface{}) {
	level.Debug(s.logger).Log("msg", fmt.Sprint(spaceSep(msg)...))
}

func (s *loggerShim) Infoln(msg ...interface{}) {
	level.Info(s.logger).Log("msg", fmt.Sprint(spaceSep(msg)...))
}

func (s *loggerShim) Warnln(msg ...interface{}) {
	level.Warn(s.logger).Log("msg", fmt.Sprint(spaceSep(msg)...))
}

func (s *loggerShim) Errorln(msg ...interface{}) {
	level.Error(s.logger).Log("msg", fmt.Sprint(spaceSep(msg)...))
}

func (s *loggerShim) Debugf(format string, args ...interface{}) {
	level.Debug(s.logger).Log("msg", fmt.Sprintf(format, args...))
}

func (s *loggerShim) Infof(format string, args ...interface{}) {
	level.Info(s.logger).Log("msg", fmt.Sprintf(format, args...))
}

func (s *loggerShim) Warnf(format string, args ...interface{}) {
	level.Warn(s.logger).Log("msg", fmt.Sprintf(format, args...))
}

func (s *loggerShim) Errorf(format string, args ...interface{}) {
	level.Error(s.logger).Log("msg", fmt.Sprintf(format, args...))
}

// WithFields will return a new logger derived from the original
// kitlog logger, with the provided fields added to the log string,
// as a key-value pair
func (s *loggerShim) WithFields(fields log.Fields) log.Logger {
	var keyvals []interface{} // nolint: prealloc

	for key, value := range fields {
		keyvals = append(keyvals, key, value)
	}

	return &loggerShim{
		logger: kitlog.With(s.logger, keyvals...),
	}
}
